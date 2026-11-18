
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

## Windows software OpenGL fallback

Fyne apps need OpenGL, which VMs and locked-down hosts often lack. KrankyBear Timer automatically probes for usable hardware OpenGL at launch on Windows and, only if that fails, falls back to a bundled Mesa3D software renderer and relaunches once. Nothing changes on a machine with a working GPU. The installer bundles the fallback into a `mesa-fallback` folder next to `KrankyBearTimer-windows.exe`. The renderer is Mesa's fast `llvmpipe`; on CPUs without SSE4.1 (e.g. Proxmox kvm64/qemu64) it uses the slower but safe `softpipe`. If the app hangs after the first frame, start it with `-mesa-gl-compat` (or rename `mesa-fallback\.force-mesa-compat.sample` to `.force-mesa-compat`); `-mesa-gl-performance` forces `llvmpipe`. On a known no-GPU host (e.g. a golden VM image), moving `opengl32.dll` and `libgallium_wgl.dll` from `mesa-fallback` next to the `.exe` skips the one-time probe-and-relaunch. Help > About shows which renderer is in use. The DLLs are not in git: run `build/tools/get-mesa.sh --install` on each build host.

# To-do / known problems
* Available for Windows with Setup.exe and MacOS .dmg installers, as well as zip files for portable apps 

# To-do / known problems
- Allow optional always on top, save in prefs - may not be possible on Mac
https://www.google.com/search?q=fyne+golang+always+on+top&oq=fyne+golang+always+on+top&gs_lcrp=EgZjaHJvbWUyBggAEEUYOTIKCAEQABiABBiiBDIKCAIQABiABBiiBDIKCAMQABiABBiiBDIKCAQQABiABBiiBNIBCDg5MTBqMGoxqAIAsAIA&sourceid=chrome&ie=UTF-8
- See ReleaseNotes.txt for all recent changes and future plans

- Known problems - needs OpenGL drivers on some Windows


# License
See LICENSE - GPL v3
This is 100% free for anyone to use or misuse any way you like with no warranty as
to suitability or anything else, other than it has no viruses when I compile and
commit to git. But you should always check and scan anything you download from the
internet for viruses anyway. Don't be reckless.

All KrankyBear icons, images, logos used are copyright (c) Allan Marillier, 2024, 2025 ...

<img width="737" alt="image" src="https://github.com/user-attachments/assets/2da952b2-5058-481c-8744-9c938206967b" />
<img width="1297" alt="image" src="https://github.com/user-attachments/assets/1e8bf129-b981-48ca-9b25-d1def88e85a4" />
<img width="977" alt="image" src="https://github.com/user-attachments/assets/33e2cf63-8c0a-42df-ae93-91e78afc5008" />
