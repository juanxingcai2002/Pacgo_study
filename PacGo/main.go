package main

import (
	"bufio"
	"fmt"
	"os"
)

var maze []string

func loadMaze(file string) error {
	f, err := os.Open(file) //The `os.Open()` function returns a pair of values: a file and an error
	if err != nil {
		return err
	}
	defer f.Close() // 当loadMaze函数结束后再调用f.close()

	scanner := bufio.NewScanner(f) //将文件读到内存中（即读到一个变量上）
	for scanner.Scan() {           //如果还可以从文件中读取内容，.scan会返回truel
		line := scanner.Text() // scanner.Text返回文件下一行的输入
		maze = append(maze, line)
	}
	return nil

}

func printScreen() {
	for _, line := range maze { // 对于数组和切片，rang-for第一个参数是索引，第二个参数是数据
		fmt.Println(line)
	}
}

func main() {

	//initialize game

	// load resources
	err := loadMaze("maze01.txt")
	if err != nil {
		fmt.Println("failed ti load maze:", err)
		return
	}

	// game loop
	for {
		// update screen
		printScreen()

		// process input

		// process movement

		// process collisions

		// check game over

		// temp:break infinite loop
		break
		// repeat
	}

}
