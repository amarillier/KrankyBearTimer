#! /bin/sh

version="0.9.2"
cp README.md installers/KrankyBearTimer/Resources
cp ReleaseNotes.txt installers/KrankyBearTimer/Resources
cd installers || exit
cp ../bin/WinAMD64/KrankyBearTimer.exe KrankyBearTimer
zip -r KrankyBearTimerWinAMD.zip KrankyBearTimer
rm KrankyBearTimer/KrankyBearTimer.exe

cp ../bin/MacOSAMD64/KrankyBearTimer KrankyBearTimer
zip -r KrankyBearTimerMacOSAMD.zip KrankyBearTimer
rm KrankyBearTimer/KrankyBearTimer

cp ../bin/MacOSARM64/KrankyBearTimer KrankyBearTimer
zip -r KrankyBearTimerMacOSARM.zip KrankyBearTimer
rm KrankyBearTimer/KrankyBearTimer

# see gh docs: https://cli.github.com/manual/gh_release_create
awk '/0.9.2/{flag=1}/^$/{flag=0}flag' ../ReleaseNotes.txt > latestReleaseNotes.txt
gh release create --title v"$version" v"$version" --draft --notes-file latestReleaseNotes.txt --prerelease KrankyBearTimerWinAMD.zip KrankyBearTimerMacOSAMD.zip KrankyBearTimerMacOSARM.zip KrankyBearTimerSetup.exe KrankyBearTimerARM.dmg KrankyBearTimerIntel.dmg

echo "Created draft release $version"
echo "Remember to publish when ready"
echo "gh release edit v$version --draft=false --prerelease=false"