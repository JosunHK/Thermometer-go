package main

import (
	"fmt"

	"github.com/golang-collections/collections/stack"
)

type Node struct {
	Points
	next   *Node
	prev   *Node
	active bool
}

type Points struct {
	x int
	y int
}

type r struct {
	this Points
	to   Points
}

var heads = []Points{
	{0, 2},
	{1, 0},
	{1, 1},
	{1, 2},
	{2, 3},
	{3, 3},
}

var relations = []r{
	{Points{0, 2}, Points{0, 1}},
	{Points{0, 1}, Points{0, 0}},
	{Points{1, 0}, Points{2, 0}},
	{Points{2, 0}, Points{3, 0}},
	{Points{1, 1}, Points{2, 1}},
	{Points{2, 1}, Points{3, 1}},
	{Points{1, 2}, Points{2, 2}},
	{Points{2, 3}, Points{1, 3}},
	{Points{1, 3}, Points{0, 3}},
	{Points{1, 3}, Points{0, 3}},
	{Points{3, 3}, Points{3, 2}},
}

var targetX = []int{2, 2, 2, 1}
var targetY = []int{2, 2, 1, 2}

func main() {
	const N int = 4
	m := make([][]*Node, N)
	for i := range m {
		m[i] = make([]*Node, N)
	}

	for _, head := range heads {
		m[head.y][head.x] = &Node{Points: head, active: false}
	}

	for _, r := range relations {
		prev := m[r.this.y][r.this.x]
		next := &Node{Points: r.to, prev: prev, active: false}
		prev.next = next
		m[r.this.y][r.this.x] = prev
		m[r.to.y][r.to.x] = next
	}

	play(m, N)

	printMap(m)
}

func play(m [][]*Node, N int) {
	nodeStack := stack.New()
	currentX := make([]int, N)
	currentY := make([]int, N)

	for _, head := range heads {
		//start with the max possible head
		x, y := head.x, head.y
		for {
			nodeStack.Push(m[y][x])
			currentX[head.x]++
			currentY[head.x]++
			boom := isBoom(currentX[head.x], currentY[head.y], head.x, head.y)
			if boom {

			}

		}

		for {
			currentX[head.x]++
			currentY[head.x]++

			boom := isBoom(currentX[head.x], currentY[head.y], head.x, head.y)
			if boom {
				nodeStack.Pop()
			}
		}
	}
}

func isGG(currentX, currentY []int) bool {
	for i := range currentX {
		if currentX[i] != targetX[i] || currentY[i] != targetY[i] {
			return false
		}
	}

	return true
}

func isBoom(valX, valY, indexX, indexY int) bool {
	if valX > targetX[indexX] || valY > targetY[indexY] {
		return true
	}

	return false
}

func printMap(m [][]*Node) {
	for _, row := range m {
		for _, node := range row {
			fmt.Print("┋")
			if node != nil {
				if node.prev != nil {
					if node.prev.y == node.y {
						fmt.Print(" ━ ┋")
					} else {
						fmt.Print(" | ┋")
					}
				} else {
					fmt.Print(" o ┋")
				}
			} else {
				fmt.Print("   ┋")
			}
		}
		fmt.Println()
		for range row {
			fmt.Print("┉")
			fmt.Print("┉")
			fmt.Print("┉")
			fmt.Print("┉")
			fmt.Print("┉")
		}
		fmt.Println()
	}
}
