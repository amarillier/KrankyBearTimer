#! /bin/sh

# fast update fyne and compile
go get fyne.io/fyne/v2@latest # or a specific version like @v2.4.0
go mod tidy
go mod vendor

rm KrankyBearTimer
rm timer
go build -ldflags="-w -s" -o KrankyBearTimer .
./setIcon.sh Resources/Images/KrankyBearFedoraRed.png KrankyBearTimer

killall timer
rm ~/bin/timer
cp KrankyBearTimer ~/bin/timer