package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"
)

type Node struct {
	edges []Edge
}

func NewNode() Node {
	return Node{}
}

func (n Node) IsLeaf() bool {
	return len(n.edges) == 0
}

func (n Node) Print(depth int) {
	for _, edge := range n.edges {
		key, count := edge.key, edge.count
		fmt.Printf("%s %s:%d\n", strings.Repeat(" ", depth*2), key, count)
		edge.node.Print(depth + 1)
	}
}

type Edge struct {
	count int
	key   []byte
	value interface{}
	node  Node
}

func NewEdge(key []byte, value interface{}) Edge {
	return Edge{
		key:   key,
		value: value,
		count: 1,
		node:  NewNode(),
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func sharedPrefix(s, t []byte) int {
	minLen := min(len(s), len(t))
	for i := 0; i < minLen; i++ {
		if s[i] != t[i] {
			return i
		}
	}
	return minLen
}

func insert(root *Node, key []byte, value interface{}) {
	if root == nil || len(key) == 0 {
		return
	}
	var p, pos int
	for i, child := range root.edges {
		if len(child.key) == 0 {
			continue
		}
		if child.key[0] != key[0] {
			continue
		}
		// This is expensive to compute, hence the checking above will
		// reduce the computation.
		p = sharedPrefix(child.key, key)
		if p == 0 {
			continue
		}
		pos = i
		break
	}
	if p == 0 {
		root.edges = append(root.edges, NewEdge(key, value))
		return
	}
	nextEdge := root.edges[pos]
	if bytes.Equal(nextEdge.key, key) {
		root.edges[pos].count++
		return
	}
	if len(nextEdge.key) == p {
		// if bytes.HasPrefix(key, nextEdge.key) {
		// fmt.Println(string(key), string(nextEdge.key), p)
		root.edges[pos].count++
		insert(&(root.edges[pos].node), key[p:], value)
		return
	}
	split(key, p, pos, root)
}

func split(key []byte, pos, p int, root *Node) {
	var rem int
	edge := root.edges[p]
	for k, v := range root.edges {
		if bytes.Equal(v.key, edge.key) {
			rem = k
			break
		}
	}
	root.edges = append(root.edges[:rem], root.edges[rem+1:]...)
	prefix, left, right := edge.key[:pos], edge.key[pos:], key[pos:]

	newEdge := NewEdge(prefix, "")
	newEdge.count += edge.count

	edge.key = left

	insert(&(newEdge.node), right, "")
	// if newEdge.node == nil {
	//         newEdge.node = make([]Edge, 0)
	// }
	newEdge.node.edges = append(newEdge.node.edges, edge)
	root.edges = append(root.edges, newEdge)
}

func find(root *Node, key []byte) {
	var foundElements int
	traverseNode := root
	for traverseNode != nil && !traverseNode.IsLeaf() && foundElements < len(key) {
		var nextEdge *Edge
		for _, edge := range traverseNode.edges {
			if len(edge.key) == 0 || edge.key[0] != key[foundElements] {
				continue
			}
			nextEdge = &edge
			break
		}
		if nextEdge == nil {
			break
		}
		foundElements += len(nextEdge.key)
		traverseNode = &(nextEdge.node)
	}
	if traverseNode == nil || traverseNode.IsLeaf() {
		return
	}
	for _, edge := range traverseNode.edges {
		func(in []byte) {
			out := make([]byte, len(key))
			copy(out, key)
			out = append(out, in...)
			fmt.Printf("%s %s\n", out, in)
			find(root, out)
		}(edge.key)
	}
}

// TODO: lowercase the characters.
func main() {
	cpufile, err := os.Create("cpuprofile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpufile)
	defer pprof.StopCPUProfile()

	root := NewNode()
	f, err := os.Open("/usr/share/dict/words")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// var i int
	for scanner.Scan() {
		// i++
		// if i > 50 {
		//         break
		// }
		b := scanner.Bytes()
		// fmt.Println(string(b))
		insert(&root, bytes.ToLower(b), nil)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	find(&root, []byte("abac"))
	// root.Print(0)

	memfile, err := os.Create("memprofile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(memfile)
	memfile.Close()

	// reader := bufio.NewScanner(os.Stdin)
	// defer reader.Close()
	// for reader.Scan() {
	//         b := reader.Bytes()
	//         fmt.Printf("searching for %s\n", b)
	//         find(&root, b)
	// }
	// if err := reader.Err(); err != nil {
	//         log.Fatal(err)
	// }
}
