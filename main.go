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

var m [][]*Node
var targetX = []int{2, 2, 2, 1}
var targetY = []int{2, 2, 1, 2}

func main() {
	const N int = 4
	m = make([][]*Node, N)
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

	// play(N)

	printMap(m)
}

func play(N int) error {
	nodeStack := stack.New()
	currentX := make([]int, N)
	currentY := make([]int, N)

	fmt.Println(exhaustNode(m[2][1], currentX, currentY))

	for i, h1 := range heads {
		curr := m[h1.y][h1.x]
		ptr := 0
		for {
			//back track
			if curr == nil {
				nodeStack.Pop()
				//func decrement curr node(curr, currentX, currentY)
				curr = nodeStack.Peek().(*Node).prev
				ptr--
				if ptr == i {
					ptr--
				}
			}

			//end reached, start back track
			if ptr == N {
				curr = curr.prev
			}

			eNode := exhaustNode(curr, currentX, currentY)
			nodeStack.Push(eNode)

			if i == ptr {
				ptr++
			}

			//get next head
			nextHead := heads[ptr]
			curr = m[nextHead.y][nextHead.x]
			ptr++

			//check end game
			if isGG(currentX, currentY) {
				return nil
			}
		}
	}

	return fmt.Errorf("no solution")
}

func exhaustNode(node *Node, currentX, currentY []int) *Node {
	boom := isBoom(currentX[node.x]+1, currentY[node.y]+1, node.x, node.y)
	if boom {
		return node.prev
	} else {
		node.active = true
		m[node.y][node.x] = node
		currentX[node.x]++
		currentY[node.y]++
		if node.next == nil {
			return node
		}

		return exhaustNode(node.next, currentX, currentY)
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
	//println color
	// red := "\033[31m"
	// Reset = "\033[0m"
	fmt.Printf("     ")
	for _, val := range targetX {
		fmt.Printf(" %v  ", val)
	}

	fmt.Println()
	fmt.Print("┉┉┉┉")

	for range targetX {
		fmt.Print("┉┉┉┉")
	}
	fmt.Println("┉")

	for i, row := range m {
		fmt.Printf("  %v ┋", targetY[i])
		for _, node := range row {
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
			fmt.Print("┉┉┉┉")
		}
		fmt.Print("┉┉┉┉")
		fmt.Println("┉")
	}
}
