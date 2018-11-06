package typeahead

import (
	"strings"
)

// http://www.cs.yale.edu/homes/aspnes/pinewiki/RadixSearch.html?highlight=%28CategoryAlgorithmNotes%29
// PerByte = 8
var BitsPerByte = 8

const TrieBase = 2

// func main() {
//         var t *Trie
//         t = trieInsert(t, "hello")
//         t = trieInsert(t, "car")
//         fmt.Println("has hello", trieContains(t, "hello"))
//         fmt.Println("has cac", trieContains(t, "cac"))
//         fmt.Println("has cac", trieContains(t, "car"))
// }

func GetBit(key string, n int) int {
	if len(key) == n/BitsPerByte {
		return 0
	}
	if key[n/BitsPerByte]&(0x1<<uint(BitsPerByte-1-n%BitsPerByte)) != 0 {
		return 1
	}
	return 0
}

type Trie struct {
	key      string
	children [TrieBase]*Trie
}

func NewTrie(key string) *Trie {
	return &Trie{
		key: key,
		// children: make([TrieBase]*Trie, TrieBase),
	}
}

func isLeaf(trie *Trie) bool {
	if trie == nil {
		return true
	}
	return (trie.children[0] == nil) && (trie.children[1] == nil)
}

func TrieInsert(trie *Trie, key string) *Trie {
	var bit, bitvalue int
	var t, kid *Trie
	var oldKey string

	if trie == nil {
		return NewTrie(key)
	}

	// Search for the key.
	for t = trie; !isLeaf(t); bit, t = bit+1, kid {
		bitvalue = GetBit(key, bit)
		kid = t.children[bitvalue]
		if kid == nil {
			t.children[bitvalue] = NewTrie(key)
			return trie
		}
	}
	// Nothing to do here.
	// if strings.EqualFold(t.key, key) {
	if strings.Compare(t.key, key) == 0 {
		// if t.key == key {
		return trie
	}

	// Extend the trie.
	oldKey = t.key
	t.key = ""

	// Walk the common prefix.
	// bitvalue = GetBit(key, bit)
	bitvalue = GetBit(key, bit)
	for GetBit(oldKey, bit) == bitvalue {
		kid = NewTrie("")
		t.children[bitvalue] = kid
		bit++
		t = kid
		bitvalue = GetBit(key, bit)
	}
	// Then split.
	t.children[bitvalue] = NewTrie(key)

	// There is no NOT operator in golang.
	t.children[1^bitvalue] = NewTrie(oldKey)
	return trie
}

func TrieContains(trie *Trie, target string) bool {
	for bit := 0; trie != nil && !isLeaf(trie); bit++ {
		trie = trie.children[GetBit(target, bit)]
	}
	if trie == nil {
		return false
	}
	return strings.EqualFold(trie.key, target)
}
