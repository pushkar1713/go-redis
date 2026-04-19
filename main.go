package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
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

var mmap = make(map[string]string)
var dbMutex sync.RWMutex
var aofBuffer = []string{}
var lock atomic.Bool

type Value struct {
	RespType       string
	TypeString     string
	TypeArray      []Value
	ProcessedArray []string
}

func fsync(s []string, file *os.File) {
	for _, v := range s {
		file.Write([]byte(v))
	}
}

func handleConn(conn net.Conn, file *os.File) {
	defer conn.Close()
	// buffer := make([]byte, 1024)
	// var array []string
	reader := bufio.NewReader(conn)
	for {
		// line, err := reader.ReadString('\n')
		// fmt.Println(line[1])
		// fmt.Println(strconv.Atoi(strings.TrimSpace(line[1:])))
		// b1, err := reader.ReadByte()
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Printf("%v %q\n", b1, b1)
		// b2, err := reader.ReadByte()
		// if err != nil {
		// 	panic(err)
		// }
		// buffer := make([]byte, 1024)
		// var n int
		// for range int(b2 - '0') {
		// 	n, err = reader.Read(buffer)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }

		// array = append(array, string(b))
		// n, err := conn.Read(buffer)
		// conn.Write([]byte("+OK\r\n"))
		// fmt.Println(string(buffer[:n]))
		// fmt.Println(array)
		//
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
		response := handleCase(v.ProcessedArray, file)
		// fmt.Println(v)
		// fmt.Printf("%+v\n", v)
		// conn.Write([]byte("+OK\r\n"))
		conn.Write([]byte(response))
	}
}

func checkAllowedCommands(command string) bool {
	// for _, c := range allowedCommands {
	// 	if c == command {
	// 		return true
	// 	}
	// }
	return slices.Contains(allowedCommands, command)
}

func handleCase(commands []string, file *os.File) string {
	// dbMutex.RLock()
	// fmt.Println(mmap)
	// dbMutex.RUnlock()
	command := strings.ToUpper(commands[0])
	aofLog := strings.Join(commands, " ") + "\r\n"
	fmt.Println(command)
	if checkAllowedCommands(command) == false {
		return "-ERR command not supported\r\n"
	}
	switch command {
	case PING:
		return "+PONG\r\n"
	case SET:
		{
			if len(commands) > 3 {
				return "-ERR SET commands only support 2 arguments, key and value\r\n"
			}
			dbMutex.Lock()
			mmap[commands[1]] = commands[2]
			dbMutex.Unlock()
			if lock.Load() {
				aofBuffer = append(aofBuffer, aofLog)
				return "+OK\r\n"
			}
			_, err := file.Write([]byte(aofLog))
			if err != nil {
				panic(err)
			}
			return "+OK\r\n"
		}
	case GET:
		{
			{
				if len(commands) > 2 {
					return "-ERR GET commands only support 1 argument, key\r\n"
				}
				dbMutex.RLock()
				result, ok := mmap[commands[1]]
				dbMutex.RUnlock()
				if ok {
					return fmt.Sprintf("$%d\r\n%s\r\n", len(result), result)
				}
				return "$-1\r\n"
			}
		}
	case DEL:
		{
			if len(commands) > 2 {
				return "-ERR DEL commands only support 1 argument, key\r\n"
			}
			dbMutex.Lock()
			delete(mmap, commands[1])
			dbMutex.Unlock()
			return ":1\r\n"
		}
	case EXISTS:
		{
			dbMutex.RLock()
			inputKeys := commands[1:]
			existsCount := 0
			for _, val := range inputKeys {
				_, exists := mmap[val]
				if exists {
					existsCount++
				}
			}
			dbMutex.RUnlock()
			return fmt.Sprintf(":%v\r\n", existsCount)
		}

	}

	return "+OK\r\n"

}

func parseIntegers(reader *bufio.Reader) {
	b1, err := reader.ReadByte()
	if rune(b1) != ':' {
		panic("not an integer")
	}
	if err != nil {
		panic(err)
	}
	b2, err := reader.ReadByte()
	if rune(b2) == '-' || rune(b2) == '+' {
	}
	for {

	}
}

func parseBulkStrings(reader *bufio.Reader, v *Value) {
	b1, err := reader.ReadByte()
	if rune(b1) != '$' {
		panic("not a bulk string")
	}
	if err != nil {
		panic(err)
	}
	//this assumes that the lenght of the argument is in single digit but that is actully wrong
	// args, err := reader.ReadByte()
	args, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	// fmt.Println(args)
	numArgs, err := strconv.Atoi(strings.TrimSpace(args))
	_ = numArgs
	// fmt.Printf("this is numArgs %+v\n", numArgs)
	if err != nil {
		panic(err)
	}
	buffer := make([]byte, numArgs)
	res, err := io.ReadFull(reader, buffer)
	_ = res
	// res, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	command := string(buffer)
	// set := "set"
	// if strings.Contains(command, set) {
	// 	_, err := file.Write([]byte(res))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	v.ProcessedArray = append(v.ProcessedArray, command)
	ele := Value{
		RespType:   string(b1),
		TypeString: string(command),
	}
	v.TypeArray = append(v.TypeArray, ele)
	reader.ReadByte()
	reader.ReadByte()
	// fmt.Println(string(buffer))
}

func aofRestore(file *os.File) {
	info, err := file.Stat()
	if err != nil {
		panic(err)
	}
	if info.Size() == 0 {
		fmt.Println("aof file is empty")
		return
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		commands := strings.Split(line, " ")
		command1 := strings.ToUpper(commands[0])
		set := "SET"
		if strings.Contains(command1, set) {
			mmap[commands[1]] = commands[2]
		}
	}

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

	time.AfterFunc(5*time.Second, func() {
		fmt.Println("doing the fsync compaction")
		if !lock.CompareAndSwap(false, true) {
			return
		}
		defer lock.Store(false)
		aofCompaction(mmap, file)
		fsync(aofBuffer, file)
	})

	defer listener.Close()
	fmt.Printf("tcp server listening on port %+v\n", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn, file)
	}

}
