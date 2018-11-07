package typeahead

// func main() {
//         tree := NewTernaryTree()
//         tree.Add([]rune("hello"))
//         tree.Add([]rune("ac"))
//
//         fmt.Println("Hello, playground", tree.Contains([]rune("hello")))
//         fmt.Println("Hello, playground", tree.Contains([]rune("xxdda")))
//
// }

type TernaryNode struct {
	char    rune
	left    *TernaryNode
	right   *TernaryNode
	center  *TernaryNode
	endword bool
}

func NewTernaryNode(char rune, endword bool) *TernaryNode {
	return &TernaryNode{
		char:    char,
		endword: endword,
	}
}

type TernaryTree struct {
	root *TernaryNode
}

func NewTernaryTree() *TernaryTree { return &TernaryTree{} }

func (t *TernaryTree) Add(s []rune) {
	t.root = t.add(s, 0, t.root)
}

func (t *TernaryTree) add(s []rune, pos int, node *TernaryNode) *TernaryNode {
	if node == nil {
		node = NewTernaryNode(s[pos], false)
	}
	if s[pos] < node.char {
		node.left = t.add(s, pos, node.left)
	} else if s[pos] > node.char {
		node.right = t.add(s, pos, node.right)
	} else {
		if pos+1 == len(s) {
			node.endword = true
		} else {
			node.center = t.add(s, pos+1, node.center)
		}
	}
	return node
}

func (t *TernaryTree) Contains(s []rune) bool {
	var pos int
	node := t.root
	for node != nil {
		if s[pos] < node.char {
			node = node.left
		} else if s[pos] > node.char {
			node = node.right
		} else {
			pos++
			if pos == len(s) {
				return node.endword
			}
			node = node.center
		}
	}
	return false
}
