@echo off
:: file-encoding=GBK
rem by iTanken
rem ���� footer �ű������ڴ�ӡ���ʱ�估���� BatNoPause ���������ж��Ƿ���ʾ�����������

color 0f
rem ��ȡ��ǰʱ��
set _now=%date:~0,4%-%date:~5,2%-%date:~8,2% %time:~0,2%:%time:~3,2%:%time:~6,2%.%time:~9,3%

rem ��ȡ���ڱ���
for /f "usebackq delims=" %%t in (`powershell -noprofile -c "[Console]::Title.Replace(' - '+[Environment]::CommandLine,'')"`) do (
  set _title=%%t
)

rem ��ӡ�����Ϣ
echo. & echo [%_now%] ִ��%_title%��ɣ� & echo.
if "%BatNoPause%" NEQ "1" (
  pause
)
