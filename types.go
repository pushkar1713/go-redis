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

const INVALID_ERR string = "-ERR invalid input\r\n"
const INTERNAL_ERR string = "-ERR internal server error, contact dev\r\n"
