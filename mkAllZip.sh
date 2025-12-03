#! /bin/sh

version="0.9.5"
cp README.md installers/KrankyBearTimer/Resources
cp ReleaseNotes.txt installers/KrankyBearTimer/Resources
cd installers || exit
if [ ! -d KrankyBearTimer ]
then
    mkdir -p KrankyBearTimer
fi
rm KrankyBearTimer/KrankyBearTimer*

rm KrankyBearTimer/timer.exe
cp ../bin/KrankyBearTimer-windows.exe KrankyBearTimer/timer.exe
zip -r KrankyBearTimerWinAMD64.zip KrankyBearTimer
rm KrankyBearTimer/timer.exe


cp ../bin/KrankyBearTimer-macos-amd64 KrankyBearTimer/timer
zip -r KrankyBearTimerMacOSAMD64.zip KrankyBearTimer
rm KrankyBearTimer/timer

cp ../bin/KrankyBearTimer-macos-arm64 KrankyBearTimer/timer
zip -r KrankyBearTimerMacOSARM64.zip KrankyBearTimer
rm KrankyBearTimer/timer

cp ../bin/KrankyBearTimer-linux-amd64 KrankyBearTimer/timer
zip -r KrankyBearTimerLinuxAMD64.zip KrankyBearTimer
rm KrankyBearTimer/timer

cp ../bin/KrankyBearTimer-linux-arm64 KrankyBearTimer/timer
zip -r KrankyBearTimerLinuxARM64.zip KrankyBearTimer
rm KrankyBearTimer/timer

if [ $# -eq 1 ] && [ "$1" = "-release" ]
then
    # see gh docs: https://cli.github.com/manual/gh_release_create
    awk '/0.9.5/{flag=1}/^$/{flag=0}flag' ../ReleaseNotes.txt > latestReleaseNotes.txt
    gh release create --title v"$version" v"$version" --draft --notes-file latestReleaseNotes.txt --prerelease krankybear-resources_1.0.0-1_amd64.deb KrankyBearTimer_0.9.5-1_x86_64.rpm krankybear-resources_1.0.0-1_x86_64.rpm KrankyBearTimerLinuxAMD64.zip KrankyBearTimerLinuxARM64.zip KrankyBearTimer_0.9.5-1_aarch64.rpm     KrankyBearTimerMacOSAMD64.zip KrankyBearTimer_0.9.5-1_amd64.deb KrankyBearTimerMacOSARM64.zip KrankyBearTimer_0.9.5-1_amd64.pkg KrankyBearTimerSetup.exe KrankyBearTimer_0.9.5-1_arm64.deb KrankyBearTimerWinAMD64.zip KrankyBearTimer_0.9.5-1_arm64.pkg
    echo "Created draft release $version"
    echo "Remember to publish when ready"
    echo "gh release edit v$version --draft=false --prerelease=false"
fi