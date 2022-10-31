package main

import (
	"encoding/json"
	"fmt"
)

type barstruct struct {
	Key   int
	Value int
}

type foo struct {
	Key int
	Bar []barstruct
}

func main() {
	data := &foo{1, []barstruct{{1, 2}, {2, 3}}}
	out, _ := json.Marshal(data)
	fmt.Println(string(out))
}
