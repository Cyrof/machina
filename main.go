//go:build windows

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modNetapi32                = windows.NewLazySystemDLL("Netapi32.dll")
	procNetUnjoinDomain        = modNetapi32.NewProc("NetUnjoinDomain")
	procNetJoinDomain          = modNetapi32.NewProc("NetJoinDomain")
	procNetRenameMachineDomain = modNetapi32.NewProc("NetRenameMachineInDomain") // not used; we rename locally below

	modKernel32           = windows.NewLazySystemDLL("Kernel32.dll")
	procSetComputerNameEx = modKernel32.NewProc("SetComputerNameExW")
)

const (
	// COMPUTER_NAME_FORMAT
	ComputerNamePhysicalDnsHostname = 5

	// NetJoin/Unjoin flags (see lmjoin.h)
	NETSETUP_JOIN_DOMAIN      = 0x00000001
	NETSETUP_ACCT_CREATE      = 0x00000002
	NETSETUP_ACCT_DELETE      = 0x00000004
	NETSETUP_WIN9X_UPGRADE    = 0x00000010
	NETSETUP_DOMAIN_JOIN_IF_JOINED = 0x00000020
	NETSETUP_JOIN_UNSECURE    = 0x00000040
	NETSETUP_MACHINE_PWD_PASSED = 0x00000080
	NETSETUP_DEFER_SPN_SET    = 0x00000100
	NETSETUP_INSTALL_INVOCATION = 0x00040000
)

func mustAdmin() {
	// Quick admin check: try opening SCM with CREATE_SERVICE (fails when not admin)
	scm, err := windows.OpenSCManager(nil, nil, windows.SC_MANAGER_CREATE_SERVICE)
	if err != nil {
		fail("This tool must be run as Administrator (elevated). Error: %v", err)
	}
	windows.CloseServiceHandle(scm)
}

func utf16Ptr(s string) *uint16 {
	if s == "" {
		return nil
	}
	ptr, _ := windows.UTF16PtrFromString(s)
	return ptr
}

func callNetUnjoinDomain(server, account, password *uint16, flags uint32) error {
	r1, _, e := procNetUnjoinDomain.Call(
		uintptr(unsafe.Pointer(server)),
		uintptr(unsafe.Pointer(account)),
		uintptr(unsafe.Pointer(password)),
		uintptr(flags),
	)
	if r1 != 0 { // NERR_Success == 0
		return fmt.Errorf("NetUnjoinDomain failed: %v (code=%d)", e, r1)
	}
	return nil
}

func callNetJoinDomain(server, domain, account, password, ou *uint16, flags uint32) error {
	r1, _, e := procNetJoinDomain.Call(
		uintptr(unsafe.Pointer(server)),
		uintptr(unsafe.Pointer(domain)),
		uintptr(unsafe.Pointer(account)),
		uintptr(unsafe.Pointer(password)),
		uintptr(unsafe.Pointer(ou)),
		uintptr(flags),
	)
	if r1 != 0 { // NERR_Success == 0
		return fmt.Errorf("NetJoinDomain failed: %v (code=%d)", e, r1)
	}
	return nil
}

func setComputerNameEx(newName string) error {
	r1, _, e := procSetComputerNameEx.Call(
		uintptr(uint32(ComputerNamePhysicalDnsHostname)),
		uintptr(unsafe.Pointer(utf16Ptr(newName))),
	)
	if r1 == 0 {
		return fmt.Errorf("SetComputerNameExW failed: %v", e)
	}
	return nil
}

func rebootNow() {
	// simplest: use shutdown.exe
	exec.Command("shutdown", "/r", "/t", "0").Start()
}

func runPowershell(ps string, args ...string) error {
	all := append([]string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", ps}, args...)
	cmd := exec.Command("powershell.exe", all...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("powershell error: %v\n%s", err, string(out))
	}
	return nil
}

func cmdUnjoin(args []string) {
	fs := flag.NewFlagSet("unjoin", flag.ExitOnError)
	user := fs.String("user", "", `Domain user (e.g. "2d\sa2dadmin")`)
	pass := fs.String("pass", "", "Password for domain user")
	deleteAcct := fs.Bool("delete-acct", false, "Also delete the computer account in AD (if it exists)")
	reboot := fs.Bool("reboot", false, "Reboot after unjoin")
	fs.Parse(args)

	mustAdmin()

	// Try clean unjoin with credentials (if supplied). If not, still try; if it fails, fall back.
	var flags uint32 = 0
	if *deleteAcct {
		flags |= NETSETUP_ACCT_DELETE
	}
	var err error
	if *user != "" && *pass != "" {
		fmt.Println("Attempting clean NetUnjoinDomain with provided credentials...")
		err = callNetUnjoinDomain(nil, utf16Ptr(*user), utf16Ptr(*pass), flags)
	} else {
		fmt.Println("Attempting NetUnjoinDomain without credentials (may fail if AD contact required)...")
		err = callNetUnjoinDomain(nil, nil, nil, flags)
	}
	if err != nil {
		fmt.Printf("Clean unjoin failed (%v). Forcing WORKGROUP locally...\n", err)
		// Force to WORKGROUP using PowerShell Add-Computer with LOCAL credentials suppressed.
		// This path doesn't need domain contact.
		ps := `Add-Computer -WorkGroupName "WORKGROUP" -Force`
		if e := runPowershell(ps); e != nil {
			fail("Fallback to WORKGROUP failed: %v", e)
		}
	}

	fmt.Println("Unjoin complete (or forced to WORKGROUP).")
	if *reboot {
		fmt.Println("Rebooting...")
		rebootNow()
	}
}

func cmdRename(args []string) {
	fs := flag.NewFlagSet("rename", flag.ExitOnError)
	newName := fs.String("newname", "", "New hostname (e.g. NV-NWC-N1VS01)")
	reboot := fs.Bool("reboot", false, "Reboot after rename")
	fs.Parse(args)

	if *newName == "" {
		fail("missing --newname")
	}
	mustAdmin()

	fmt.Printf("Renaming computer to %q...\n", *newName)
	if err := setComputerNameEx(*newName); err != nil {
		fail("%v", err)
	}
	fmt.Println("Rename scheduled (takes effect after reboot).")
	if *reboot {
		rebootNow()
	}
}

func cmdSetDNS(args []string) {
	fs := flag.NewFlagSet("setdns", flag.ExitOnError)
	iface := fs.String("iface", "Ethernet", "Interface alias (e.g. Ethernet)")
	dns := fs.String("dns", "", "Primary DNS server IP (e.g. 192.168.100.60)")
	fs.Parse(args)

	if *dns == "" {
		fail("missing --dns")
	}
	mustAdmin()

	ps := fmt.Sprintf(`Set-DnsClientServerAddress -InterfaceAlias "%s" -ServerAddresses %s`, *iface, *dns)
	fmt.Printf("Setting DNS on interface %q to %s...\n", *iface, *dns)
	if err := runPowershell(ps); err != nil {
		fail("%v", err)
	}
	fmt.Println("DNS updated.")
}

func cmdRejoin(args []string) {
	fs := flag.NewFlagSet("rejoin", flag.ExitOnError)
	domain := fs.String("domain", "", "AD domain (e.g. 2d.com)")
	user := fs.String("user", "", `Domain user (e.g. "2d\sa2dadmin")`)
	pass := fs.String("pass", "", "Password")
	ou := fs.String("ou", "", `Target OU (e.g. OU=Workstations,DC=2d,DC=com)`)
	reboot := fs.Bool("reboot", false, "Reboot after join")
	fs.Parse(args)

	mustAdmin()
	if *domain == "" || *user == "" || *pass == "" {
		fail("missing required flags: --domain, --user, --pass")
	}

	// Flags: create account if missing, allow join even if previously joined
	flags := uint32(NETSETUP_JOIN_DOMAIN | NETSETUP_ACCT_CREATE | NETSETUP_DOMAIN_JOIN_IF_JOINED)
	fmt.Printf("Joining domain %q...\n", *domain)
	if err := callNetJoinDomain(nil, utf16Ptr(*domain), utf16Ptr(*user), utf16Ptr(*pass), utf16Ptr(*ou), flags); err != nil {
		fail("%v", err)
	}
	fmt.Println("Domain join complete.")
	if *reboot {
		rebootNow()
	}
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	switch strings.ToLower(os.Args[1]) {
	case "unjoin":
		cmdUnjoin(os.Args[2:])
	case "rename":
		cmdRename(os.Args[2:])
	case "setdns":
		cmdSetDNS(os.Args[2:])
	case "rejoin":
		cmdRejoin(os.Args[2:])
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `adtool (Windows)

Commands:
  unjoin   --user "DOMAIN\admin" --pass "..." [--delete-acct] [--reboot]
  rename   --newname NEW-NAME [--reboot]
  setdns   --iface "Ethernet" --dns 192.168.100.60
  rejoin   --domain 2d.com --user "DOMAIN\admin" --pass "..." [--ou "OU=Workstations,DC=2d,DC=com"] [--reboot]

Examples:
  adtool unjoin --user "2d\sa2dadmin" --pass "P@ssw0rd" --delete-acct --reboot
  adtool rename --newname NV-NWC-N1VS01 --reboot
  adtool setdns --iface "Ethernet" --dns 192.168.100.60
  adtool rejoin --domain 2d.com --user "2d\sa2dadmin" --pass "P@ssw0rd" --ou "OU=Workstations,DC=2d,DC=com" --reboot
`)
}

func fail(f string, a ...any) {
	fmt.Fprintf(os.Stderr, "ERROR: "+f+"\n", a...)
	os.Exit(1)
}
