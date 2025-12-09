@echo off

cd %~dp0

go mod tidy


:serverd
go build -trimpath -ldflags="-w -s" -v -o release\stealth-dns.exe main.go
IF %ERRORLEVEL% NEQ 0 goto :exit
if not exist release\etc mkdir release\etc
if not exist release\etc\cert mkdir release\etc\cert
copy  etc\*.* release\etc
copy  sdk\nhp-agent.* release\sdk
copy  etc\cert\rootCA.pem release\etc\cert\

:exit
IF %ERRORLEVEL% NEQ 0 (
    echo [Error] %ERRORLEVEL%
) ELSE (
    echo [Done] StealthDNS for platform %OS% built!
)