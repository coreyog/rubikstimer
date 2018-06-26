RubiksTimer [-u|--undecorated]

Adding an undecorated flag will remove the border from the window.
Use Escape to close the program.

Hold both control keys on your keyboard to arm the timer.
The timer will start when you release either control key.
Press both control keys at the same time again to stop the timer.
Pressing R will restart the timer and wait for both controls to be pressed again.
F12 will flip between a Black and Magenta background (for use as a chroma key).

Compile without command prompt window (i.e. double clicking an exe will only build a GUI window and not a console window): go build -ldflags -H=windowsgui main.go