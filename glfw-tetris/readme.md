Golang demo to show use of basic glfw window, opengl graphics, and keyboard
input.

You’ll need to install glfw and go-gl:

`go get github.com/go-gl/gl/v{3.2,3.3,4.1,4.4,4.5}-{core,compatibility}/gl`

`go get github.com/go-gl/gl/v3.3-core/gl`

`go get -u github.com/go-gl/glfw/v3.2/glfw`

 

and my little wrapper for glfw:

`go get -u github.com/Omustardo/window`

 

and then this repo:

`go get -u github.com/omustardo/tetris/glfw-tetris`

 

run tetris with:

`go run github.com/omustardo/tetris/glfw-tetris/main.go`

 

I didn’t make any endgame conditions or tracking of score, but the gameplay
works.
