package ahead

import (
	"fmt"
	"strings"
)

// Node holds an array of edge.
type Node struct {
	Edges []*Edge
}

// NewNode returns a new node value.
func NewNode() *Node {
	return &Node{}
}

// IsLeaf returns true if the node does not have any edges.
func (n Node) IsLeaf() bool {
	return len(n.Edges) == 0
}

// Print iteratively prints all the node edges.
func (n Node) Print(depth int) {
	for _, edge := range n.Edges {
		key, count := edge.Key, edge.Count
		fmt.Printf("%s %s:%d\n", strings.Repeat(" ", depth*2), key, count)
		edge.Node.Print(depth + 1)
	}
}
