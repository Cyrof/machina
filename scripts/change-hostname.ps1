param(
    [Parameter(Mandatory=$true)] 
    [string]$NewHostname

    [switch]$Restart
)

$ErrorActionPreference = "Stop"

Write-Host "[*] Current hostname: $env:COMPUTERNAME"
Write-Host "[*] Changing hostname to: $NewHostname"

try {
    Rename-Computer -NewName $NewHostname -Force -ErrorAction Stop
    Write-Host "[+] Hostname updated in registry."
} catch {
    Write-Error "[x] Failed to change hostname: $($_.Exception.Message)"
    exit 1
}

if ($Restart) {
    Write-Host "[*] Restarting to apply hostname..."
    Restart-Computer -Force
} else {
    Write-Host "[i] Reboot required for hostname change to take effect."
}