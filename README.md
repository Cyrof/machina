# Machina
`machina` is a CLI tool for automating common Window system administration tasks.

Currently it supports **unjoining from AD**, **joining AD**, and **changing the hostname**.

---

## Installation 
1. Download the binary from the [releases](https://github.com/Cyrof/machina/releases) page.
2. Place it anywhere on your Windows machine. 
3. Run it from Powershell or Command Prompt.

---

## Prerequisites
For the `unjoin` and `join` commands: 
- Ensure the client can resolve the AD domain with DNS:
    ```powershell
    nslookup your-domain.com
    ```
- Ensure the client DNS points to either:
    - Your central DNS server, or 
    - The AD server IP directly.
Run all commands in **elevated shell (Run as Administrator)**.

---

## Usage 
### 1. Unjoin
Unjoins the computer from an AD domain and places it into a workgroup (WORKGROUP by default).
```powershell
.\machina unjoin --restart
```
- Prompts for confirmation before running.
- Required a local Administrator account for login after unjoining.
- You can specify the different workgroup with `--workgroup NAME`.

### 2. Join 
Joins the computer to an AD domain.
```powershell
# Interactive secure credential prompt 
.\machina join --domain "your-domain.com" --prompt

# Non-interactive (explicit user + password)
.\machina join --domain "your-domain.com" --user "DOMAIN\admin" --password "Secret123"
```
Options:
- `--domain` (_required_): AD domain to join.
- `--prompt`: Securely prompt for credentials.
- `--user` and `--password`: Non-interactive credentials (not recommended).
- `--dns`: Optional, set DNS servers before join (e.g., `--dns "1.1.1.1"`).
- `--restart`: Restart after successful join.

### 3. Hostname
Changes the computer hostname.
```powershell
# Change hostname (restart later)
.\machina hostname --name NEW-PC

# Change hostname and restart immediately
.\machina hostname --name NEW-PC --restart
```
#### Force rename via registry (bypasses domain restrictions)
Use this if the machine is still domain-joined and `Rename-Computer` is blocked:
```powershell
.\machina hostname --name NEW-PC --registry
```
You will be asked to confirm unless you use `--yes`:
```powershell
.\machina hostname --name NEW-PC --registry --yes
```
Options:
- `--name` (_require_): New hostname.
- `--restart`: Restart immediately after renaming.
- `--registry`: Force rename hostname via registry (WILL not change in AD).
- `--yes`: Skip confirmation when using `--registry`.

> **Important note**: Using `--registry` does **not** update the AD computer object. After reboot, the local machine will use the new name, but AD may still reference the old one until you manually update it or rejoin.

---

## Notes 
- All commands must be run with Administrator privileges.
- Both `unjoin` and `join` require a reboot before domain logon changes take effect.
- `hostname` requires a reboot before the new name is applied fully.