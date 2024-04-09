# go-sdl2
Example of Go using SDL2 Library

## Build
```
go mod init go-sdl2
go mod tidy
go build test_event.go
go build test_joystick.go
go build test_render.go
go build -tags=gles2 test_opengles3.1.go

#Built for RG35XX
CGO_CFLAGS="-Os -marm -march=armv7-a -mtune=cortex-a9 -mfpu=neon-fp16 -mfloat-abi=hard" GOOS=linux GOARCH=arm CGO_ENABLED=1 go build test_render.go

#Built for RG353P
GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build test_render.go
```
## Test Joystick in Steamdeck
When the executable run in console, joystick won't work because joystick event is redirect as keyboard event.  
In order to run properly in Steamdeck, add the executable via "Add non steam Game" in Steam GUI.  
Set gamepad layout in Controller Setting, choose : Gamepad With Joystick Trackpad.  
Run it from Steam GUI.  
