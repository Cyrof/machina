@echo off 
setlocal

set "WG=%~1"
if "%WG%"=="" set "WG=WORKGROUP"

echo Unjoining this computer from the domain...
wmic.exe /interactive:off ComputerSystem Where "Name='%computersystem%'" Call UnJoinDomainOrWorkgroup FUnjoinOptions=0
if errorlevel 1 (
    echo Failed to unjoin domain (errorlevel %errorlevel%)
    exit /b %errorlevel%
)

echo Joining workgroup: %WG%...
wmic.exe /interactive:off ComputerSystem Where "Name='%computersystem%'" Call JoinDomainOrWorkgroup name="WORKGROUP"
if errorlevel 1 {
    echo Failed to join workgroup %WG% (errorlevel %errorlevel%)
    exit /b %errorlevel%
}

echo Successfully unjoined domain and joined workgroup: %WG%
exit /b 0