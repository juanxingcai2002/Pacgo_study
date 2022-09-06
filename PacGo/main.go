package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/danicat/simpleansi" //用来我们游戏for循环完之后给迭代给我们一个空白的屏幕，用的是转义序列
)

//创建精灵类型保存精灵这个玩家的信息用结构体
type sprite struct {
	row int
	col int
}

var player sprite //创建一个精灵类型的玩家

var ghosts []*sprite //创建幽灵，和精灵的创建不同没有直接存储在内存上，我们用以系列的指针去指向，方便高效的更新精灵的位置

var score int   //记录精灵的得分
var numDots int // 记录点数，精灵吃到一个点就得分
var lives = 1   //表示精灵的生命

// 启动Cbreak模式
func initialise() {
	cbTerm := exec.Command("stty", "cbreak", "-echo") // 返回的是是一个命令（c*Cmd），是执行完stty的命令，参数是cbreak和-echo
	cbTerm.Stdin = os.Stdin                           // 将标准输入读入

	err := cbTerm.Run() //func (c *Cmd) Run() error如果该命令运行，则返回的错误为nil，复制stdin、stdout和stderr没有问题，并以零退出状态退出,如果命令启动但没有成功完成，则错误类型为*ExitError。 对于其他情况，可能会返回其他错误类型。

	if err != nil {
		log.Fatalln("unable to active cbreak mode:", err) //Fatalln等价于Println()之后调用os.Exit(1)。说明是错误退出
		//func Fatalln(v ...any)函数原型，说明后面跟各种类型的参数都可以
	}
}

//生成随机数控制幽灵随机的移动，用随机数的整数通过map这个数据结构映射成字符串
func drawDirection() string {
	dir := rand.Intn(4)
	move := map[int]string{
		0: "UP",
		1: "DOWN",
		2: "RIGHT",
		3: "LEFT",
	}
	return move[dir]
}

// 处理幽灵的移动
func moveGhosts() {
	for _, g := range ghosts {
		dir := drawDirection()                     // 随机得到一个方向移动的指令
		g.row, g.col = makeMove(g.row, g.col, dir) // makeMove return 移动过后的位置

	}
}

// 启动Cooked mode
func cleanup() {
	cookedTerm := exec.Command("stty", "-cbreak", "echo")
	cookedTerm.Stdin = os.Stdin

	err := cookedTerm.Run()
	if err != nil {
		log.Fatalln("unable to restore cooked mode:", err)
	}
}

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

	//地图中P的位置就是精灵的位置，找到精灵的位置（行和列）
	for row, line := range maze { //遍历每一行
		for col, char := range line { // 遍历每一行的每一列
			switch char {
			case 'P':
				player = sprite{row: row, col: col}
			case 'G':
				ghosts = append(ghosts, &sprite{row, col}) //这意味着我们不是在切片中添加一个精灵对象，而是添加一个指向它的指针。
			case '.':
				numDots++ //记录地图上点的数量
			}
		}
	}

	return nil

}

//处理精灵的移动
func makeMove(oldRow, oldCol int, dir string) (newRow, newCol int) {
	newRow, newCol = oldRow, oldCol //精灵可能不移动，新位置等于原来的位置
	switch dir {                    // 处理从标准输入中得到的按键字符，游戏里的坐标和平常数学上的是相反的
	case "UP":
		newRow = newRow - 1 // 特性是：如果移动越界，那么回到当前位置的正对面
		if newRow < 0 {
			newRow = len(maze) - 1
		}
	case "DOWN":
		newRow = newRow + 1
		if newRow == len(maze) {
			newRow = 0
		}
	case "RIGHT":
		newCol = newCol + 1
		if newCol == len(maze[0]) {
			newCol = 0
		}
	case "LEFT":
		newCol = newCol - 1
		if newCol < 0 {
			newCol = len(maze) - 1
		}
	}
	if maze[newRow][newCol] == '#' { // #代表障碍物，如果最终的的位置是障碍物，那么回到原来的位置
		newRow = oldRow
		newCol = oldCol
	}
	return // 可以给返回值取名，这里取做 newrow 和 oldRow，因此返回值可以不需要显示的指出
}

// 更新精灵玩家的信息
func movePlayer(dir string) {
	player.row, player.col = makeMove(player.row, player.col, dir)
	switch maze[player.row][player.col] {
	case '.': //玩家的新位置在点上的话，玩家要吃了那个点然后得分
		numDots--
		score++
		maze[player.row] = maze[player.row][0:player.col] + " " + maze[player.row][player.col+1:]
	}
}

func printScreen() {
	simpleansi.ClearScreen()
	for _, line := range maze { // 对于数组和切片，rang-for第一个参数是索引，第二个参数是数据
		for _, chr := range line {
			switch chr { // 刷新屏幕的时候，不将全部的图片打印下来，只打印障碍物
			case '#':
				fallthrough //fmt.Printf("%c", chr) // fallthrough 会去掉go中隐式的 break；
			case '.':
				fmt.Printf("%c", chr) //需要把点字符在屏幕上也要显示
			default:
				fmt.Print(" ")
			}
		}
	}
	simpleansi.MoveCursor(player.row, player.row) // 在任意位置打印玩家
	fmt.Print("P")

	for _, g := range ghosts {
		simpleansi.MoveCursor(g.row, g.col)
		fmt.Print("G")
	}

	simpleansi.MoveCursor(len(maze)+1, 0)
	fmt.Println("Score:", score, "\tLives:", lives)
}

//从标准输入中读取内容
func readInput() (string, error) {
	buffer := make([]byte, 100) // 创建一个大小为100字节的数组来装从标准输入中读取的东西

	cnt, err := os.Stdin.Read(buffer) //从标准输入中读取东西到buffer中，返回成功读取的字节数和错误
	if err != nil {                   // 如果发生错误
		return "", err
	}
	if cnt == 1 && buffer[0] == 0x1b { // ESC的16进制编码就是0x1b
		return "ESC", nil
	} else if cnt >= 3 {
		if buffer[0] == 0x1b && buffer[1] == '[' { //箭头键的转义序列为3字节长，以' ESC+['开始，然后从A到D的一个字母。
			switch buffer[2] {
			case 'A':
				return "UP", nil
			case 'B':
				return "DOWN", nil
			case 'C':
				return "RIGHT", nil
			case 'D':
				return "LEFT", nil
			}
		}
	}
	return "", nil //没有出错的其他情况，返回无措nil和空字符串
}

func main() {

	//initialize game
	initialise()
	defer cleanup() // 游戏模式下是cbcooked，平常的证明模式下都是cooked

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
		input := make(chan string) //创建一个名字叫input的通道
		go func(ch chan<- string) {
			for {
				input, err := readInput()
				if err != nil {
					log.Println("error reading input", err)
					ch <- "ESC"
				}
				ch <- input
			}
		}(input)
		// process movement
		select {
		case inp := <-input:
			if inp == "ESC" {
				lives = 0
			}
			movePlayer(inp)
		default:
		}
		// process collisions

		for _, g := range ghosts {
			if player == *g { // 玩家和幽灵碰在了一起，玩家死亡
				lives = 0
			}
		}

		// check game over
		if numDots == 0 || lives <= 0 { //得不到分或者生命为0就退出
			break
		}

		// temp:break infinite loop
		// repeat
		time.Sleep(200 * time.Microsecond)
	}

}
