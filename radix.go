package typeahead

import (
	"fmt"
	"strings"
)

// func main() {
//         root := NewTrieNode("^")
//
//         m := make(map[string]int)
//         words := []string{"alexas", "alexander", "alexanders", "alexandreid", "alexandra", "alexandrian", "alexandrianism", "alexandrine", "alexandrina", "alexandrite", "alexic", "alexia", "alexian", "alexin", "alexinic", "alexipyretic", "alexipharmic", "alexipharmical", "alexipharmacum", "alexipharmacon", "alexis", "alexiteric", "alexiterical", "alexius"}
//         for _, w := range words {
//                 root.Add(w)
//
//                 if _, found := m[w]; found {
//                         m[w]++
//                 } else {
//                         m[w] = 0
//                 }
//         }
//
//         root.Print(0)
//         fmt.Println(len(m))
//
//         fmt.Println("")
//         fmt.Println("search results")
//         result := root.Search("alex")
//         fmt.Println("ttoal", len(result))
//         fmt.Println(result)
// }

type TrieNode struct {
	key      string
	endword  bool
	count    int
	pos      int
	children []*TrieNode
}

func NewTrieNode(key string) *TrieNode {
	return &TrieNode{
		key:      key,
		endword:  true, // By default, it should all be true, except when branching.
		count:    1,
		children: make([]*TrieNode, 0),
	}
}
func (n *TrieNode) String() string {
	if n.endword {
		return fmt.Sprintf("%s:%d$", n.key, n.count)
	}
	return fmt.Sprintf("%s:%d", n.key, n.count)

}

// Print prints the whole tree.
func (n *TrieNode) Print(i int) {
	for _, child := range n.children {
		fmt.Printf("%*s %s\n", i*2, "", child)
		child.Print(i + 1)
	}
}

func (n *TrieNode) IsLeaf() bool {
	return len(n.children) == 0
}
func (n *TrieNode) Add(key string) {

	// We need to check if there is an existing key with similar prefix in order to branch it out.

	node := n
	for {
		var child *TrieNode
		var i int

		// Indicates that there are no edges yet, add the first one.
		if node.IsLeaf() {
			// fmt.Println("condition leaf")
			node.children = append(node.children, NewTrieNode(key))
			return
		}
		for _, child = range node.children {
			// We need to find out the index of the prefix that matches.
			i = matchPrefix(child.key, key)
			if i == -1 {
				// We continue if the current child does not have a matching prefix.
				// If this continue until the end of the iteration, we would be left
				// with the last node, and i will be -1.
				continue
			}
			// Break if we found one with a matching prefix. Remember, we can only have
			// one node with a matching prefix.
			break
		}
		if i == -1 {
			// fmt.Println("condition no match")
			// No children found, add one and return.
			node.children = append(node.children, NewTrieNode(key))
			return
		}
		// We already have an exact match, update the count and return.
		if child.key == key {
			child.count++
			// child.endword = true
			break
		}
		// Set the node to be equal the current child with the given prefix.
		// Here there are several conditions to check - can we iterate the child further?

		if child.key == key[:i+1] {
			// fmt.Println("condition 1", child.key, key, key[:i+1], child.key[:i+1])
			node = child
			node.count++

			// Also, update the key
			key = key[i+1:]
			continue
		}
		// child.key is longer than key, which means key is a substring of child.key
		if child.key[:i+1] == key {
			// fmt.Println("condition 2", child.key, key, key[:i+1], child.key[:i+1])
			oldKey := child.key
			child.key = oldKey[:i+1]
			child.count++
			// child.endword = true

			// This must be endword too
			child.children = append(child.children, NewTrieNode(oldKey[i+1:]))
			break
		}
		// E.g. john and jane. We know the first 'j' is the prefix, and john is already in the trie.
		// So we first create a new copy of john, and the key.
		if child.key[:i+1] == key[:i+1] {
			// fmt.Println("condition 3", child.key, key, key[:i+1], child.key[:i+1])
			// Create a copy of the key.
			oldKey := child.key

			// Create a copy of john.
			var nodecpy TrieNode
			nodecpy = *child

			// Then update the key by taking the suffix 'ohn'.
			nodecpy.key = nodecpy.key[i+1:]

			// Override the old node with the prefix 'j'.
			*child = *NewTrieNode(oldKey[:i+1])
			child.endword = false // This is a split node, so it should be false.

			// Append the 'ohn' to the new node.
			child.children = append(child.children, &nodecpy)

			// Now, append 'ane' from 'jane'
			child.children = append(child.children, NewTrieNode(key[i+1:]))
			break
		}
	}
}

func (n *TrieNode) Search(key string) []string {
	var matches int
	var str strings.Builder
	table := make(map[string]struct{})
	node := n
outer:
	for {
		var i int
		var child *TrieNode
		if node.IsLeaf() && node.key != key {
			break outer
		}
		for _, child = range node.children {
			if i = matchPrefix(child.key, key); i == -1 {
				continue
			}
			break
		}
		if i == -1 {
			break outer
		}
		if strings.EqualFold(key, child.key) {
			str.WriteString(key)
			s := str.String()

			var head *TrieNode
			queue := []*TrieNode{child}
			for len(queue) > 0 {
				head, queue = queue[0], queue[1:]
				if len(head.children) == 0 {
					str.Reset()
					str.WriteString(s)
					str.WriteString(head.key[i+1:])
					// result = append(result, str.String())
					table[str.String()] = struct{}{}
				}
				for _, c := range head.children {
					c.key = head.key + c.key
					// Important to include those with endword too.
					if c.endword {
						str.Reset()
						str.WriteString(s)
						str.WriteString(c.key[i+1:])
						// result = append(result, str.String())
						table[str.String()] = struct{}{}
					}
					queue = append(queue, c)
				}
			}
			break outer
		}
		if strings.Contains(key, child.key) {
			node = child
			matches += len(child.key)
			str.WriteString(child.key)
			key = key[i+1:]
		} else {
			break outer
		}
	}
	result := make([]string, len(table))
	var i int
	for r := range table {
		result[i] = r
		i++
	}
	return result
}

// matchPrefix will return -1 if the prefix does not match.
func matchPrefix(s, t string) int {
	// When either one has len zero, it would not match.
	if len(s) == 0 || len(t) == 0 {
		return -1
	}
	// If the first character does not match, it would be -1.
	if s[0] != t[0] {
		return -1
	}

	// We should only compare up to the shortest string index to avoid panic.
	n := len(s)
	// If s is longer than t, take t instead.
	if len(s) > len(t) {
		n = len(t)
	}
	var j int
	for i := 0; i < n; i++ {
		if s[i] == t[i] {
			j = i
			continue
		}
		return j
	}
	return j
}
