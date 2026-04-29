package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

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
