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

// http://hacktalks.blogspot.com/2012/03/implementing-auto-complete-with-ternary.html
func main() {
	tree := NewTernaryTree()
	tree.Add("hello")
	tree.Add("helsinki")
	tree.Add("dobby")
	tree.Add("sell")
	tree.Add("hallo")
	tree.Add("hi")
	tree.Add("car")
	fmt.Println("has hello:", tree.Contains("hello"))
	result := tree.NewSearch("h")
	fmt.Println(result)
}

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

// Add adds the item to the tree recursively.
func (t *TernaryTree) Add(s string) {
	t.root = t.radd([]rune(s), 0, t.root)
}

func (t *TernaryTree) radd(s []rune, pos int, node *TernaryNode) *TernaryNode {
	if node == nil {
		node = NewTernaryNode(s[pos], false)
	}
	if s[pos] < node.char {
		node.left = t.radd(s, pos, node.left)
	} else if s[pos] > node.char {
		node.right = t.radd(s, pos, node.right)
	} else {
		if pos+1 == len(s) {
			node.endword = true
		} else {
			node.center = t.radd(s, pos+1, node.center)
		}
	}
	return node
}

func (t *TernaryTree) Contains(str string) bool {
	s := []rune(str)
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

// Search implements an autocomplete for ternary search tree.
func (t *TernaryTree) Search(str string) {
	t.search(t.root, []rune(str), "")
}

func (t *TernaryTree) NewSearch(str string) []string {
	s := []rune(str)
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
				break
			}
			node = node.center
		}
	}
	var result []string
	t.runedfs(node.center, append(s, node.center.char), &result)
	return result
}

func (t *TernaryTree) search(node *TernaryNode, str []rune, match string) {
	if len(str) > 0 {
		key := str[0]
		if key < node.char {
			t.search(node.left, str, match)
		} else if key > node.char {
			t.search(node.right, str, match)
		} else {
			if len(str) == 1 {
				if node.endword {
					// Match!
					fmt.Println("match", match, string(node.char))
				}
				t.dfs(node.center, match+string(node.char)+string(node.center.char))
				return
			}
			t.search(node.center, str[1:], string(match)+string(key))
		}
	}
}

// Depth first search for ternary tree.
func (t *TernaryTree) dfs(node *TernaryNode, match string) {
	if node.endword {
		fmt.Println("match", match)
	}
	if node.center == nil && node.left == nil && node.right == nil {
		return
	}
	if node.center != nil {
		t.dfs(node.center, match+string(node.center.char))
	}
	if node.right != nil {
		t.dfs(node.right, match[:len(match)-1]+string(node.right.char))
	}
	if node.left != nil {
		t.dfs(node.left, match[:len(match)-1]+string(node.left.char))
	}
}

func (t *TernaryTree) runedfs(node *TernaryNode, match []rune, result *[]string) {
	if node.endword {
		*result = append(*result, string(match))
	}
	if node.center == nil && node.left == nil && node.right == nil {
		return
	}
	if node.center != nil {
		t.runedfs(node.center, append(match, node.center.char), result)
	}
	if node.right != nil {
		match[len(match)-1] = node.right.char
		t.runedfs(node.right, match, result)
	}
	if node.left != nil {
		match[len(match)-1] = node.left.char
		t.runedfs(node.left, match, result)
	}
}
