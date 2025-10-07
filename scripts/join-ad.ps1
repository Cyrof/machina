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
func Log($level, $msg) {
    $ts = (Get-Date).ToString('yyyy-MM-dd HH:mm:ss')
    Write-Host "[$ts] [$level] $msg"
}

function Fail($m){ 
    Log "ERROR" $m
    exit 1 
}

try{
    Log "INFO" "Starting domain join process for $ComputerName to $DomainName"

    if ($DNSServers) {
        try {
            Log "INFO" "Setting DNS servers to: $($DNSServers -join ', ')"
            Get-NetAdapter | Where-Object Status -eq 'UP' | ForEach-Object {
                Set-DnsClientServerAddress -InterfaceIndex $_.IfIndex -ServerAddresses $DNSServers
            }
            Log "OK" "DNS servers updated successfully."
        } catch {
            Fail "Failed to set DNS servers: $($_.Exception.Message)"
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

    Log "INFO" "Joining computer '$ComputerName' to domain '$DomainName'"
    Add-Computer -ComputerName $ComputerName -DomainName $DomainName -Credential $cred -ErrorAction Stop -Force
    LOG "OK" "Successfully joined $ComputerName to $DomainName"

    $cs = Get-WmiObject Win32_ComputerSystem
    if ($cs.PartOfDomain -and $cs.Domain -ieq $DomainName) {
        LOG "OK" "Computer is now part of domain '$DomainName'"
    } else {
        Fail "Join command ran but verification failed. Current domain: $($cs.Domain)"
    }

    if ($Restart) {
        LOG "InFO" "Restarting computer to apply changes..."
        Restart-Computer -Force
    } else {
        LOG "INFO" "Restart not requested. Please restart the computer manually to apply changes."
    }
} catch {
    Fail "Unhandled error: $($_.Exception.Message)"
}