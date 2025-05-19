.DEFAULT_GOAL := hello
.ONESHELL:

hello:
	echo "Hello KrankyBear!"
	echo "make fmt to format the code"
	echo "make lint to run golint"
	echo "make vet to vet the code"
	echo "make run to run the main.go"
	echo "make build to build for the current system"
	# echo "	make linuxamd64"
	# echo "	make linuxarm64"
	echo "	make macamd64"
	echo "	make macarm64"
	echo "	make winamd64"
	echo "	make winarm64"
	echo "	make all"
	echo "make clean to remove compiled files from bin/*"
	echo "make doc to generate some docs based on func names"
	echo "  grepped | tee doc.md, on display and in file"
.PHONY:hello

fmt:
	go fmt ./...
.PHONY:fmt

lint:
	~/go/bin/golint ./...
.PHONY:lint

vet: fmt
	go vet ...
.PHONY:vet

run:
	go run .
.PHONY:run

# Supported cross compile GOOS and GOARCH https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
build:
	./setver.sh
	go build -ldflags="-w -s" -o KrankyBearTimer .
	./setIcon.sh KrankyBearTimer.png KrankyBearTimer
.PHONY:build


linuxamd64:
	echo "This doesn't work right now on Mac ARM or Win AMD64 - no action"
 	# GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/LinuxAMD64/
.PHONY:linuxamd64

linuxarm64:
	echo "This doesn't work right now on Mac ARM or Win AMD64 - no action"
 	# GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/LinuxARM64/
.PHONY:linuxarm64

macamd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/MacOSAMD64/
	./setIcon.sh KrankyBearTimer.png bin/MacOSAMD64/KrankyBearTimer
.PHONY:macamd64

macarm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/MacOSARM64/
	./setIcon.sh KrankyBearTimer.png bin/MacOSARM64/KrankyBearTimer
.PHONY:macarm64

winamd64:
	go-winres make
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" go build -ldflags="-w -s -H windowsgui -r KrankyBearTimer.rc" -o bin/WinAMD64/
.PHONY:winamd64

winarm64:
	echo "This doesn't work right now on Mac ARM or Win AMD64 - no action"
	# go-winres make
	# GOOS=windows GOARCH=arm64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" go build -ldflags="-w -s -H windowsgui -r KrankyBearTimer.rc" -o bin/WinARM64/
.PHONY:winarm64


buildall: linuxamd64 linuxarm64 macamd64 macarm64 winamd64 winarm64
.PHONY:buildall

dmg: 
	./dmgbuildIntel.sh
	./dmgbuildARM.sh
.PHONY:dmg

clean:
	rm bin/*/*
.PHONY:clean

doc:
	grep -e "^// .* \.\.\. .*" -e "^.. .* \.\.\." -e "^func .*" *.go | tee doc.md
.PHONY:doc

