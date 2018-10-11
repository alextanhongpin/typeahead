package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"
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
		node := root.edges[pos].node
		insert(&node, key[p:], value)
		root.edges[pos].node = node
		return
	}
	split(root, key, p, pos)
}

func split(root *Node, key []byte, p, pos int) {
	var rem int
	edge := root.edges[pos]
	for k, v := range root.edges {
		if bytes.Equal(v.key, edge.key) {
			rem = k
			break
		}
	}
	root.edges = append(root.edges[:rem], root.edges[rem+1:]...)
	prefix, left, right := edge.key[:p], edge.key[p:], key[p:]

	newEdge := NewEdge(prefix, nil)
	newEdge.count += edge.count

	edge.key = left

	insert(&(newEdge.node), right, nil)
	newEdge.node.edges = append(newEdge.node.edges, edge)
	root.edges = append(root.edges, newEdge)
}

func find(root *Node, in []byte) map[string]bool {
	if len(in) == 0 {
		return nil
	}

	result := make(map[string]bool)
	queue := make([][]byte, 0)
	queue = append(queue, in)
	var key []byte
	for len(queue) > 0 {
		key, queue = queue[0], queue[1:]
		if valid, found := result[string(key)]; found || valid {
			continue
		} else {
			result[string(key)] = found
		}
		var foundElements int

		var foundKey []byte
		foundKey = make([]byte, len(key))
		copy(foundKey, key)
		traverseNode := root

		for traverseNode != nil && !traverseNode.IsLeaf() && foundElements < len(key) {
			var nextEdge *Edge
			for _, edge := range traverseNode.edges {
				if bytes.HasPrefix(key[foundElements:], edge.key) {
					nextEdge = &edge
					break
				}
			}
			if nextEdge == nil {
				traverseNode = nil
				break
			}
			foundElements += len(nextEdge.key)
			traverseNode = &(nextEdge.node)
		}
		if traverseNode == nil || traverseNode.IsLeaf() {
			continue
		}
		for _, edge := range traverseNode.edges {
			func(edge Edge) {
				out := append(foundKey, edge.key...)
				if _, found := result[string(out)]; !found {
					queue = append([][]byte{out}, queue...)
				}
				if edge.node.IsLeaf() {
					result[string(out)] = true
				}
			}(edge)
		}
	}
	return result
}

func main() {
	var (
		cpuprofile  = flag.String("cpu", "", "The file to output the cpu profiling, e.g. cpu.out")
		memprofile  = flag.String("mem", "", "The file to output the memory profiling, e.g. mem.out")
		interactive = flag.Bool("i", false, "Whether to allow interactive mode or not")
	)
	flag.Parse()
	if *cpuprofile != "" {
		cpufile, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(cpufile)
		defer pprof.StopCPUProfile()
	}

	root := NewNode()
	f, err := os.Open("/usr/share/dict/words")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var count int
	for scanner.Scan() {
		count++
		b := bytes.ToLower(scanner.Bytes())
		/* if bytes.HasPrefix(b, []byte("alex")) { */
		/*         fmt.Println(string(b)) */
		/*         iii++ */
		/* } */
		insert(&root, b, nil)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Dictionary has", count, "results")

	// root.Print(0)
	if *memprofile != "" {
		memfile, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(memfile)
		memfile.Close()
	}

	if *interactive {
		fmt.Println("Enter a search keyword:")
		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
			b := bytes.TrimSpace(reader.Bytes())
			if len(b) == 0 {
				continue
			}
			fmt.Printf("searching for %s:\n", b)
			start := time.Now()
			result := find(&root, b)
			var count int
			fmt.Printf("found %d results in %s\n", len(result), time.Since(start))
			for r, endWord := range result {
				if endWord {
					count++
					fmt.Println(r)
				}
			}
			fmt.Printf("found %d results in %s\n", count, time.Since(start))
			fmt.Println()
		}
		if err := reader.Err(); err != nil {
			log.Fatal(err)
		}
	}
}
