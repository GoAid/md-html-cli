@echo off
color 07
title ��̬���
:: file-encoding=GBK
rem by iTanken
echo ��ʼ���о�̬���... & echo.

cd /d %~dp0/../../
echo [revive.run]
go run github.com/mgechev/revive@latest -config ./_dev/config/revive.toml -exclude ./vendor/... -formatter stylish ./...
echo. & echo [staticcheck.io]
go run honnef.co/go/tools/cmd/staticcheck@latest -f text ./...

call "%~dp0/done-time-pause.bat"
