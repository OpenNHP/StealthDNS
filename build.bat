@echo off

set "ROOT_DIR=%~dp0"
cd /d "%ROOT_DIR%"

REM Version management
if exist .version (
    set /p BASE_VERSION=<.version
) else (
    set BASE_VERSION=1.0.0
)
for /f %%i in ('git rev-list --count HEAD 2^>nul') do set BUILD_NUMBER=%%i
if "%BUILD_NUMBER%"=="" set BUILD_NUMBER=0
for /f %%i in ('git rev-parse --short HEAD 2^>nul') do set COMMIT_ID=%%i
if "%COMMIT_ID%"=="" set COMMIT_ID=unknown
for /f "tokens=*" %%a in ('powershell -Command "Get-Date -Format 'yyyy-MM-dd_HH:mm:ss'"') do set BUILD_TIME=%%a
set GOMODULE=github.com/OpenNHP/StealthDNS/version
set "VERSION_LDFLAGS=-X %GOMODULE%.Version=%BASE_VERSION% -X %GOMODULE%.BuildNumber=%BUILD_NUMBER% -X %GOMODULE%.CommitID=%COMMIT_ID% -X %GOMODULE%.BuildTime=%BUILD_TIME%"

if "%1"=="ui" goto :buildui
if "%1"=="full" goto :buildfull
if "%1"=="" goto :builddns
goto :builddns


:builddns
echo [StealthDNS] Initializing...
echo [StealthDNS] Version: %BASE_VERSION% (Build: %BUILD_NUMBER%, Commit: %COMMIT_ID%)
git submodule update --init --recursive
go mod tidy

echo [StealthDNS] Building OpenNHP SDK from submodule...
IF %ERRORLEVEL% NEQ 0 exit /b %ERRORLEVEL%


echo [StealthDNS] Building Windows SDK (nhp-agent.dll)...
if not exist sdk mkdir sdk
set CGO_ENABLED=1

cd third_party\opennhp\nhp
go mod tidy
IF %ERRORLEVEL% NEQ 0 exit /b %ERRORLEVEL%

cd ..\endpoints
go mod tidy
IF %ERRORLEVEL% NEQ 0 exit /b %ERRORLEVEL%

go build -a -trimpath -buildmode=c-shared -ldflags="-w -s" -v -o ..\..\..\sdk\nhp-agent.dll agent\main\main.go agent\main\export.go
IF %ERRORLEVEL% NEQ 0 exit /b %ERRORLEVEL%

cd ..\..\..
echo [StealthDNS] Windows SDK built successfully!
echo [StealthDNS] Building StealthDNS...
echo [StealthDNS] Injecting version info: Version=%BASE_VERSION%, Build=%BUILD_NUMBER%, Commit=%COMMIT_ID%
go build -trimpath -ldflags="-w -s %VERSION_LDFLAGS%" -v -o release\stealth-dns.exe main.go
IF %ERRORLEVEL% NEQ 0 exit /b %ERRORLEVEL%
if not exist release\etc mkdir release\etc
if not exist release\etc\cert mkdir release\etc\cert
if not exist release\sdk mkdir release\sdk
copy  etc\*.* release\etc
copy  sdk\nhp-agent.* release\sdk
copy  etc\cert\rootCA.pem release\etc\cert\

if "%1"=="full" goto :buildui
goto :done

:buildui
echo [StealthDNS] Building UI...


cd /d "%ROOT_DIR%"


cd ui
IF %ERRORLEVEL% NEQ 0 (
    echo [Error] Cannot find ui directory
    exit /b %ERRORLEVEL%
)


call go mod tidy
IF %ERRORLEVEL% NEQ 0 exit /b %ERRORLEVEL%


cd frontend
IF %ERRORLEVEL% NEQ 0 (
    echo [Error] Cannot find frontend directory
    exit /b %ERRORLEVEL%
)

call npm ci
IF %ERRORLEVEL% NEQ 0 exit /b %ERRORLEVEL%


cd ..


echo [StealthDNS] Running wails build...
echo [StealthDNS UI] Injecting version info: Version=%BASE_VERSION%, Build=%BUILD_NUMBER%, Commit=%COMMIT_ID%
set "UI_LDFLAGS=-X stealthdns-ui/version.Version=%BASE_VERSION% -X stealthdns-ui/version.BuildNumber=%BUILD_NUMBER% -X stealthdns-ui/version.CommitID=%COMMIT_ID% -X stealthdns-ui/version.BuildTime=%BUILD_TIME%"
call wails build -ldflags="%UI_LDFLAGS%" -platform windows/amd64
IF %ERRORLEVEL% NEQ 0 exit /b %ERRORLEVEL%


cd /d "%ROOT_DIR%"


if not exist release mkdir release
if exist "ui\build\bin\stealthdns-ui.exe" (
    copy /Y "ui\build\bin\stealthdns-ui.exe" release\
    echo [Done] StealthDNS UI built and copied to release\
) else if exist "ui\build\bin\stealthdns-ui\stealthdns-ui.exe" (
    copy /Y "ui\build\bin\stealthdns-ui\stealthdns-ui.exe" release\
    echo [Done] StealthDNS UI built and copied to release\
) else (
    echo [Warning] Could not find stealthdns-ui.exe in expected locations
    echo Checking ui\build\bin\ contents:
    dir /b ui\build\bin\
)
goto :done

:buildfull
echo [StealthDNS] Building full package (DNS + UI)...
goto :builddns


:done
echo [StealthDNS] Cleaning up submodule changes...
cd /d "%ROOT_DIR%"
cd third_party\opennhp
git restore .
cd /d "%ROOT_DIR%"
echo [Done] StealthDNS for platform %OS% built!
goto :eof

:exit
echo [Error] Build failed with error code %ERRORLEVEL%
