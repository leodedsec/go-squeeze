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
	Red:      "\u001B[31m",
	Green:    "\u001B[32m",
	SysReset: "\u001B[0m",
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

func Table(header string, rows map[string]string) {
	// -----------------------------------
	// || One  | ... ||
	// || Two  | ... ||
	// || Free | ... ||
	// -----------------------------------

	var result [][]string

	left := fmt.Sprintf("%s|| %s", Color.Green, Color.SysReset)
	middle := " | "
	right := fmt.Sprintf(" %s||%s", Color.Green, Color.SysReset)

	for key, value := range rows {
		tsKey := strings.TrimSpace(key)
		tsValue := strings.TrimSpace(value)

		slicedRow := []string{left, tsKey, middle, tsValue, right}
		result = append(result, slicedRow)
	}

	if result == nil {
		return
	}

	keyIndex := 1
	valueIndex := 3

	maxKeyLen := 0
	maxValueLen := 0

	for _, r := range result {
		maxKeyLen = max(maxKeyLen, len(r[keyIndex]))       // rewrite key len to max
		maxValueLen = max(maxValueLen, len(r[valueIndex])) // rewrite value len to max
	}

	for i, r := range result {
		key := r[keyIndex]
		value := r[valueIndex]

		keyDiff := maxKeyLen - len(key)
		if keyDiff > 0 {
			result[i][keyIndex] = fmt.Sprintf("%s%s", key, strings.Repeat(" ", keyDiff))
		}
		valueDiff := maxValueLen - len(value)
		if valueDiff > 0 {
			result[i][valueIndex] = fmt.Sprintf("%s%s", value, strings.Repeat(" ", valueDiff))
		}
	}

	width := 0
	strTable := ""

	for _, r := range result {
		strRow := strings.Join(r, "")
		width = max(width, len(strRow)-18)
		strTable += strRow + "\n"
	}

	if len(header) != 0 {
		fmt.Printf("%s\n", header)
	}
	fmt.Printf("%s%s%s\n", Color.Green, strings.Repeat("-", width), Color.SysReset)
	fmt.Printf(strTable)
	fmt.Printf("%s%s%s\n", Color.Green, strings.Repeat("-", width), Color.SysReset)
}
