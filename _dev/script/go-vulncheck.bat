@echo off
color 07
title ©�����
:: file-encoding=GBK
rem by iTanken

rem ��ȡ go �汾
set gv=99999999999999999999
for /f "tokens=*" %%i in ('go version') do (
    set gv=%%i
)
set ver=%gv:~13,5%
:del-right
if "%ver:~-1%" equ "." set ver=%ver:~0,-1%&&goto del-right
if "%ver:~-1%" equ " " set ver=%ver:~0,-1%&&goto del-right
:goon
rem go �汾����С�� 1.18
if %ver% leq 1.17 (
  color 04
  echo. & echo ��ʹ�� go1.18 �����ϰ汾����©����飡 & echo.
  pause & exit
)

echo ��ʼ����©�����... & echo.

cd /d %~dp0/../../
go run golang.org/x/vuln/cmd/govulncheck@latest ./...

call "%~dp0/done-time-pause.bat"
