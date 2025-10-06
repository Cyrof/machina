@echo off
setlocal

set "WG=%~1"
if "%WG%"=="" set "WG=WORKGROUP"

echo Unjoining this computer from the domain...
for /f "tokens=2 delims== " %%A in ('
  wmic computersystem call UnJoinDomainOrWorkgroup 0 ^| find "ReturnValue"
') do set "RV=%%A"
if not "%RV%"=="0" exit /b %RV%

echo Joining workgroup: %WG%...
for /f "tokens=2 delims== " %%A in ('
  wmic computersystem call JoinDomainOrWorkgroup name^="%WG%" ^| find "ReturnValue"
') do set "RV=%%A"
if not "%RV%"=="0" exit /b %RV%

exit /b 0
