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
	count   int
	key     []byte
	value   interface{}
	node    Node
	endword bool
}

func NewEdge(key []byte, value interface{}) Edge {
	return Edge{
		key:   key,
		value: value,
		count: 1,
		node:  NewNode(),
	}
}

func (e Edge) String() string {
	return string(e.key)
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
		edge := NewEdge(key, value)
		edge.endword = true
		root.edges = append(root.edges, edge)
		return
	}
	nextEdge := root.edges[pos]
	if bytes.Equal(nextEdge.key, key) {
		root.edges[pos].count++
		root.edges[pos].endword = true
		return
	}
	if len(nextEdge.key) == p {
		root.edges[pos].count++
		// root.edges[pos].endword = true
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
	// newEdge.endword = true

	edge.key = left

	insert(&(newEdge.node), right, nil)
	newEdge.node.edges = append(newEdge.node.edges, edge)
	root.edges = append(root.edges, newEdge)
}

func find(root *Node, in []byte) map[string]Edge {
	if len(in) == 0 {
		return nil
	}

	result := make(map[string]Edge)
	queue := make([]Edge, 0)
	queue = append(queue, NewEdge(in, nil))
	for len(queue) > 0 {
		var key Edge
		var foundKey []byte
		key, queue = queue[0], queue[1:]
		// if _, found := result[key.String()]; found {
		//         // result[key.String()] = key.endword
		//         // continue
		// } else {
		//         result[key.String()] = key
		// }
		foundKey = make([]byte, len(key.key))
		copy(foundKey, key.key)
		// if key.endword {
		//         result[key.String()] = true
		// }
		var foundElements int
		traverseNode := root
		for traverseNode != nil && !traverseNode.IsLeaf() && foundElements < len(foundKey) {
			var nextEdge *Edge
			for _, edge := range traverseNode.edges {
				if bytes.HasPrefix(foundKey[foundElements:], edge.key) {
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
			func(key, edge Edge) {
				out := append(foundKey, edge.key...)
				// if _, found := result[string(out)]; !found {
				queue = append([]Edge{NewEdge(out, nil)}, queue...)
				result[string(out)] = edge
				// }

				// if edge.node.IsLeaf() {
				// }
			}(key, edge)
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
	var ii int
	for scanner.Scan() {
		count++
		b := bytes.ToLower(scanner.Bytes())
		if bytes.HasPrefix(b, []byte("paper")) {
			ii++
		}
		insert(&root, b, nil)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Dictionary has", count, "results", ii)
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
			for r, edge := range result {
				if edge.endword {
					fmt.Println(r)
					count++
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
