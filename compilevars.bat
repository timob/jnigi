@echo off

set jdkpath=%1%

if x%jdkpath%==x goto fail

set CGO_CFLAGS=-I%jdkpath%/include -I%jdkpath%/include/win32

goto end

:fail
echo "jdk path not given" 
echo "usage: ./compilevars.bat <jdk path>"

:end
