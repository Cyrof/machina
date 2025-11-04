param(
    [Parameter(Mandatory = $true)] [string]$ConfigPath,
    [switch]$DryRun,
    [switch]$VerboseMode
)

$ErrorActionPreference = "Stop"

function Log {
    param([string]$level, [string]$msg)
    $ts = (Get-Date).ToString('yyyy-MM-dd HH:mm:ss')
    $line = "[$ts] [$level] $msg"
    Write-Host $line
    Add-Content -Path $Global:LogPath -Value $line
}

# setup logging targets
$LogDirDefault = 'C:\Machina'
try {
    if (-not (Test-Path $LogDirDefault)) { New-Item -ItemType Directory -Path $LogDirDefault -Force | Out-Null }
} catch { }

$Global:LogPath = Join-Path $LogDirDefault "cleanup.log"

Log "INFO" "Starting cleanup with config: $ConfigPath"

if (-not (Test-Path $ConfigPath)) {
    Log "ERROR" "Config not found: $ConfigPath"
    exit 1
}

# Json only to avoid external modules 
try {
    $cfgRaw = Get-Content -LiteralPath $ConfigPath -Raw -ErrorAction Stop
    $cfg = $cfgRaw | ConvertFrom-Json -ErrorAction Stop
} catch {
    Log "ERROR" "Failed to parse JSON config: $($_.Exception.Message)"
    exit 1
}

if ($VerboseMode) {
    Log "DEBUG" "Parsed config: $($cfg | ConvertTo-Json -Depth 5)"
}

function Remove-PathSafe {
    param([string]$p)
    if (-not $p) { return }
    if ($DryRun) { Log "DRYRUN" "Would remove: $p"; return }

    try {
        if (Test-Path -LiteralPath $p) {
            $item = Get-Item -LiteralPath $p -ErrorAction Stop
            if ($item.PSIsContainer) {
                Remove-Item -LiteralPath $p -Recurse -Force -ErrorAction Stop
            } else {
                Remove-Item -LiteralPath $p -Force -ErrorAction Stop
            }
            Log "OK" "Removed: $p"
        } else {
            # try wildcard matches 
            $matches = Get-ChildItem -Path $p -ErrorAction SilentlyContinue
            if ($matches) {
                foreach ($m in $matches) {
                    if ($DryRun) {
                        Log "DRYRUN" "Would remove: $($m.FullName)"
                    } else {
                        Remove-Item -LiteralPath $m.FullName -Recurse -Force -ErrorAction Stop
                        Log "OK" "Removed: $($m.FullName)"
                    }
                }
            } else {
                Log "WARN" "Not found or no matches: $p"
            }
        }
    } catch {
        Log "ERROR" "Failed removing $p: $($_.Exception.Message)"
    }
}

function Remove-RegKeySafe {
    param([string]$rk)
    if (-not $rk) { return }
    if ($DryRun) { Log "DRYRUN" "Would remove registry key: $rk"; return }
    try {
        if (Test-Path $rk) {
            Remove-Item -Path $rk -Recurse -Force -ErrorAction Stop
            Log "OK" "Removed registry: $rk"
        } else {
            Log "WARN" "Registry key not found: $rk"
        }
    } catch {
        Log "ERROR" "Failed removing registry $rk: $($_.Exception.Message)"
    }
}

function Stop-Delete-ServiceSafe {
    param([string]$svcName)
    if (-not $svcName) { return }
    if ($DryRun) { Log "DRYRUN" "Woudl stop & delete service: $svcName"; return }
    try {
        $svc = Get-Service -Name $svcName -ErrorAction SilentlyContinue
        if ($svc) {
            if ($svc.Status -ne 'Stopped') {
                try { Stop-Service -Name $svcName -Force -ErrorAction Stop } catch {}
            }
            sc.exe delete $svcName | Out-Null
            Log "OK" "Deleted service: $svcName"
        } else {
            Log "WARN" "Service not found: $svcName"
        }
    } catch {
        Log "ERROR" "Failed deleting service $svcName: $($_.Exception.Message)"
    }
}

try {
    if ($cfg.LogDir) { foreach ($d in $cfg.LogDir) { Remove-PathSafe -p $d } }
    if ($cfg.TempFiles) { foreach ($f in $cfg.TempFiles) { Remove-PathSafe -p $f } }
    if ($cfg.ExtraPaths) { foreach ($p in $cfg.ExtraPaths) { Remove-PathSafe -p $p } }
    if ($cfg.RegKeys) { foreach ($rk in $cfg.RegKeys) { Remove-RegKeySafe -rk $rk } }
    if ($cfg.Services) { foreach ($s in $cfg.Services) { Stop-Delete-ServiceSafe -svcName $s } }
} catch {
    Log "Error" "Unhandled error during cleanup: $($_.Exception.Message)"
    exit 1
}

Log "OK" "Cleanup completed. Log: $Global:LogPath"
