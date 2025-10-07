param(
    [Parameter(Mandatory=$true)] [string]$DomainName, 
    [string]$OUPath,
    [string]$ComputerName = $env:COMPUTERNAME,
    [string]$User,
    [string]$Password,
    [switch]$PromptForCredentials,
    [string[]]$DNSServers,
    [switch]$Restart
)

$ErrorActionPreference = "Stop"
function Fail($m){ Write-Error $m; exit 1 }

if ($DNSServers) {
    Get-NetAdapter -Physical | Where-Object Status -eq 'UP' | ForEach-Object {
        Set-DnsClientServerAddress -InterfaceIndex $_.IfIndex -ServerAddresses $DNSServers
    }
}

if ($PromptForCredentials) {
    $cred = Get-Credential -Message "Enter credentials to join $DomainName"
} elseif ($User -and $Password) {
    $secpasswd = ConvertTo-SecureString $Password -AsPlainText -Force
    $cred = New-Object System.Management.Automation.PSCredential ($User, $secpasswd)
} else {
    Fail "Either provide User and Password or use -PromptForCredentials"
}

$params = @{
    ComputerName = $ComputerName
    DomainName = $DomainName
    Credential = $cred
    Force = $true
}
if ($OUPath) { $params.OUPath = $OUPath }

if ($Restart) { Restart-Computer -Force } else {
    Write-Host "Join complete. Reboot is typically required to apply changes."
}