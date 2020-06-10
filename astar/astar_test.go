package astar

import (
	"testing"
)

func generateGrid() [][]Node {
	var grid = [][]Node{}
	for y := 0; y <= 8; y++ {
		var row = []Node{}

		for x := 0; x <= 8; x++ {
			var node = Node{
				Position: &Vertex{
					X: x,
					Y: y,
				},
				Walkable: true,
			}

			row = append(row, node)
		}

		grid = append(grid, row)
	}

	return grid
}

func TestFindPath(t *testing.T) {
	var startingVertex = Vertex{
		X: 0,
		Y: 0,
	}

	var endingVertex = Vertex{
		X: 5,
		Y: 5,
	}

	var testingAstar = Astar{
		Grid: generateGrid(),
	}

	aStarStack := testingAstar.FindPath(&startingVertex, &endingVertex, 10)

	if aStarStack.count <= 0 {
		t.Fail()
	}
}
