package ahead

import (
	"bytes"
)

// Root represents the root of the radix tree.
type Root struct {
	Node Node
}

// New returns a new tree.
func New() *Root {
	return &Root{
		Node: NewNode(),
	}
}

// Insert adds a key value pair into the tree.
func (r *Root) Insert(key []byte, value interface{}) {
	insert(&(r.Node), key, value)
}

// Find searches for the edge of the node that matches the given prefix.
func (r *Root) Find(key []byte) map[string]Edge {
	return find(&(r.Node), key)
}

func (r *Root) FindRecursive(key []byte) [][]byte {
	return findRecursive(&(r.Node), key)
}

func insert(root *Node, key []byte, value interface{}) {
	if root == nil || len(key) == 0 {
		return
	}
	var p, pos int
	for i, child := range root.Edges {
		if len(child.Key) == 0 || child.Key[0] != key[0] {
			continue
		}
		p = sharedPrefix(child.Key, key)
		pos = i
		break
	}
	if p == 0 {
		edge := NewEdge(key, value)
		edge.Endword = true
		root.Edges = append(root.Edges, edge)
		return
	}
	currKey := root.Edges[pos].Key
	if bytes.Equal(currKey, key) {
		root.Edges[pos].Count++
		root.Edges[pos].Endword = true
		return
	}
	if len(currKey) == p {
		root.Edges[pos].Count++
		node := root.Edges[pos].Node
		insert(&node, key[p:], value)
		root.Edges[pos].Node = node
		return
	}
	split(root, key, p, pos)
}

func split(root *Node, key []byte, p, pos int) {
	var rem int
	edge := root.Edges[pos]
	for k, v := range root.Edges {
		if bytes.Equal(v.Key, edge.Key) {
			rem = k
			break
		}
	}
	// Pop the old edge.
	root.Edges = append(root.Edges[:rem], root.Edges[rem+1:]...)
	prefix, left, right := edge.Key[:p], edge.Key[p:], key[p:]

	newEdge := NewEdge(prefix, nil)
	newEdge.Count += edge.Count

	edge.Key = left

	insert(&(newEdge.Node), right, nil)
	newEdge.Node.Edges = append(newEdge.Node.Edges, edge)
	root.Edges = append(root.Edges, newEdge)
}

func complete(root *Node, orikey []byte) [][]byte {
	key := make([]byte, len(orikey))
	copy(key, orikey)
	out := make([][]byte, 0)
	for _, edge := range root.Edges {
		newext := append(key, edge.Key...)
		if edge.Endword {
			out = append(out, newext)
		}
		res := complete(&(edge.Node), newext)
		out = append(out, res...)
	}
	return out
}

func findRecursive(root *Node, key []byte) [][]byte {
	if root == nil || len(key) == 0 {
		return nil
	}
	foundElements := 0
	traverseNode := root
	for traverseNode != nil && !traverseNode.IsLeaf() && foundElements < len(key) {
		var nextEdge *Edge
		for _, edge := range traverseNode.Edges {
			if bytes.HasPrefix(key[foundElements:], edge.Key) {
				nextEdge = &edge
				break
			}
		}
		if nextEdge == nil {
			traverseNode = nil
			break
		}
		foundElements += len(nextEdge.Key)
		traverseNode = &(nextEdge.Node)
	}
	if traverseNode == nil || traverseNode.IsLeaf() {
		return nil
	}
	return complete(traverseNode, key)
}

func find(root *Node, in []byte) map[string]Edge {
	if root == nil || len(in) == 0 {
		return nil
	}
	result := make(map[string]Edge)

	var stack []Edge
	stack = append(stack, NewEdge(in, nil))
	for len(stack) > 0 {
		var key Edge
		var foundKey []byte
		key, stack = stack[0], stack[1:]
		if edge, found := result[key.String()]; found && edge.Endword && edge.Node.IsLeaf() {
			continue
		}
		// Need to make a copy of the byte.
		foundKey = make([]byte, len(key.Key))
		copy(foundKey, key.Key)
		var foundElements int
		traverseNode := root
		for traverseNode != nil && !traverseNode.IsLeaf() && foundElements < len(foundKey) {
			var nextEdge *Edge
			for _, edge := range traverseNode.Edges {
				if bytes.HasPrefix(foundKey[foundElements:], edge.Key) {
					nextEdge = &edge
					break
				}
			}
			if nextEdge == nil {
				traverseNode = nil
				break
			}
			foundElements += len(nextEdge.Key)
			traverseNode = &(nextEdge.Node)
		}
		if traverseNode == nil || traverseNode.IsLeaf() {
			continue
		}
		for _, edge := range traverseNode.Edges {
			func(key, edge Edge) {
				out := append(foundKey, edge.Key...)
				stack = append([]Edge{NewEdge(out, nil)}, stack...)
				result[string(out)] = edge
			}(key, edge)
		}
	}

	for id, edge := range result {
		if !edge.Endword {
			delete(result, id)
		}
	}
	return result
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
