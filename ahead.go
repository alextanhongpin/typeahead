package ahead

import (
	"strings"
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
func (r *Root) Insert(key string, value interface{}) {
	insert(&(r.Node), key, value)
}

// Find searches for the edge of the node that matches the given prefix.
func (r *Root) Find(key string) map[string]Edge {
	return find(&(r.Node), key)
}

func (r *Root) FindRecursive(key string) []string {
	return findRecursive(&(r.Node), key)
}

func insert(root *Node, key string, value interface{}) {
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
	if currKey == key {
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

func split(root *Node, key string, p, pos int) {
	var rem int
	edge := root.Edges[pos]
	for k, v := range root.Edges {
		if v.Key == edge.Key {
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

func complete(root *Node, key string) []string {
	// key := make(string, len(orikey))
	// copy(key, orikey)
	out := make([]string, 0)
	var sb strings.Builder
	for _, edge := range root.Edges {
		sb.Reset()
		sb.WriteString(key)
		sb.WriteString(edge.Key)
		str := sb.String()
		if edge.Endword {
			out = append(out, str)
		}
		res := complete(&(edge.Node), str)
		out = append(out, res...)
	}
	return out
}

func findRecursive(root *Node, key string) []string {
	if root == nil || len(key) == 0 {
		return nil
	}
	foundElements := 0
	traverseNode := root
	for traverseNode != nil && !traverseNode.IsLeaf() && foundElements < len(key) {
		var nextEdge *Edge
		for _, edge := range traverseNode.Edges {
			if strings.HasPrefix(key[foundElements:], edge.Key) {
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

func find(root *Node, in string) map[string]Edge {
	if root == nil || len(in) == 0 {
		return nil
	}
	result := make(map[string]Edge)

	var stack []Edge
	stack = append(stack, NewEdge(in, nil))
	for len(stack) > 0 {
		var key Edge
		var foundKey string
		key, stack = stack[0], stack[1:]
		if edge, found := result[key.String()]; found && edge.Endword && edge.Node.IsLeaf() {
			continue
		}
		// Need to make a copy of the byte.
		// foundKey = make(string, len(key.Key))
		// copy(foundKey, key.Key)
		var foundElements int
		traverseNode := root
		for traverseNode != nil && !traverseNode.IsLeaf() && foundElements < len(key.Key) {
			var nextEdge *Edge
			for _, edge := range traverseNode.Edges {
				if strings.HasPrefix(key.Key[foundElements:], edge.Key) {
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
		var sb strings.Builder
		for _, edge := range traverseNode.Edges {
			sb.Reset()

			sb.WriteString(foundKey)
			sb.WriteString(edge.Key)
			str := sb.String()
			// out := append(foundKey, edge.Key...)
			stack = append([]Edge{NewEdge(str, nil)}, stack...)
			result[str] = edge
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

func sharedPrefix(s, t string) int {
	minLen := min(len(s), len(t))
	for i := 0; i < minLen; i++ {
		if s[i] != t[i] {
			return i
		}
	}
	return minLen
}
