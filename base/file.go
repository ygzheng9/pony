package base

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// PrettyPrint 使用缩进打印 struct
func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("pretty error: %+v\n", err)
		return
	}

	fmt.Println(string(b))
}

// ReadLines 逐行读取文件
func ReadLines(path string) ([]string, error) {
	sugar := Sugar()

	file, err := os.Open(path)
	if err != nil {
		sugar.Error("can not open file", "file", path)
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		lines = append(lines, trimmed)
	}
	return lines, scanner.Err()
}

// WriteLines writes the lines to the given file.
func WriteLines(lines []string, path string) error {
	sugar := Sugar()

	file, err := os.Create(path)
	if err != nil {
		sugar.Error("can not create file", "file", path)
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
