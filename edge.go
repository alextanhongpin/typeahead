package ahead

// Edge represents the edge of a node.
type Edge struct {
	Count   int
	Key     []byte
	Value   interface{}
	Node    *Node
	Endword bool
}

// NewEdge creates a new Edge with the given key value pair.
func NewEdge(key []byte, value interface{}) *Edge {
	return &Edge{
		Key:   key,
		Value: value,
		Count: 1,
		Node:  NewNode(),
	}
}

func (e Edge) String() string {
	return string(e.Key)
}
