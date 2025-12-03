@echo off

cd %~dp0

go mod tidy


:serverd
go build -trimpath -ldflags="-w -s" -v -o release\stealth-dns.exe main.go
IF %ERRORLEVEL% NEQ 0 goto :exit
if not exist release\etc mkdir release\etc
copy  etc\*.* release\etc
copy  sdk\nhp-agent.* release\sdk

:exit
IF %ERRORLEVEL% NEQ 0 (
    echo [Error] %ERRORLEVEL%
) ELSE (
    echo [Done] StealthDNS for platform %OS% built!
)