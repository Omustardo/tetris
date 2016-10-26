Golang demo to show use of basic glfw window, opengl graphics, and keyboard
input. When run on desktop it runs with regular OpenGL, and when served with
`gopherjs serve` or by compiling using gopherjs, it uses WebGL.   I didn’t make
any endgame conditions or tracking of score, but the gameplay works.

To run on desktop:

`go run main.go`

To run in browser:

`gopherjs serve`

and then go to http://localhost:8080/github.com/omustardo/tetris/webgl-tetris/

This is my first real experience with low level OpenGL. I added some useful
links in my [documentation
page](<https://github.com/Omustardo/docs/tree/master/opengl>).
