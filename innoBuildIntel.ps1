# compile, then create an Inno setup installer package

# go build .
# GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" go build -ldflags="-w -s -H windowsgui -r TaniumClock.rc" -o bin/WinAMD64/
go build -ldflags="-w -s -H windowsgui -r TaniumTimer.rc" -o bin/WinAMD64/

Copy-Item bin/WinAMD64/TaniumTimer.exe .\TaniumTimer.exe

& 'C:\Program Files (x86)\Inno Setup 6\ISCC.exe' .\Inno\TaniumTimer.iss
