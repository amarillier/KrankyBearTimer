
# KrankyBear Timer  ![image](https://github.com/user-attachments/assets/9d653e80-6a4a-429d-a669-d607ba1bfcf3)

preferences stored via fyne preferences API land in
* MacOS: ~/Library/Preferences/fyne/com.github.amarillier.KrankyBearTimer/preferences.json
* Windows: ~\AppData\Roaming\fyne\com.github.amarillier.KrankyBearTimer\preferences.json
* MacOS resource location (sounds and backgrounds): /Applications/KrankyBear Timer.app/Contents/Resources
* Theme preferences are in ~/Library/Preferences/fyne/settings.json


## Features

Basic list - see help in the app for more detail - this is a countdown timer. Default times shown below can be modified to user preferred defaults via settings.
* Ad hoc time settable in 5 minute steps
* Bio break timer 10 minutes
* Lunch break timer 60 minutes
* Break timer to end at user specified time
* Notifications when the timer is done
* Color highlight when time is running out
* Optional customizable desktop clock with seconds, UTC time, date, hourly chime
* System tray access
* Optional desktop clock embedded - basically identical clock codebase from https://github.com/amarillier/KrankyBearClock allowing running a single timer / clock application

# To-do / known problems
* Available for Windows with Setup.exe and MacOS .dmg installers, as well as zip files for portable apps 

# To-do / known problems
- Allow optional always on top, save in prefs - may not be possible on Mac
https://www.google.com/search?q=fyne+golang+always+on+top&oq=fyne+golang+always+on+top&gs_lcrp=EgZjaHJvbWUyBggAEEUYOTIKCAEQABiABBiiBDIKCAIQABiABBiiBDIKCAMQABiABBiiBDIKCAQQABiABBiiBNIBCDg5MTBqMGoxqAIAsAIA&sourceid=chrome&ie=UTF-8
- See ReleaseNotes.txt for all recent changes and future plans

- Known problems - needs OpenGL drivers on some Windows


# License
See license.txt
This is 100% free for anyone to use or misuse any way you like with no warranty as
to suitability or anything else, other than it has no viruses when I compile and
commit to git. But you should always check and scan anything you download from the
internet for viruses anyway. Don't be reckless.
All KrankyBear icons, images, logos used are copyright (c) Allan Marillier, 2024, 2025 ...

<img width="737" alt="image" src="https://github.com/user-attachments/assets/2da952b2-5058-481c-8744-9c938206967b" />
<img width="1297" alt="image" src="https://github.com/user-attachments/assets/1e8bf129-b981-48ca-9b25-d1def88e85a4" />
<img width="977" alt="image" src="https://github.com/user-attachments/assets/33e2cf63-8c0a-42df-ae93-91e78afc5008" />
