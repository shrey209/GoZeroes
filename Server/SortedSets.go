package main

import (
	"fmt"
	"math/rand"
)

type Node struct {
	value    int
	member   string
	forward  []*Node
	backward *Node
}

type SkipList struct {
	maxLevel int
	curLevel int
	head     *Node
	members  map[string]*Node
	tail     *Node
}

func NewNode(value int, member string, level int) *Node {
	return &Node{
		value:    value,
		member:   member,
		forward:  make([]*Node, level),
		backward: nil,
	}
}

func NewSkipList(maxLevel int) *SkipList {
	head := NewNode(-1, "head", maxLevel)
	return &SkipList{
		maxLevel: maxLevel,
		curLevel: 0,
		head:     head,
		members:  make(map[string]*Node),
		tail:     head,
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
		n1.forward[i] = update[i].forward[i]
		update[i].forward[i] = n1
	}

	n1.backward = update[0]
	if n1.forward[0] != nil {
		n1.forward[0].backward = n1
	}

	if n1.forward[0] == nil {
		sl.tail = n1
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

		if target.forward[0] != nil {
			target.forward[0].backward = target.backward
		} else {
			sl.tail = target.backward
		}
		for sl.curLevel > 0 && sl.head.forward[sl.curLevel-1] == nil {
			sl.curLevel--
		}
		delete(sl.members, key)
	}
}

func (sl *SkipList) DisplayAll() {
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

func (sl *SkipList) DisplayForward() {
	fmt.Println("Displat forward ->")
	cur := sl.head
	for cur != nil {
		fmt.Printf("%d ", cur.value)
		cur = cur.forward[0]
	}
	fmt.Println()
}

func (sl *SkipList) DisplayReverse() {
	fmt.Println("Reverse ->")
	cur := sl.tail
	for cur != nil {
		fmt.Printf("%d ", cur.value)
		cur = cur.backward
	}
	fmt.Println()
}

func (sl *SkipList) DisplayForwardN(n int) {
	fmt.Println("Display forward ->")
	cur := sl.head
	for cur != nil && n > 0 {
		fmt.Printf("%d ", cur.value)
		cur = cur.forward[0]
		n--
	}
	fmt.Println()
}

func (sl *SkipList) DisplayReverseN(n int) {
	fmt.Println("Reverse ->")
	cur := sl.tail
	for cur != nil && n > 0 {
		fmt.Printf("%d ", cur.value)
		cur = cur.backward
		n--
	}
	fmt.Println()
}

func main() {
	skipList := NewSkipList(5)

	skipList.Insert(1, "one")
	skipList.Insert(2, "two")
	skipList.Insert(7, "seven")
	skipList.Insert(4, "four")
	skipList.Insert(5, "five")

	fmt.Println("Before deletion:")
	skipList.DisplayAll()

	skipList.Delete("one")
	skipList.Delete("seven")

	fmt.Println("After deletion:")
	skipList.DisplayAll()
	skipList.DisplayReverse()
	skipList.DisplayForward()
}
