package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"slices"
	"strconv"
	"strings"
)

const (
	SET     = "SET"
	GET     = "GET"
	COMMAND = "COMMAND"
	PING    = "PING"
)

var allowedCommands = []string{
	SET,
	GET,
	PING,
	COMMAND,
}

var mmap = make(map[string]string)

type Value struct {
	RespType       string
	TypeString     string
	TypeArray      []Value
	ProcessedArray []string
}

func handleConn(conn net.Conn) {
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
			panic("the command is not an array")
		}
		v.RespType = "*"
		if err != nil {
			panic(err)
		}
		args, err := reader.ReadByte()
		if err != nil {
			panic(err)
		}

		fmt.Println((string(args)))
		numArgs, err := strconv.Atoi(string(args))
		if err != nil {
			panic(err)
		}
		reader.ReadByte()
		reader.ReadByte()
		for range numArgs {
			parseBulkStrings(reader, &v)
		}
		response := handleCase(v.ProcessedArray)
		// fmt.Println(v)
		fmt.Printf("%+v\n", v)
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

func handleCase(commands []string) string {
	command := commands[0]
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
				return "+SET commands only support 2 arguments, key and value\r\n"
			}
			mmap[commands[1]] = commands[2]
			return "+DONE\r\n"
		}
	case GET:
		{
			{
				if len(commands) > 2 {
					return "+GET commands only support 1 argument, key\r\n"
				}
				result, ok := mmap[commands[1]]
				if ok {
					return fmt.Sprintf("+%s\r\n", result)
				}
				return "+KEY does not exist\r\n"
			}
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
	args, err := reader.ReadByte()
	if err != nil {
		panic(err)
	}
	fmt.Println((string(args)))
	numArgs, err := strconv.Atoi(string(args))
	if err != nil {
		panic(err)
	}
	buffer := make([]byte, numArgs)
	reader.ReadByte()
	reader.ReadByte()
	res, err := io.ReadFull(reader, buffer)
	command := string(buffer)
	v.ProcessedArray = append(v.ProcessedArray, strings.ToUpper(command))
	if err != nil {
		_ = res
		panic(err)
	}
	ele := Value{
		RespType:   string(b1),
		TypeString: string(buffer),
	}
	v.TypeArray = append(v.TypeArray, ele)
	reader.ReadByte()
	reader.ReadByte()
	fmt.Println(string(buffer))
}

func main() {
	fmt.Println("this is a redis like key value store writtern in go")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Printf("tcp server listening on port %+v\n", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		handleConn(conn)
	}

}
