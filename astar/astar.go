package astar

import (
	"math"
	"sort"
)

// Vertex a vertex represents a point in our grid.
type Vertex struct {
	X int
	Y int
}

// Node anode represents a particular unit.
type Node struct {
	Parent           *Node
	Position         *Vertex
	DistanceToTarget int
	Cost             int
	Weight           int
	Walkable         bool
}

// F calculates the distance and the cost of movement.
func (n *Node) F() int {
	if n.DistanceToTarget != -1 && n.Cost != -1 {
		return n.DistanceToTarget + n.Cost
	}

	return -1
}

// Stack A basic stack structure.
type Stack struct {
	nodes []*Node
	count int
}

// Astar is an implmentation of Astar pathfinding in GoLang
type Astar struct {
	Grid [][]Node
}

// Push adds a node to the stack.
func (s *Stack) Push(n *Node) {
	s.nodes = append(s.nodes, n)
	s.count++
}

// Pop removes and returns a node from the stack in last to first order.
func (s *Stack) Pop() *Node {
	if s.count == 0 {
		return nil
	}
	node := s.nodes[s.count-1]
	s.count--
	return node
}

// Contains determins if a slice contains a particular entry.
func Contains(a []*Node, x *Node) bool {
	for _, n := range a {
		if x.Position.X == n.Position.X && x.Position.Y == n.Position.Y {
			return true
		}
	}
	return false
}

// RemoveIndex Removes an entry from a slice
func (a *Astar) RemoveIndex(s []*Node, index int) []*Node {
	return append(s[:index], s[index+1:]...)
}

func indexOf(element *Node, data []*Node) int {
	for k, v := range data {
		if element.Position.X == v.Position.X && element.Position.Y == v.Position.Y {
			return k
		}
	}
	return -1 //not found.
}

// FindPath Finds an A* path to your desired endLocation if possible.
func (a *Astar) FindPath(startLocation *Vertex, endLocation *Vertex, avaliableAp int) *Stack {
	start := &Node{
		Parent:           nil,
		Position:         startLocation,
		DistanceToTarget: -1,
		Cost:             1,
		Weight:           1,
		Walkable:         true,
	}

	end := &Node{
		Parent:           nil,
		Position:         endLocation,
		DistanceToTarget: -1,
		Cost:             1,
		Weight:           1,
		Walkable:         true,
	}

	path := &Stack{
		nodes: make([]*Node, 0),
	}
	openList := []*Node{}
	closedList := []*Node{}
	adjacencies := []*Node{}

	current := start

	openList = append(openList, start)

	for len(openList) != 0 && Contains(closedList, end) == false {
		current = openList[0]
		openList = a.RemoveIndex(openList, 0)
		closedList = append(closedList, current)
		adjacencies = a.GetAdjacentNodes(current)

		for _, n := range adjacencies {
			if Contains(closedList, n) == false && n.Walkable {
				if Contains(openList, n) == false {
					n.Parent = current
					n.DistanceToTarget = int(math.Abs(float64(n.Position.X-end.Position.X)) + math.Abs(float64(n.Position.Y-end.Position.Y)))
					n.Cost = n.Weight + n.Parent.Cost
					openList = append(openList, n)
					sort.SliceStable(openList, func(i, j int) bool { return openList[i].F() < openList[j].F() })
				}
			}
		}
	}

	if Contains(closedList, end) == false {
		return nil
	}

	temp := closedList[indexOf(current, closedList)]

	if temp == nil {
		return nil
	}

	for {
		if temp == nil || temp == start {
			break
		}

		if temp.Cost <= avaliableAp+1 {
			path.Push(temp)
			temp = temp.Parent
		} else {
			return nil
		}
	}

	return path
}

// GridRows Returns the number of rows in the current grid.
func (a *Astar) GridRows() int {
	return len(a.Grid[0])
}

// GridCols Returns the number of columns in the current grid.
func (a *Astar) GridCols() int {
	return len(a.Grid)
}

// GetAdjacentNodes returns the adjacent nodes starting from a node.
func (a *Astar) GetAdjacentNodes(n *Node) []*Node {
	temp := []*Node{}

	row := n.Position.X
	col := n.Position.Y

	if row+1 < a.GridRows() {
		temp = append(temp, &a.Grid[col][row+1])
	}
	if row-1 >= 0 {
		temp = append(temp, &a.Grid[col][row-1])
	}
	if col-1 >= 0 {
		temp = append(temp, &a.Grid[col-1][row])
	}
	if col+1 < a.GridCols() {
		temp = append(temp, &a.Grid[col+1][row])
	}

	return temp
}
