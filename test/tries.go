package main

func main() {
	setOfStrings:=[...]string{
		"pqrs",
		"pprt",
		"psst",
		"qqrs",
		"pqrs",
	}
}

type Trie struct{
	root *TrieNode
}
func (t *Trie) query(s string){
	current:=t.root.next(s[0])
}

type TrieNode struct{
	terminating int
	trieNodes *[]TrieNode
}

func (tn *TrieNode) next(c string){
	return tn.trieNodes[c]
}