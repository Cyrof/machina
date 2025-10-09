param(
    [string]$Workgroup = "WORKGROUP",
    [switch]$Restart
)

$ErrorActionPreference = "Stop"

function Log($level, $msg) {
    $ts = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    Write-Host "$ts [$level] $msg"
}

function Fail($m) {
    Log "ERROR" $m
    exit 1
}
try {
    Log "INFO" "Starting unjoin process. Target workgroup: $Workgroup"
    $cs = Get-WmiObject -Class Win32_ComputerSystem -ErrorAction Stop
    Log "INFO" "Computer: $($cs.Name) Domain: $($cs.Domain) PartOfDomain: $($cs.PartOfDomain)"

    if ($cs.PartOfDomain){
        Log "INFO" "Unjoining domain..."
        $r = $cs.UnjoinDomainOrWorkgroup(0, $null, $null)
        if ($r.ReturnValue -ne 0) { Fail "Unjoin failed (ReturnValue=$($r.ReturnValue))" }
        Log "OK" "Successfully unjoined from domain."
    } else {
        Log "INFO" "Machine is not part of a domain."
    }

    Log "INFO" "Joining workgroup '$Workgroup'..."
    $r = $cs.JoinDomainOrWorkgroup($Workgroup, $null, $null)
    Log "OK" "Successfully joined workgroup '$WORKGROUP'."

    $cs2 = Get-WmiObject -Class Win32_ComputerSystem
    if (-not $cs2.PartOfDomain -and $cs2.Workgroup -eq $Workgroup) {
        Log "OK" "Verificatoin passed. Current workgroup: $($cs2.Workgroup)"
    } else {
        Fail "Verification failed. Current Domain: $($cs2.Domain) Workgroup: $($cs2.Workgroup)"
    }

    if ($Restart) {
        Log "INFO" "Restarting computer to apply changes..."
        Restart-Computer -Force
    } else {
        Log "INFO" "Unjoin completed. A reboot is recommended to apply changes."
    }
} catch {
    Fail "Unhandled error: $($_.Exception.Message)"
}