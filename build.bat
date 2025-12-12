@echo off

set "ROOT_DIR=%~dp0"
cd /d "%ROOT_DIR%"


if "%1"=="ui" goto :buildui
if "%1"=="full" goto :buildfull
if "%1"=="" goto :builddns
goto :builddns


:builddns
echo [StealthDNS] Initializing...
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
go build -trimpath -ldflags="-w -s" -v -o release\stealth-dns.exe main.go
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

call npm install
IF %ERRORLEVEL% NEQ 0 exit /b %ERRORLEVEL%


cd ..


echo [StealthDNS] Running wails build...
call wails build -platform windows/amd64
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
echo [Done] StealthDNS for platform %OS% built!
goto :eof

:exit
echo [Error] Build failed with error code %ERRORLEVEL%
