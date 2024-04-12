@echo off
setlocal enabledelayedexpansion

REM Enable verbose mode
@echo on

REM Execute the commands
upx dist\*\supervisord*
