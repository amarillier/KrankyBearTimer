#! /bin/sh

if [ $# -ge 1 ]
then
    ver=$1
else
    echo "Enter a version number"
    exit
fi

echo "version: $ver"
echo "main.go"
sed -i '' "s/Version = \".*\"/Version = \"$ver\"/" main.go

echo "FyneApp.toml"
sed -i '' "s/Version = \".*\"/Version = \"$ver\"/" FyneApp.toml

echo "Inno Setup Inno/TaniumClock.iss"
sed -i '' "s/MyAppVersion \".*\"/MyAppVersion \"$ver\"/" ./Inno/TaniumTimer.iss
