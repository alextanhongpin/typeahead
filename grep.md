# Boyer-Moore algorithm

Look into how to implement them.

```go
package main

import (
	"fmt"
	"strings"
)

const NO_OF_CHARS = 256

// Best case is O(n/m)
// Worst case is O(mn)
func badCharHeuristic(str []rune) [NO_OF_CHARS]int {
	size := len(str)
	var badChar [NO_OF_CHARS]int
	for i := 0; i < NO_OF_CHARS; i++ {
		badChar[i] = -1
	}
	for i := 0; i < size; i++ {
		badChar[str[i]] = i
	}
	return badChar
}

func search(txt, pat []rune) {
	m, n := len(pat), len(txt)
	badChar := badCharHeuristic(pat)

	// The shift of the pattern with respect to text.
	s := 0
	for s <= n-m {
		// Index j starts with the len of the pattern - 1.
		j := m - 1

		// Reduce index j while the pattern matches those from txt.
		// Note that we are comparing the text from behind.
		for j >= 0 && pat[j] == txt[s+j] {
			j--
		}

		// If pattern is present at current shift, then the index j will be -1.
		if j < 0 {
			fmt.Println("index is", s)

			if s+m < n {
				s += (m - badChar[txt[s+m]])
			} else {
				s++
			}
		} else {
			s += max(1, j-badChar[txt[s+j]])
		}
	}

}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func main() {
	txt, pat := "AABAACAADAABAABA", "AABA"

	fmt.Println(strings.Index(txt, pat))
	search([]rune(txt), []rune(pat))
}
```

References:

- http://web.cs.ucdavis.edu/~gusfield/cs224f11/bnotes.pdf
- https://www.geeksforgeeks.org/boyer-moore-algorithm-for-pattern-searching/
- https://www.geeksforgeeks.org/boyer-moore-algorithm-good-suffix-heuristic/
- https://en.wikipedia.org/wiki/Aho%E2%80%93Corasick_algorithm
