# go-ahead
A text autocompleter with golang. The goal is to create a google-like
  autocomplete. Based on some research done, it seems like google is
  implementing it using radix/patricia tree. The implementation here is
  probably a slightly more naive version than what google has. 

Radix tree is more storage efficient than a trie (?). 

Another interesting implementation is the _typeahead_ by Facebook.


## TODO

- Create an actual server that can add new search suggestions. 
- Improve the scoring/ranking algorithm.
- Look into a better way to identify performance issue
- Perform testing with quickcheck

## References

- http://dhruvbird.com/autocomplete.pdf
