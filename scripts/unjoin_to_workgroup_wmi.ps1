[CmdletBinding()]
param(
    [string]$Workgroup = "WORKGROUP"
)

$cs = Get-WmiObject -Class Win32_ComputerSystem

$r = $cs.UnjoinDomainOrWorkgroup(0, $null, $null)
if ($r.ReturnValue -ne 0) { exit $r.ReturnValue }

$r = $cs.JoinDomainOrWorkgroup($Workgroup, $null, $null)
exit $r.ReturnValue