package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	SET     = "SET"
	GET     = "GET"
	COMMAND = "COMMAND"
	PING    = "PING"
	DEL     = "DEL"
	EXISTS  = "EXISTS"
)

var allowedCommands = []string{
	SET,
	GET,
	PING,
	COMMAND,
	DEL,
	EXISTS,
}

var mainChan = make(chan Request)
var mmap = make(map[string]string)

// var dbMutex sync.RWMutex
var aofBuffer = []string{}
var lock atomic.Bool

func fsync(s []string, file *os.File) {
	for _, v := range s {
		file.Write([]byte(v))
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		v := Value{}
		b1, err := reader.ReadByte()
		if string(b1) != "*" {
			// panic("the command is not an array")
			continue
		}
		v.RespType = "*"
		if err != nil {
			panic(err)
		}
		// args, err := reader.ReadByte()
		args, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		// fmt.Println((string(args)))
		numArgs, err := strconv.Atoi(strings.TrimSpace(args))
		if err != nil {
			panic(err)
		}
		for range numArgs {
			parseBulkStrings(reader, &v)
		}

		respChan := make(chan string)
		mainChan <- Request{commands: v.ProcessedArray, respChan: respChan}
		// response := handleCase(v.ProcessedArray, file)
		response := <-respChan
		conn.Write([]byte(response))
	}
}

func checkAllowedCommands(command string) bool {
	return slices.Contains(allowedCommands, command)
}

func aofCompactionAndBufferSync(file *os.File) {
	if !lock.CompareAndSwap(false, true) {
		return
	}
	defer lock.Store(false)
	aofCompaction(mmap, file)
	fmt.Println(aofBuffer)
	fsync(aofBuffer, file)
}

func main() {
	fmt.Println("this is a redis like key value store writtern in go")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile("aof.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println(mmap)
	aofRestore(file)
	fmt.Println(mmap)

	go func() {
		for range time.Tick(60 * time.Second) {
			fmt.Println("this is compactions starting")
			aofCompactionAndBufferSync(file)
		}
	}()

	defer listener.Close()
	fmt.Printf("tcp server listening on port %+v\n", listener.Addr())

	go func() {
		for req := range mainChan {
			response := handleCase(req.commands, file)
			req.respChan <- response
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}

}
