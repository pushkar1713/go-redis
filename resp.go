package main

import (
	"bufio"
	"errors"
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

func parseBulkStrings(reader *bufio.Reader) (string, error) {
	b1, err := reader.ReadByte()
	if rune(b1) != '$' {
		return "", errors.New("this is not a bulk string")
	}
	if err != nil {
		return "", err
	}
	//this assumes that the lenght of the argument is in single digit but that is actully wrong
	// args, err := reader.ReadByte()
	args, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	// fmt.Println(args)
	numArgs, err := strconv.Atoi(strings.TrimSpace(args))
	_ = numArgs
	// fmt.Printf("this is numArgs %+v\n", numArgs)
	if err != nil || numArgs <= 0 {

		var newError error

		if err == nil {
			newError = errors.New(INVALID_ERR)

		} else {
			newError = err
		}

		return "", newError
	}
	buffer := make([]byte, numArgs)
	res, err := io.ReadFull(reader, buffer)
	_ = res
	// res, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	command := string(buffer)
	// set := "set"
	// if strings.Contains(command, set) {
	// 	_, err := file.Write([]byte(res))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// v.ProcessedArray = append(v.ProcessedArray, command)

	// type array is not required right now, we will improve and use this later
	// ele := Value{
	// 	RespType:   string(b1),
	// 	TypeString: string(command),
	// }
	// v.TypeArray = append(v.TypeArray, ele)
	reader.ReadByte()
	reader.ReadByte()
	// fmt.Println(string(buffer))

	return command, nil
}

func handleCase(commands []string, file *os.File) (string, error) {
	if len(commands) < 1 {
		handleCaseError := errors.New(INVALID_ERR)
		return "", handleCaseError
	}
	command := strings.ToUpper(commands[0])
	aofLog := strings.Join(commands, " ") + "\r\n"
	fmt.Println(command)
	if checkAllowedCommands(command) == false {
		return "", errors.New("-ERR command not supported\r\n")

	}
	switch command {
	case PING:
		return "+PONG\r\n", nil
	case SET:
		{
			if len(commands) < 3 || len(commands) > 3 {
				return "", errors.New("-ERR SET commands only support 2 arguments, key and value\r\n")
			}
			// dbMutex.Lock()
			mmap[commands[1]] = commands[2]
			// dbMutex.Unlock()
			if lock.Load() {
				fmt.Print("writing in buffer")
				aofBuffer = append(aofBuffer, aofLog)
				return "+OK\r\n", nil
			}
			_, err := file.Write([]byte(aofLog))
			if err != nil {
				fmt.Print(err)
				return "", errors.New(INTERNAL_ERR)
			}
			return "+OK\r\n", nil
		}
	case GET:
		{
			{
				if len(commands) < 2 || len(commands) > 2 {
					return "", errors.New("-ERR GET commands only support 1 argument, key\r\n")
				}
				// dbMutex.RLock()
				result, ok := mmap[commands[1]]
				// dbMutex.RUnlock()
				if ok {
					return fmt.Sprintf("$%d\r\n%s\r\n", len(result), result), nil
				}
				return "$-1\r\n", nil
			}
		}
	case DEL:
		{
			if len(commands) < 2 {
				return "", errors.New("-ERR Invalid number of arguments, key\r\n")
			}
			// dbMutex.Lock()
			var count int
			for i := 1; i < len(commands); i++ {
				_, ok := mmap[commands[i]]
				if ok {
					count++
				}
				delete(mmap, commands[1])
			}
			// dbMutex.Unlock()
			return fmt.Sprintf(":%d\r\n", count), nil
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
			return fmt.Sprintf(":%v\r\n", existsCount), nil
		}

	}

	return "+OK\r\n", nil

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
