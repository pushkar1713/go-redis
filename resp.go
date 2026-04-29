package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func endcodeResp(operation, key, value string) string {
	operationLength := len(operation)
	keyLenght := len(key)
	valueLenght := len(value)

	return fmt.Sprintf("*3\r\n$%d\r\n$%s\r\n$%d\r\n$%s\r\n%d\r\n$%s\r\n", operationLength, operation, keyLenght, key, valueLenght, value)
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

func handleCase(commands []string, file *os.File) string {
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
			// dbMutex.Lock()
			mmap[commands[1]] = commands[2]
			// dbMutex.Unlock()
			if lock.Load() {
				fmt.Print("writing in buffer")
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
				// dbMutex.RLock()
				result, ok := mmap[commands[1]]
				// dbMutex.RUnlock()
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
			// dbMutex.Lock()
			delete(mmap, commands[1])
			// dbMutex.Unlock()
			return ":1\r\n"
		}
	case EXISTS:
		{
			// dbMutex.RLock()
			inputKeys := commands[1:]
			existsCount := 0
			for _, val := range inputKeys {
				_, exists := mmap[val]
				if exists {
					existsCount++
				}
			}
			// dbMutex.RUnlock()
			return fmt.Sprintf(":%v\r\n", existsCount)
		}

	}

	return "+OK\r\n"

}

// Todo complete this function
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
