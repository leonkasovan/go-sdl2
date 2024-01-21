# go-sdl2
Example of Go using SDL2 Library

## Build
```
go mod init go-sdl2
go mod tidy
go build test_event.go
go build test_joystick.go
```
## Test Joystick in Steamdeck
When the executable run in console, joystick won't work because joystick event is redirect as keyboard event.  
In order to run properly in Steamdeck, add the executable via "Add non steam Game" in Steam GUI.  
Set gamepad layout in Controller Setting, choose : Gamepad With Joystick Trackpad.  
Run it from Steam GUI.  
