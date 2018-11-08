package typeahead

package main

import (
	"fmt"
)

// REFERENCES:
// https://www.cs.upc.edu/~ps/downloads/tst/tst.html
// http://hacktalks.blogspot.com/2012/03/implementing-auto-complete-with-ternary.html
// https://www.javaworld.com/article/2075027/java-app-dev/plant-your-data-in-a-ternary-search-tree.html?page=2
// http://www.drdobbs.com/database/ternary-search-trees/184410528?pgno=3
func main() {
	tree := NewTernaryTree()
	tree.Add("hello")
	tree.Add("dobby")
	tree.Add("debby")
	tree.Add("dobbs")
	tree.Add("hells")
	tree.Add("helrs")
	tree.Add("helrs")
	tree.Add("helsinki")
	tree.Add("hobby")
	tree.Add("sell")
	tree.Add("hallo")
	tree.Add("anglo")
	tree.Add("hi")
	tree.Add("car")
	fmt.Println("has hello:", tree.Contains("dobby"))
	result := tree.Search("he")
	fmt.Println(result)
	fmt.Println(tree.Traverse())
	tree.NearSearch("dobbs")
}

type TernaryNode struct {
	char    rune
	left    *TernaryNode
	right   *TernaryNode
	center  *TernaryNode
	endword bool
}

func NewTernaryNode(char rune, endword bool) *TernaryNode {
	return &TernaryNode{
		char:    char,
		endword: endword,
	}
}

type TernaryTree struct {
	root *TernaryNode
}

func NewTernaryTree() *TernaryTree { return &TernaryTree{} }

// Add adds the item to the tree recursively.
func (t *TernaryTree) Add(s string) {
	t.root = t.radd([]rune(s), 0, t.root)
}

func (t *TernaryTree) radd(s []rune, pos int, node *TernaryNode) *TernaryNode {
	if node == nil {
		node = NewTernaryNode(s[pos], false)
	}
	if s[pos] < node.char {
		node.left = t.radd(s, pos, node.left)
	} else if s[pos] > node.char {
		node.right = t.radd(s, pos, node.right)
	} else {
		if pos+1 == len(s) {
			node.endword = true
		} else {
			node.center = t.radd(s, pos+1, node.center)
		}
	}
	return node
}

func (t *TernaryTree) Contains(str string) bool {
	r := []rune(str)
	result := traverse(t.root, r)
	if result == nil {
		return false
	}
	return result.endword
}

// traverse returns the last node that the traversal of the tree returns
func traverse(node *TernaryNode, r []rune) *TernaryNode {
	var pos int
	for node != nil {
		switch {
		case r[pos] < node.char:
			node = node.left
		case r[pos] > node.char:
			node = node.right
		default:
			pos++
			if pos == len(r) {
				return node
			}
			node = node.center
		}
	}
	return node
}

// Search implements an autocomplete for ternary search tree.
func (t *TernaryTree) Search(str string) (result []string) {
	r := []rune(str)
	node := traverse(t.root, r)
	if node == nil {
		return
	}
	t.dfs(node.center, append(r, node.center.char), &result)
	return
}

func (t *TernaryTree) Traverse() (result []string) {
	match := []rune{t.root.char}
	t.dfs(t.root, match, &result)
	return
}

// dfs performs depth-first-search for ternary tree.
func (t *TernaryTree) dfs(node *TernaryNode, match []rune, result *[]string) {
	if node.endword {
		*result = append(*result, string(match))
	}
	if node.center == nil && node.left == nil && node.right == nil {
		return
	}
	if node.center != nil {
		t.dfs(node.center, append(match, node.center.char), result)
	}
	if node.right != nil {
		match[len(match)-1] = node.right.char
		t.dfs(node.right, match, result)
	}
	if node.left != nil {
		match[len(match)-1] = node.left.char
		t.dfs(node.left, match, result)
	}
}

// TODO: Implement nearest neighbour with hamming distance.
