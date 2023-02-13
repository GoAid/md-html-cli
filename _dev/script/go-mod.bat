@echo off
color 07
title 更新并整理 GO 模块依赖
:: file-encoding=GBK
rem by iTanken

cd /d %~dp0/../../
echo 1. 更新三方依赖...
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

echo 2. 整理模块依赖...
go mod tidy & echo.

:: echo 3. 导入模块依赖到 vendor 目录...
:: go mod vendor & echo.

git add .

call "%~dp0/done-time-pause.bat"
