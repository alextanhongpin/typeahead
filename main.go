package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"
)

type Node struct {
	edges []*Edge
}

func NewNode() *Node {
	return &Node{
		edges: make([]*Edge, 0),
	}
}

func (n *Node) Print(depth int) {
	for _, edge := range n.edges {
		fmt.Println(strings.Repeat(" ", depth*2), edge.key, edge.count)
		if edge.node != nil {
			edge.node.Print(depth + 1)
		}
	}
}

type Edge struct {
	count int
	key   string
	value string
	node  *Node
}

func NewEdge(key, value string) *Edge {
	return &Edge{key: key, value: value, count: 1, node: NewNode()}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func sharedPrefix(s, t string) int {
	minLen := min(len(s), len(t))
	for i := 0; i < minLen; i++ {
		if s[i] != t[i] {
			return i
		}
	}
	return minLen
}

func insert(root *Node, key, value string) {
	var p int
	var nextEdge *Edge
	for _, child := range root.edges {
		p = sharedPrefix(child.key, key)
		if p == 0 {
			continue
		} else {
			nextEdge = child
			break
		}
	}
	// fmt.Println("insert", key, p)
	if nextEdge == nil || p == 0 {
		root.edges = append(root.edges, NewEdge(key, value))
		return
	}
	if nextEdge.key == key {
		nextEdge.count++
		return
	}
	if strings.HasPrefix(key, nextEdge.key) {
		nextEdge.count++
		insert(nextEdge.node, key[p:], value)
		return
	}
	split(key, p, root, nextEdge)
}

func split(key string, pos int, root *Node, edge *Edge) {
	var rem int
	for k, v := range root.edges {
		if v.key == edge.key {
			rem = k
			break
		}
	}
	root.edges = append(root.edges[:rem], root.edges[rem+1:]...)
	prefix, left, right := edge.key[:pos], edge.key[pos:], key[pos:]

	newEdge := NewEdge(prefix, "")
	newEdge.count += edge.count

	edge.key = left

	insert(newEdge.node, right, "")
	newEdge.node.edges = append(newEdge.node.edges, edge)
	root.edges = append(root.edges, newEdge)
}

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
	for scanner.Scan() {
		insert(root, scanner.Text(), "")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	memfile, err := os.Create("memprofile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(memfile)
	memfile.Close()
}
