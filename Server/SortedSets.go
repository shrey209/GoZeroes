package main

import (
	"fmt"
	"math/rand"
)

type Node struct {
	value   int
	member  string
	forward []*Node
}

type SkipList struct {
	maxLevel int
	curLevel int
	head     *Node
	members  map[string]*Node
}

func NewNode(value int, member string, level int) *Node {
	return &Node{
		value:   value,
		member:  member,
		forward: make([]*Node, level),
	}
}

func NewSkipList(maxLevel int) *SkipList {
	return &SkipList{
		maxLevel: maxLevel,
		curLevel: 0,
		head:     NewNode(-1, "head", maxLevel),
		members:  make(map[string]*Node),
	}
}

func (sl *SkipList) generateMaxLevel() int {
	level := 1
	for rand.Intn(2) == 1 && level < sl.maxLevel {
		level++
	}
	return level
}

func (sl *SkipList) Insert(value int, key string) {
	if _, exists := sl.members[key]; exists {
		sl.Delete(key)
	}

	n1 := NewNode(value, key, sl.generateMaxLevel())
	update := make([]*Node, sl.maxLevel)
	current := sl.head

	for i := sl.curLevel - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].value < value {
			current = current.forward[i]
		}
		update[i] = current
	}

	if len(n1.forward) > sl.curLevel {
		for i := sl.curLevel; i < len(n1.forward); i++ {
			update[i] = sl.head
		}
		sl.curLevel = len(n1.forward)
	}

	for i := 0; i < len(n1.forward); i++ {
		if update[i].forward[i] != nil {
			n1.forward[i] = update[i].forward[i]
		}
		update[i].forward[i] = n1
	}

	sl.members[key] = n1
}

func (sl *SkipList) Delete(key string) {
	node, exists := sl.members[key]
	if !exists {
		fmt.Println("No such key exists")
		return
	}

	value := node.value
	update := make([]*Node, sl.maxLevel)
	current := sl.head

	for i := sl.curLevel - 1; i >= 0; i-- {
		for current.forward[i] != nil && (current.forward[i].value < value || (current.forward[i].value == value && current.forward[i].member != key)) {
			current = current.forward[i]
		}
		update[i] = current
	}

	target := current.forward[0]
	if target != nil && target.value == value && target.member == key {
		for i := 0; i < len(target.forward); i++ {
			if update[i].forward[i] == target {
				update[i].forward[i] = target.forward[i]
			}
		}
		for sl.curLevel > 0 && sl.head.forward[sl.curLevel-1] == nil {
			sl.curLevel--
		}
		delete(sl.members, key)
	}
}

func (sl *SkipList) Display() {
	fmt.Println("Skip List:")
	for lvl := sl.curLevel - 1; lvl >= 0; lvl-- {
		current := sl.head
		fmt.Printf("Level %d: ", lvl)
		for current.forward[lvl] != nil {
			current = current.forward[lvl]
			fmt.Printf("%d ", current.value)
		}
		fmt.Println()
	}
}

func main() {
	skipList := NewSkipList(5)

	skipList.Insert(1, "one")
	skipList.Insert(2, "two")
	skipList.Insert(7, "seven")
	skipList.Insert(4, "four")
	skipList.Insert(5, "five")

	fmt.Println("Before deletion:")
	skipList.Display()

	skipList.Delete("one")
	skipList.Delete("seven")

	fmt.Println("After deletion:")
	skipList.Display()
}
