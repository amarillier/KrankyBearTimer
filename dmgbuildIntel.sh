#! /bin/sh
# compile, then create a dmg package
# https://github.com/create-dmg/create-dmg

# go build .
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/MacOSAMD64/

# set executable icon
./setIcon.sh KrankyBear.png bin/MacOSAMD64/KrankyBearTimer

cp bin/MacOSAMD64/KrankyBearTimer KrankyBearTimer.app/Contents/MacOS

test -f KrankyBearTimerIntel.dmg && rm KrankyBearTimerIntel.dmg
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

# set dmg icon
./setIcon.sh KrankyBear.png KrankyBearTimerIntel.dmg
if [ ! -d installers ]
then
  mkdir installers
fi
cp KrankyBearTimerIntel.dmg installers

# "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942