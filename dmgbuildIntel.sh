#! /bin/sh
# compile, then create a dmg package
# https://github.com/create-dmg/create-dmg

# go build .
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/MacOSAMD64/

# set executable icon
./setIcon.sh Resources/Images/KrankyBearFedoraRed.png bin/MacOSAMD64/KrankyBearTimer

test -f KrankyBearTimerIntel.dmg && rm KrankyBearTimerIntel.dmg
unzip -o timer.zip
cp bin/MacOSAMD64/KrankyBearTimer KrankyBearTimer.app/Contents/MacOS
#   --volicon "KrankyBearTimer.icns" \
create-dmg \
  --volname "KrankyBearTimer" \
  --window-pos 200 120 \
  --window-size 800 400 \
  --icon-size 100 \
  --icon "KrankyBearTimer.app" 200 190 \
  --hide-extension "KrankyBearTimer.app" \
  --app-drop-link 600 185 \
  --eula license.txt \
  "KrankyBearTimerIntel.dmg" \
  "KrankyBearTimer.app"
  # --add-file KrankyBearTimer.app ./KrankyBearTimer.app
  # "./"
rm -rf KrankyBearTimer.app

# set dmg icon
./setIcon.sh Resources/Images/KrankyBearFedoraRed.png KrankyBearTimerIntel.dmg
if [ ! -d installers ]
then
  mkdir installers
fi
cp KrankyBearTimerIntel.dmg installers

# "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942