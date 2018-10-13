package ahead

import (
	"bytes"
)

// Root represents the root of the radix tree.
type Root struct {
	node Node
}

// New returns a new tree.
func New() *Root {
	return &Root{
		node: NewNode(),
	}
}

// Insert adds a key value pair into the tree.
func (r *Root) Insert(key []byte, value interface{}) {
	insert(&(r.node), key, value)
}

// Find searches for the edge of the node that matches the given prefix.
func (r *Root) Find(key []byte) map[string]Edge {
	return find(&(r.node), key)
}

func (r *Root) FindRecursive(key []byte) [][]byte {
	return findRecursive(&(r.node), key)
}

func insert(root *Node, key []byte, value interface{}) {
	if root == nil || len(key) == 0 {
		return
	}
	var p, pos int
	for i, child := range root.edges {
		if len(child.key) == 0 || child.key[0] != key[0] {
			continue
		}
		p = sharedPrefix(child.key, key)
		pos = i
		break
	}
	if p == 0 {
		edge := NewEdge(key, value)
		edge.endword = true
		root.edges = append(root.edges, edge)
		return
	}
	currKey := root.edges[pos].key
	if bytes.Equal(currKey, key) {
		root.edges[pos].count++
		root.edges[pos].endword = true
		return
	}
	if len(currKey) == p {
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
	// Pop the old edge.
	root.edges = append(root.edges[:rem], root.edges[rem+1:]...)
	prefix, left, right := edge.key[:p], edge.key[p:], key[p:]

	newEdge := NewEdge(prefix, nil)
	newEdge.count += edge.count

	edge.key = left

	insert(&(newEdge.node), right, nil)
	newEdge.node.edges = append(newEdge.node.edges, edge)
	root.edges = append(root.edges, newEdge)
}

func complete(root *Node, orikey []byte) [][]byte {
	key := make([]byte, len(orikey))
	copy(key, orikey)
	out := make([][]byte, 0)
	for _, edge := range root.edges {
		newext := append(key, edge.key...)
		if edge.endword {
			out = append(out, newext)
		}
		res := complete(&(edge.node), newext)
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
		if edge, found := result[key.String()]; found && edge.endword && edge.node.IsLeaf() {
			continue
		}
		// Need to make a copy of the byte.
		foundKey = make([]byte, len(key.key))
		copy(foundKey, key.key)
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
				stack = append([]Edge{NewEdge(out, nil)}, stack...)
				result[string(out)] = edge
			}(key, edge)
		}
	}

	for id, edge := range result {
		if !edge.endword {
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
