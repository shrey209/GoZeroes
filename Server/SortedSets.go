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
	span     []int
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
		span:     make([]int, level),
	}
}

func NewSkipList(maxLevel int) *SkipList {
	head := NewNode(-1, "head", maxLevel)
	return &SkipList{
		maxLevel: maxLevel,
		curLevel: 1,
		head:     head,
		members:  make(map[string]*Node),
		tail:     head,
	}
}

/*
this method randomly generatess the level for a node
*/
func (sl *SkipList) generateMaxLevel() int {
	level := 1
	for rand.Intn(2) == 1 && level < sl.maxLevel {
		level++
	}
	return level
}

/*
For inserting a node
*/
func (sl *SkipList) Insert(value int, key string) {
	if _, exists := sl.members[key]; exists {
		sl.Delete(key)
	}

	// Create a new node with a randomly generated level.
	n1 := NewNode(value, key, sl.generateMaxLevel())
	update := make([]*Node, sl.maxLevel)
	rank := make([]int, sl.maxLevel)
	current := sl.head

	// Step 1: Determine update array and rank array
	for i := sl.curLevel - 1; i >= 0; i-- {
		rank[i] = 0 // Initialize rank for each level
		for current.forward[i] != nil && current.forward[i].value < value {
			rank[i] += current.span[i] // Add the span as we traverse
			current = current.forward[i]
		}
		update[i] = current
	}

	// Step 2: Adjust the current level of the skip list
	if len(n1.forward) > sl.curLevel {
		for i := sl.curLevel; i < len(n1.forward); i++ {
			update[i] = sl.head
			rank[i] = 0
			sl.head.span[i] = len(sl.members) + 1 // Update head span for higher levels
		}
		sl.curLevel = len(n1.forward)
	}

	// Step 3: Update spans and forward pointers
	for i := 0; i < len(n1.forward); i++ {
		// Update forward pointers
		n1.forward[i] = update[i].forward[i]
		update[i].forward[i] = n1

		// Calculate spans for the new node
		if n1.forward[i] != nil {
			n1.span[i] = update[i].span[i] - (rank[0] - rank[i]) + 1
		} else {
			n1.span[i] = 1 // Last node at this level
		}

		// Update spans for the update array
		if update[i].forward[i] != nil {
			update[i].span[i] = rank[0] - rank[i] + 1
		} else {
			update[i].span[i] = 1
		}
	}

	// Step 4: Handle backward pointer
	n1.backward = update[0]
	if n1.forward[0] != nil {
		n1.forward[0].backward = n1
	} else {
		sl.tail = n1
	}

	// Add the node to the members map
	sl.members[key] = n1
}

/*
for deleteing
*/
func (sl *SkipList) Delete(key string) {
	node, exists := sl.members[key]
	if !exists {
		fmt.Println("No such key exists")
		return
	}

	value := node.value
	update := make([]*Node, sl.maxLevel)
	current := sl.head

	// Find update array and calculate spans
	for i := sl.curLevel - 1; i >= 0; i-- {
		for current.forward[i] != nil && (current.forward[i].value < value || (current.forward[i].value == value && current.forward[i].member != key)) {
			current = current.forward[i]
		}
		update[i] = current
	}

	target := current.forward[0]
	if target != nil && target.value == value && target.member == key {
		// Update pointers and spans
		for i := 0; i < len(target.forward); i++ {
			if update[i].forward[i] == target {
				// Update span
				update[i].span[i] += target.span[i] - 1
				// Update forward pointer
				update[i].forward[i] = target.forward[i]
			}
		}

		// Update backward pointer of the next node
		if target.forward[0] != nil {
			target.forward[0].backward = target.backward
		} else {
			sl.tail = target.backward
		}

		// Reduce the skip list level if necessary
		for sl.curLevel > 0 && sl.head.forward[sl.curLevel-1] == nil {
			sl.curLevel--
		}

		// Remove from members map
		delete(sl.members, key)
	} else {
		fmt.Println("Node to delete not found")
	}
}

/*to display the skip list
 */
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
	fmt.Println("Display Forward:")
	cur := sl.head.forward[0]
	for cur != nil {
		fmt.Printf("%d ", cur.value)
		cur = cur.forward[0]
	}
	fmt.Println()
}

func (sl *SkipList) DisplayReverse() {
	fmt.Println("Display Reverse:")
	cur := sl.tail
	for cur != nil {
		fmt.Printf("%d ", cur.value)
		cur = cur.backward
	}
	fmt.Println()
}

func (sl *SkipList) DisplayForwardN(n int) {
	fmt.Println("Display Forward N:")
	cur := sl.head.forward[0]
	for cur != nil && n > 0 {
		fmt.Printf("%d ", cur.value)
		cur = cur.forward[0]
		n--
	}
	fmt.Println()
}

func (sl *SkipList) DisplayReverseN(n int) {
	fmt.Println("Display Reverse N:")
	cur := sl.tail
	for cur != nil && n > 0 {
		fmt.Printf("%d ", cur.value)
		cur = cur.backward
		n--
	}
	fmt.Println()
}

func (sl *SkipList) GetScore(key string) int {
	node, exists := sl.members[key]
	if exists {
		return node.value
	}
	return -1
}

func (sl *SkipList) getRank(key string) int {
	node, exists := sl.members[key]
	if !exists {
		return -1 // Key not found
	}

	value := node.value
	rank := 0
	current := sl.head

	// Traverse through the levels
	for i := sl.curLevel - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].value < value {
			rank += current.span[i] // Add span to rank
			current = current.forward[i]
		}
	}

	// Add 1 to rank if the target node exists at the current position
	if current.forward[0] != nil && current.forward[0].value == value && current.forward[0].member == key {
		rank += 1
	}

	return rank
}

func main() {
	skipList := NewSkipList(5)

	skipList.Insert(1, "one")
	skipList.Insert(2, "two")
	skipList.Insert(7, "seven")
	skipList.Insert(7, "newseven")
	skipList.Insert(4, "four")
	skipList.Insert(5, "five")

	fmt.Println("Before deletion:")
	skipList.DisplayAll()

	// skipList.Delete("one")
	// skipList.Delete("seven")

	fmt.Println(skipList.getRank("seven"))
	fmt.Println(skipList.getRank("newseven"))

	fmt.Println("After deletion:")
	skipList.DisplayAll()
	skipList.DisplayReverse()
	skipList.DisplayForward()
}
