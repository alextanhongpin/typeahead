package typeahead

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"testing/quick"
)

func TestRadix(t *testing.T) {
	f := func(in []byte) bool {
		root := NewNode()
		insert(&root, in, nil)
		out := find(&root, in[:len(in)/2])
		if len(out) != 1 {
			fmt.Println(string(in), len(out), "search", string(in[:len(in)/2]))
			fmt.Println(string(root.edges[0].key))
			return false
		}
		for _, v := range out {
			if ok := bytes.Equal(in, v.key); !ok {
				fmt.Println(string(in))
				return false
			}
		}
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		log.Fatal(err)
	}
}
