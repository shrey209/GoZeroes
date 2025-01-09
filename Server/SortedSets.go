package main

import (
	"fmt"
	"math/rand"
)

type Node struct {
	value   int
	forward []*Node
}

type SkipList struct {
	members  map[string]*Node
	maxLevel int
	curLevel int
	head     *Node
}

func NewNode(value, level int) *Node {
	return &Node{
		value:   value,
		forward: make([]*Node, level),
	}
}

func NewSkipList(maxLevel int) *SkipList {
	head := NewNode(-1, maxLevel)
	return &SkipList{
		maxLevel: maxLevel,
		curLevel: 0,
		head:     head,
	}
}

func (sl *SkipList) generateMaxLevel() int {
	level := 1
	for rand.Float32() < 0.5 && level < sl.maxLevel {
		level++
	}
	return level
}

func (sl *SkipList) Insert(value int) {
	update := make([]*Node, sl.maxLevel)
	current := sl.head

	for i := sl.curLevel - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].value < value {
			current = current.forward[i]
		}
		update[i] = current
	}

	level := sl.generateMaxLevel()

	if level > sl.curLevel {
		for i := sl.curLevel; i < level; i++ {
			update[i] = sl.head
		}
		sl.curLevel = level
	}

	newNode := NewNode(value, level)
	for i := 0; i < level; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
}

func (sl *SkipList) Delete(value int) {
	update := make([]*Node, sl.maxLevel)
	current := sl.head

	for i := sl.curLevel - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].value < value {
			current = current.forward[i]
		}
		update[i] = current
	}

	target := current.forward[0]
	if target != nil && target.value == value {

		for i := 0; i < sl.curLevel; i++ {
			if update[i].forward[i] == target {
				update[i].forward[i] = target.forward[i]
			}
		}

		for sl.curLevel > 0 && sl.head.forward[sl.curLevel-1] == nil {
			sl.curLevel--
		}
	}
}

func (sl *SkipList) Display() {
	fmt.Println("Skip List:")
	for i := sl.curLevel - 1; i >= 0; i-- {
		fmt.Printf("Level %d: ", i)
		current := sl.head
		for current.forward[i] != nil {
			fmt.Printf("%d ", current.forward[i].value)
			current = current.forward[i]
		}
		fmt.Println()
	}
}

//for testing
// func main() {
// 	rand.Seed(time.Now().UnixNano())

// 	skipList := NewSkipList(5)
// 	skipList.Insert(100)
// 	skipList.Insert(200)
// 	skipList.Insert(10)
// 	skipList.Insert(4)
// 	skipList.Insert(5)

// 	fmt.Println("Before deletion:")
// 	skipList.Display()

// 	skipList.Delete(10)
// 	skipList.Delete(4)

// 	fmt.Println("After deletion:")
// 	skipList.Display()
// }
