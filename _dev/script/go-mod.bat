@echo off
color 07
title ���²����� GO ģ������
:: file-encoding=GBK
rem by iTanken

cd /d %~dp0/../../
echo 1. ������������...
cd
::& go get -d -u & echo.

go get github.com/FurqanSoftware/goldmark-d2@latest
go get github.com/PuerkitoBio/goquery@latest
go get github.com/alecthomas/chroma/v2@latest
go get github.com/jessevdk/go-flags@latest
go get github.com/stefanfritsch/goldmark-fences@latest
go get github.com/wzshiming/winseq@latest
go get github.com/yuin/goldmark@latest
go get github.com/yuin/goldmark-emoji@latest
go get github.com/yuin/goldmark-highlighting/v2@latest
go get go.abhg.dev/goldmark/mermaid@latest
go get oss.terrastruct.com/d2@latest

echo 2. ����ģ������...
go mod tidy & echo.

:: echo 3. ����ģ�������� vendor Ŀ¼...
:: go mod vendor & echo.

git add .

call "%~dp0/done-time-pause.bat"
