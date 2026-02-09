package console

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type color struct {
	Red      string
	Green    string
	SysReset string
}

var Color = color{
	Red:      "\033[31m",
	Green:    "\033[32m",
	SysReset: "\033[0m",
}

func cleanMsg(msg string) string {
	return strings.TrimSpace(msg)
}

func Error(msg string) {
	_, _ = io.WriteString(
		os.Stdout,
		fmt.Sprintf(
			"%sError: %s%s\n",
			Color.Red,
			cleanMsg(msg),
			Color.SysReset,
		),
	)
}

func Info(msg string) {
	_, _ = io.WriteString(
		os.Stdout,
		fmt.Sprintf(
			"%s%s%s\n",
			Color.Green,
			cleanMsg(msg),
			Color.SysReset,
		),
	)
}

func PressEnterToExit() {
	fmt.Println("Press Enter to exit...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
