param(
    [Parameter(Mandatory=$true)] 
    [string]$NewName,

    [switch]$Restart,
    [switch]$Registry,
)

$ErrorActionPreference = "Stop"

function Log($lvl, $msg) {
    $ts=(Get-Date).ToString("yyyy-MM-dd HH:mm:ss")
    Write-Host "[$ts][$lvl] $msg"
}

function Fail($m) {
    Log "ERROR" $m
    exit 1
}

function Set-RegHost($name) {
    $paths = @(
        "HKLM:\SYSTEM\CurrentControlSet\Control\ComputerName\ActiveComputerName",
        "HKLM:\SYSTEM\CurrentControlSet\Control\ComputerName\ComputerName",
        "HKLM:\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters"
    )
    Log INFO "Updating registry keys for hostname: $name"
    Set-ItemProperty -Path $paths[0] -Name 'ComputerName' -value $name
    Set-ItemProperty -Path $paths[1] -Name 'ComputerName' -value $name
    Set-ItemProperty -Path $paths[2] -Name 'Hostname' -value $name
    Set-ItemProperty -Path $paths[2] -Name 'NVHostname' -value $name
    Log INFO "Registry keys updated."
}

try {
    Log INFO "Current hostname: $env:COMPUTERNAME"
    Log INFO "Request hostname: $NewName"

    if ($Registry) {
        $cs = Get-CimInstance Win32_ComputerSystem
        if ($cs.PartOfDomain) {
            Log INFO "Domain detected: $($cs.Domain) - AD object will NOT be updated."
        }
        Set-RegHost $NewName
    } else {
        Log INFO "Using Rename-Computer to change hostname."
        Rename-Computer -NewName $NewName -Force -ErrorAction Stop
        Log INFO "Hostname updated via Rename-Computer."
    }

    if ($Restart) {
        Log INFO "Restarting now..."
        Restart-Computer -Force
    } else {
        Log INFO "Reboot required for hostname change to take effect."
    }
} catch {
    Fail $_.Exception.Message
}