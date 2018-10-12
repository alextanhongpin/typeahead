package ahead

// Edge represents the edge of a node.
type Edge struct {
	count   int
	key     []byte
	value   interface{}
	node    Node
	endword bool
}

// NewEdge creates a new Edge with the given key value pair.
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
