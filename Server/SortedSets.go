package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Node struct {
	value   int
	forward []*Node
}

type SkipList struct {
	Members  map[string]*Node
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
		Members:  make(map[string]*Node),
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

func (sl *SkipList) Insert(value int, key string) {

	_, exists := sl.Members[key]
	if exists {
		sl.Delete(key)
	}

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
	sl.Members[key] = newNode
}

func (sl *SkipList) Delete(key string) {
	// Check if the key exists
	nodeToDelete, exists := sl.Members[key]
	if !exists {
		log.Println("Key does not exist")
		return
	}

	// Traverse the skip list to update forward pointers
	current := sl.head
	for i := sl.curLevel - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i] != nodeToDelete {
			current = current.forward[i]
		}

		// If the forward pointer points to the node to delete, update it
		if current.forward[i] == nodeToDelete {
			current.forward[i] = nodeToDelete.forward[i]
		}
	}

	// Adjust current level if the top levels become empty
	for sl.curLevel > 0 && sl.head.forward[sl.curLevel-1] == nil {
		sl.curLevel--
	}

	// Remove the node from the Members map
	delete(sl.Members, key)
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

func main() {
	rand.Seed(time.Now().UnixNano())

	skipList := NewSkipList(5)
	skipList.Insert(100, "one")
	skipList.Insert(200, "two")
	skipList.Insert(10, "three")
	skipList.Insert(4, "four")
	skipList.Insert(5, "five")

	fmt.Println("Before deletion:")
	skipList.Display()

	skipList.Delete("three")
	skipList.Delete("four")

	fmt.Println("After deletion:")
	skipList.Display()

	for key, value := range skipList.Members {
		fmt.Printf("key: %s, value: %d\n", key, value.value)
	}

}
