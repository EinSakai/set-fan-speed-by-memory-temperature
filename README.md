# set-fan-speed-by-memory-temperature
set nvidia rtx video card fan speed by memory temperature

build:
CGO_ENABLED=1 GOARCH=amd64 GOOS=windows go build -a -o bin/set-gpu-fan main.go