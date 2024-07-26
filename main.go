package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/stack"
)

type Node struct {
	Points
	next   *Node
	prev   *Node
	active bool
}

type Points struct {
	y int
	x int
}

type r struct {
	this Points
	to   Points
}

// 10x10
var raw = "s-0-0,~-1-0,~-2-0,~-3-0,s-0-1,~-1-1,~-2-1,~-3-1,s-4-2,~-3-2,~-2-2,~-1-2,~-0-2,s-0-3,~-1-3,~-2-3,~-3-3,~-4-3,~-5-3,~-6-3,s-2-4,~-1-4,~-0-4,s-0-6,~-0-5,s-0-7,~-0-8,~-0-9,s-5-5,~-4-5,~-3-5,~-2-5,~-1-5,s-2-6,~-1-6,s-1-7,~-2-7,~-3-7,~-4-7,~-5-7,s-1-8,~-2-8,~-3-8,~-4-8,~-5-8,s-5-9,~-4-9,~-3-9,~-2-9,~-1-9,s-3-4,~-4-4,~-5-4,s-5-6,~-4-6,~-3-6,s-4-1,~-4-0,s-5-0,~-5-1,s-6-2,~-5-2,s-8-0,~-7-0,~-6-0,s-6-1,~-7-1,~-8-1,s-6-5,~-6-4,s-6-6,~-6-7,~-6-8,~-6-9,s-7-3,~-7-2,s-7-6,~-7-5,~-7-4,s-7-8,~-7-7,s-8-9,~-7-9,s-8-2,~-8-3,~-8-4,~-8-5,~-8-6,~-8-7,~-8-8,s-9-9,~-9-8,~-9-7,~-9-6,~-9-5,~-9-4,~-9-3,~-9-2,~-9-1,~-9-0;1,4,5,3,6,8,2,6,6,8;6,8,6,5,4,2,5,7,5,1"

// 6x6
// works
// var raw = "s-2-1,~-2-0,~-1-0,~-0-0,~-0-1,~-1-1,~-1-2,~-1-3,s-0-2,~-0-3,~-0-4,~-1-4,~-2-4,~-3-4,~-4-4,s-2-5,~-1-5,~-0-5,s-2-2,~-2-3,s-4-2,~-4-1,~-4-0,~-3-0,~-3-1,~-3-2,~-3-3,s-3-5,~-4-5,~-5-5,~-5-4,s-5-2,~-5-3,~-4-3,s-5-1,~-5-0;5,4,3,2,2,1;5,3,2,2,3,2"

// works
// can find solution
// var raw = "s-0-5,~-0-4,~-0-3,~-0-2,~-0-1,~-0-0,~-1-0,~-2-0,s-4-5,~-5-5,~-5-4,~-4-4,~-3-4,~-2-4,~-2-3,~-2-2,~-2-1,~-1-1,~-1-2,s-1-5,~-1-4,~-1-3,s-2-5,~-3-5,s-3-0,~-3-1,~-4-1,~-5-1,s-4-2,~-3-2,~-3-3,~-4-3,s-5-0,~-4-0,s-5-2,~-5-3;2,1,3,2,4,5;4,3,1,3,3,3"

// 4x4
// var raw = "s-0-0,~-0-1,s-0-2,~-0-3,s-1-0,~-1-1,~-1-2,~-1-3,s-2-0,~-2-1,s-2-2,~-2-3,s-3-3,~-3-2,~-3-1,~-3-0;2,1,3,1;3,1,1,2"

var heads []Points

var relations []r

var N int

var m [][]*Node
var nodeStack = stack.New()
var currentX []int
var currentY []int

var targetX []int
var targetY []int

func main() {
	initData()
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

	printMap()

	if err := play(); err != nil {
		fmt.Println(err)
		printMap()
	} else {
		fmt.Println("GG")
		printMap()
	}
}

func play() error {
	count := 0
	for i, h := range heads {
		headsToCheck := []Points{}
		headsToCheck = append(headsToCheck, heads[:i]...)
		headsToCheck = append(headsToCheck, heads[i+1:]...)
		currentX = make([]int, N)
		currentY = make([]int, N)
		curr := m[h.y][h.x]
		ptr := 0
		for {
			count++
			// fmt.Printf("%v\n", count)
			// fmt.Printf("\033[1A\033[K")
			// printMap()
			// realTimePrintMap()
			if curr != nil {
				eNode := exhaustNode(curr)
				if eNode != nil {
					nodeStack.Push(eNode)
				}

				if isGG() {
					fmt.Println("took ", count, " steps")
					return nil
				}

				if ptr < len(headsToCheck) {
					nextHead := headsToCheck[ptr]
					curr = m[nextHead.y][nextHead.x]
					ptr++
					continue
				}
			}

			disableNode(nodeStack.Peek().(*Node))
			prev := nodeStack.Peek().(*Node).prev

			if prev != nil && slices.Index(headsToCheck, getHead(prev).Points) != len(headsToCheck)-1 {
				head := getHead(prev)
				nodeStack.Pop()
				nodeStack.Push(prev)

				ptr = slices.Index(headsToCheck, head.Points) + 1
				curr = m[headsToCheck[ptr].y][headsToCheck[ptr].x]
				continue
			}

			head := nodeStack.Pop().(*Node)
			disableNode(head)

			if nodeStack.Len() == 0 {
				break
			}
			curr = nil
		}
		cleanMap()
	}

	fmt.Println("took ", count, " steps")
	return fmt.Errorf("no solution")
}

func printStack() {
	stackCopy := *nodeStack
	for range stackCopy.Len() {
		fmt.Println(stackCopy.Pop().(*Node).Points)
	}

	for i := 0; i < stackCopy.Len(); i++ {
		fmt.Printf("\033[1A\033[K")
	}
}

// we should've need this if the logic is correct, but i suck
// its not really computationally heavy tho so it should be fine
func cleanMap() {
	for _, row := range m {
		for _, node := range row {
			node.active = false
		}
	}
}

func disableNode(node *Node) {
	//just to be safe
	if node.active {
		currentX[node.x]--
		currentY[node.y]--
	}
	node.active = false
}

func getHead(node *Node) *Node {
	if node.prev == nil {
		return node
	}
	return getHead(node.prev)
}

func exhaustNode(node *Node) *Node {
	boom := isBoom(node.x, node.y)
	if boom {
		return node.prev
	} else if !node.active {
		node.active = true
		currentY[node.y]++
		currentX[node.x]++
	}
	if node.next == nil {
		return node
	}
	return exhaustNode(node.next)
}

func isGG() bool {
	for i := range currentX {
		if currentX[i] != targetX[i] || currentY[i] != targetY[i] {
			return false
		}
	}

	return true
}

func isBoom(indexX, indexY int) bool {
	if currentX[indexX] == targetX[indexX] || currentY[indexY] == targetY[indexY] {
		return true
	}
	return false
}

func realTimePrintMap() {
	printMap()

	for i := 0; i < len(targetX)*2+2; i++ {
		fmt.Printf("\033[1A\033[K")
	}
}

func printMap() {
	//println color
	red := "\033[31m"
	reset := "\033[0m"
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
				if node.active {
					fmt.Print(red)
				}
				if node.prev != nil {
					if node.prev.y == node.y {
						fmt.Print(" ━ ")
					} else if node.prev.x == node.x {
						fmt.Print(" | ")
					}
				} else {
					fmt.Print(" o ")
				}
				if node.active {
					fmt.Print(reset)
				}
				fmt.Print("┋")
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

func initData() {
	pts := strings.Split(raw, ";")
	rawNodes := strings.Split(pts[0], ",")

	var prev Points
	for _, rawNode := range rawNodes {
		nodePts := strings.Split(rawNode, "-")
		y, _ := strconv.Atoi(nodePts[1])
		x, _ := strconv.Atoi(nodePts[2])

		if nodePts[0] == "s" {
			heads = append(heads, Points{y: y, x: x})
		} else {
			relations = append(relations, r{this: prev, to: Points{y: y, x: x}})
		}
		prev = Points{y: y, x: x}
	}

	xtargets := strings.Split(pts[1], ",")
	for _, xtarget := range xtargets {
		x, _ := strconv.Atoi(xtarget)
		targetX = append(targetX, x)
	}
	ytargets := strings.Split(pts[2], ",")
	for _, ytarget := range ytargets {
		y, _ := strconv.Atoi(ytarget)
		targetY = append(targetY, y)
	}

	N = len(targetX)
}
