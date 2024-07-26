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

const raw = "start 2 0,straight 1 0,end 0 0,start 0 1,straight 0 2,straight 0 3,straight 0 4,straight 0 5,straight 0 6,straight 0 7,straight 0 8,end 0 9,start 4 1,straight 3 1,straight 2 1,end 1 1,start 2 2,end 1 2,start 1 3,straight 1 4,straight 1 5,straight 1 6,end 1 7,start 1 9,end 1 8,start 2 3,straight 2 4,end 2 5,start 2 6,end 2 7,start 2 8,end 3 8,start 2 9,straight 3 9,straight 4 9,straight 5 9,straight 6 9,straight 7 9,straight 8 9,end 9 9,start 3 0,straight 4 0,straight 5 0,end 6 0,start 4 2,end 3 2,start 3 3,straight 3 4,straight 3 5,straight 3 6,end 3 7,start 9 3,straight 8 3,straight 7 3,straight 6 3,straight 5 3,end 4 3,start 6 4,straight 5 4,end 4 4,start 4 5,straight 4 6,straight 4 7,end 4 8,start 5 2,end 5 1,start 6 5,end 5 5,start 8 6,straight 7 6,straight 6 6,end 5 6,start 5 7,straight 6 7,straight 7 7,end 8 7,start 5 8,straight 6 8,end 7 8,start 7 1,end 6 1,start 7 2,end 6 2,start 8 0,end 7 0,start 7 4,end 8 4,start 7 5,end 8 5,start 8 2,end 8 1,start 9 8,end 8 8,start 9 2,straight 9 1,end 9 0,start 9 4,straight 9 5,straight 9 6,end 9 7;9,5,5,5,5,4,3,6,6,1;1,9,6,5,3,4,7,3,3,8;"

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
	for i, h := range heads {
		headsToCheck := []Points{}
		headsToCheck = append(headsToCheck, heads[:i]...)
		headsToCheck = append(headsToCheck, heads[i+1:]...)
		currentX = make([]int, N)
		currentY = make([]int, N)
		curr := m[h.y][h.x]
		ptr := 0
		for {
			eNode := exhaustNode(curr)
			if eNode != nil {
				nodeStack.Push(eNode)
			}

			if isGG() {
				return nil
			}

			if ptr >= len(headsToCheck) {
				disableNode(nodeStack.Peek().(*Node))
				prev := nodeStack.Peek().(*Node).prev
				if prev != nil {
					nodeStack.Pop()
					nodeStack.Push(prev)

					head := getHead(prev)
					ptr = slices.Index(headsToCheck, head.Points) + 1
					if !(ptr >= len(headsToCheck)) {
						curr = m[headsToCheck[ptr].y][headsToCheck[ptr].x]
						continue
					}
				}

				head := nodeStack.Pop().(*Node)
				disableNode(head)

				if nodeStack.Len() == 0 {
					break
				}

				if !(ptr >= len(headsToCheck)) {
					ptr = slices.Index(headsToCheck, head.Points) + 1
					curr = m[headsToCheck[ptr].y][headsToCheck[ptr].x]
					continue
				}

				// disableNode(nodeStack.Peek().(*Node))
				// prev = nodeStack.Peek().(*Node).prev
				// if prev != nil {
				// 	nodeStack.Pop()
				// 	nodeStack.Push(prev)

				// 	head := getHead(prev)
				// 	ptr = slices.Index(headsToCheck, head.Points) + 1
				// 	curr = m[headsToCheck[ptr].y][headsToCheck[ptr].x]
				// 	continue
				// }

				// head = nodeStack.Pop().(*Node)

				// if nodeStack.Len() == 0 {
				// 	break
				// }

				// ptr = slices.Index(headsToCheck, head.Points) + 1
				// curr = m[headsToCheck[ptr].y][headsToCheck[ptr].x]
				// continue

				for ptr >= len(headsToCheck) {
					disableNode(nodeStack.Peek().(*Node))
					prev = nodeStack.Peek().(*Node).prev
					if prev != nil {
						nodeStack.Pop()
						nodeStack.Push(prev)

						head := getHead(prev)
						ptr = slices.Index(headsToCheck, head.Points) + 1
						continue
					}

					head = nodeStack.Pop().(*Node)
					disableNode(head)

					if nodeStack.Len() == 0 {
						break
					}

					ptr = slices.Index(headsToCheck, head.Points) + 1
				}
				if nodeStack.Len() == 0 {
					break
				}

				curr = m[headsToCheck[ptr].y][headsToCheck[ptr].x]
				continue
			}

			//get next head
			nextHead := headsToCheck[ptr]
			curr = m[nextHead.y][nextHead.x]
			ptr++
		}
	}

	return fmt.Errorf("no solution")
}

func printStack() {
	stackCopy := *nodeStack
	for range stackCopy.Len() {
		fmt.Println(stackCopy.Pop().(*Node).Points)
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
	} else {
		node.active = true
		currentX[node.x]++
		currentY[node.y]++
		if node.next == nil {
			return node
		}

		return exhaustNode(node.next)
	}
}

func isGG() bool {
	xScore := make([]int, len(targetX))
	yScore := make([]int, len(targetY))
	for _, row := range m {
		for _, node := range row {
			if node != nil && node.active {
				xScore[node.x]++
				yScore[node.y]++
			}
		}
	}

	for i := range currentX {
		if xScore[i] != targetX[i] || yScore[i] != targetY[i] {
			return false
		}
	}

	return true
}

func isBoom(indexX, indexY int) bool {
	if currentX[indexX]+1 > targetX[indexX] || currentY[indexY]+1 > targetY[indexY] {
		return true
	}

	return false
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
		nodePts := strings.Split(rawNode, " ")
		y, _ := strconv.Atoi(nodePts[1])
		x, _ := strconv.Atoi(nodePts[2])

		if nodePts[0] == "start" {
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
