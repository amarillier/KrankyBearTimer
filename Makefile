.DEFAULT_GOAL := hello

hello:
	echo "Hello Tanium!"
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
.PHONY:tanium

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
	go build -ldflags="-w -s" -o TaniumTimer .
.PHONY:build


linuxamd64:
 	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/LinuxAMD64/
.PHONY:linuxamd64

linuxarm64:
 	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/LinuxARM64/
.PHONY:linuxarm64

macosamd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/MacOSAMD64/
.PHONY:macamd64

macosarm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/MacOSARM64/
.PHONY:macarm64

winamd64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" go build -ldflags="-w -s" -o bin/WinAMD64/
.PHONY:winamd64

winarm64:
	GOOS=windows GOARCH=arm64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" go build -ldflags="-w -s" -o bin/WinARM64/
.PHONY:winarm64


buildall: linuxamd64 linuxarm64 macosamd64 macosarm64 winamd64 winarm64
.PHONY:buildall

clean:
	rm bin/*/*
.PHONY:clean

doc:
	grep -e "^// .* \.\.\. .*" -e "^.. .* \.\.\." -e "^func .*" *.go | tee doc.md
.PHONY:doc

