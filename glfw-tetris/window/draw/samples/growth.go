// Draw "trees" in a forest.
package main

import (
	"flag"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/omustardo/window"
	"github.com/omustardo/window/draw"
)

var (
	MaxDepth int
	LeafSize int
)

func init() {
	flag.IntVar(&MaxDepth, "max_depth", 5, "Max depth for branches")
	flag.IntVar(&LeafSize, "leaf_size", 10, "Leaf size in pixels")

	flag.Parse()

	// OpenGl needs to run on one thread, evidently.
	// https://github.com/go-gl/gl/issues/13
	runtime.LockOSThread()

	rand.Seed(int64(time.Now().Nanosecond()))
}

func main() {
	// Set up gui. As the server is not meant to host a player, this is somewhat
	// unneeded, but it helps debugging to see the actual state of the world.
	gui, err := window.Initialize("Growth", 1920, 1080)
	if err != nil {
		log.Fatalln(err)
	}
	defer glfw.Terminate()

	var trees []*Tree
	for i := 0; i < int(rand.Int31n(15)+5); i++ {
		trees = append(trees, NewTree(300+rand.Float32()*1320, 0, rand.Float32()*50+75, MaxDepth+int(rand.Int31n(3))))
	}

	ticker := time.NewTicker(time.Second / 6)
	for !gui.ShouldClose() {
		draw.BeginDraw()
		for _, tree := range trees {
			tree.draw()
		}

		for _, tree := range trees {
			var children []*line
			for _, li := range tree.Lines {
				if li.MaxChildren > 0 && li.Depth < tree.MaxDepth && rand.Float32() > 0.2 { // 80% chance to make a new line
					children = append(children, li.makeChild())
					li.MaxChildren--
				}
			}
			for _, c := range children {
				tree.Lines = append(tree.Lines, c)
			}
		}

		gui.SwapBuffers()
		glfw.PollEvents()
		<-ticker.C
	}
}

type point struct {
	X, Y float32
}

type line struct {
	Start, End  point
	MaxChildren int
	Depth       int
}

func (li *line) draw() {
	draw.Line(li.Start.X, li.Start.Y, li.End.X, li.End.Y)
	if li.Depth == MaxDepth {
		draw.CircleColored(li.End.X, li.End.Y, 0, float64(LeafSize), 0, 1, 0, 1)
	}
}

func (li *line) length() float32 {
	return float32(math.Sqrt(float64(
		(li.End.X-li.Start.X)*(li.End.X-li.Start.X) +
			(li.End.Y-li.Start.Y)*(li.End.Y-li.Start.Y))))
}

// New line begins at the previous line's end point, points in a similar direction (within 60 degrees == pi/3 radians to either side), and is 60 to 90% of the length.
func (li *line) makeChild() *line {
	// Length of new line should be 60-90% of parent length
	length := (rand.Float64()*3/10 + 0.6) * float64(li.length())

	// Get angle of original line:
	angle := math.Atan2(float64(li.End.Y-li.Start.Y), float64(li.End.X-li.Start.X))
	// New line should be within 60 degrees (pi/3) in either direction.
	angle = angle - math.Pi/3 + rand.Float64()*math.Pi*2/3

	// So now we have direction and length of the new line, but we need to figure
	// out the actual point. The easiest way I know of to do this is to get the
	// unit vector in the direction we want, multiply it by the desired length,
	// and add the scaled up vector to the point we want to start from.
	return &line{
		Start:       li.End,
		End:         point{li.End.X + float32(length*math.Cos(angle)), li.End.Y + float32(length*math.Sin(angle))},
		MaxChildren: int(rand.Int31n(5)),
		Depth:       li.Depth + 1,
	}
}

type Tree struct {
	Lines    []*line
	MaxDepth int
}

func (t *Tree) draw() {
	for _, li := range t.Lines {
		li.draw()
	}
}

func NewTree(x, y, trunk_height float32, max_depth int) *Tree {
	return &Tree{
		Lines: []*line{
			{
				Start:       point{x, y},
				End:         point{x, y + trunk_height},
				MaxChildren: 3,
				Depth:       0,
			},
		},
		MaxDepth: max_depth,
	}
}
