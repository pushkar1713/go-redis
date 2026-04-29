package main

type Value struct {
	RespType       string
	TypeString     string
	TypeArray      []Value
	ProcessedArray []string
}

type Request struct {
	commands []string
	respChan chan string
}
