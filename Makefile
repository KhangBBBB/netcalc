ifeq ($(OS),Windows_NT)
    OS_DETECTED = Windows
else
    OS_DETECTED = $(shell uname -s)
endif

all:
ifeq ($(OS_DETECTED),Windows)
	go build -ldflags "-w -s -H=windowsgui" .
else ifeq ($(OS_DETECTED),Darwin)
	go build -ldflags "-w -s" .
	applify -author "netcalc authors" -name "netcalc" -version "0.0.5" ./netcalc
else ifeq ($(OS_DETECTED),Linux)
	go build -ldflags "-w -s" .
endif

clean:
	go clean
	rm -rf netcalc.app .DS_Store build
	rm -rf netcalc.exe
