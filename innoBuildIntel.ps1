# compile, then create an Inno setup installer package

# go build .
# GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" go build -ldflags="-w -s -H windowsgui -r KrankyBearClock.rc" -o bin/WinAMD64/
go-winres make
go build -ldflags="-w -s -H windowsgui -r KrankyBearTimer.rc" -o bin/WinAMD64/

Copy-Item bin/WinAMD64/KrankyBearTimer.exe .\KrankyBearTimer.exe

& 'C:\Program Files (x86)\Inno Setup 6\ISCC.exe' .\Inno\KrankyBearTimer.iss
