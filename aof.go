package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func endcodeResp(operation, key, value string) string {
	operationLength := len(operation)
	keyLenght := len(key)
	valueLenght := len(value)

	return fmt.Sprintf("*3\r\n$%d\r\n$%s\r\n$%d\r\n$%s\r\n%d\r\n$%s\r\n", operationLength, operation, keyLenght, key, valueLenght, value)
}

func aofCompaction(mmap map[string]string, file *os.File) {

	if err := file.Truncate(0); err != nil {
		panic(err)
	}

	_, err := file.Seek(0, io.SeekStart)

	if err != nil {
		panic(err)
	}

	for k, v := range mmap {
		// res := endcodeResp("SET", k, v)
		aofSimple := []string{"set", k, v}
		res := strings.Join(aofSimple, " ") + "\r\n"
		_, err := file.Write([]byte(res))
		if err != nil {
			panic(err)
		}
	}
}
