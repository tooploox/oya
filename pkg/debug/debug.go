// +build debug

package debug

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"

	logrus "github.com/sirupsen/logrus"
)

func FP(cond bool) bool {
	_, sourceFile, lineNumber, ok := runtime.Caller(1)
	if !ok {
		panic("Internal error: problem getting caller from runtime")
	}
	point := FailurePoint{
		sourceFile: sourceFile,
		lineNumber: lineNumber,
	}
	_, ok = failurePoints[point]
	failure := ok && cond
	if failure {
		logrus.Println("Triggering error")
		debug.PrintStack()
	}
	return failure
}

type FailurePoint struct {
	sourceFile string
	lineNumber int
}

var failurePoints = make(map[FailurePoint]struct{})

func init() {
	fname := "failures.txt"
	file, err := os.Open(fname)
	if err != nil {
		logrus.Println("No failures.txt!")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			panic(fmt.Sprintf("Expected lines matching <path>:<lineNumber> format in %s; actual: %q",
				fname, line))
		}
		path := parts[0]
		lineNumber, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(fmt.Sprintf("Expected line number to be an integer; actual: %q",
				parts[1]))
		}
		sourceFile, err := filepath.Abs(path)
		if err != nil {
			panic(err)
		}
		point := FailurePoint{
			sourceFile: sourceFile,
			lineNumber: lineNumber,
		}
		failurePoints[point] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
