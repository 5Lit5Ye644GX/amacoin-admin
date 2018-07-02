package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	cool "github.com/fatih/color"
)

var banner = [6]string{
	"		██████╗ ███████╗██╗   ██╗███████╗    ██╗   ██╗██╗   ██╗██╗  ████████╗\n",
	"		██╔══██╗██╔════╝██║   ██║██╔════╝    ██║   ██║██║   ██║██║  ╚══██╔══╝\n",
	"		██║  ██║█████╗  ██║   ██║███████╗    ██║   ██║██║   ██║██║     ██║\n",
	"		██║  ██║██╔══╝  ██║   ██║╚════██║    ╚██╗ ██╔╝██║   ██║██║     ██║\n",
	"		██████╔╝███████╗╚██████╔╝███████║     ╚████╔╝ ╚██████╔╝███████╗██║\n",
	"		╚═════╝ ╚══════╝ ╚═════╝ ╚══════╝      ╚═══╝   ╚═════╝ ╚══════╝╚═╝\n",
}

func print(msg string, maxtime ...int) {
	waiting := 20
	if len(maxtime) > 0 {
		waiting = maxtime[0]
	}

	runes := []rune(msg)
	for _, c := range runes {
		time.Sleep(time.Duration(rand.Intn(waiting)) * time.Millisecond)
		fmt.Printf("%c", c)
	}
}

func ok(msg string) {
	cool.New(cool.FgHiGreen).Printf("[OK] ")
	fmt.Println(msg)
}

func okf(format string, a ...interface{}) {
	cool.New(cool.FgHiGreen).Printf("[OK] ")
	fmt.Printf(format, a)
}

func fail(msg string) {
	cool.New(cool.FgHiRed).Printf("[ERROR] ")
	fmt.Println(msg)
}

func failf(format string, a ...interface{}) {
	cool.New(cool.FgHiRed).Printf("[ERROR] ")
	fmt.Printf(format, a)
}

// Clear screen
func clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
