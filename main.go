package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	// "os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
)

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func main() {
	const dumpName = ".lbm_dump"
	tw := new(tabwriter.Writer)

	command, argument := "", ""
	commands := []string{"add", "get", "rm", "ls"}

	rawArgs := os.Args

	command = rawArgs[1]

	if len(rawArgs) > 2 {
		argument = rawArgs[2]
	}

	// check if the dump file exists
	_, err := os.Stat(dumpName)
	if os.IsNotExist(err) {
		// create the dump file
		os.Create(dumpName)
	} else {
		// Stat() err
	}

	// check for command
	// Need to check if the command is in commands
	if contains(commands, command) == false {
		fmt.Println("invalid command")
		return
	}

	// Read file and split with new line
	readDat, err := ioutil.ReadFile(dumpName)
	readDatStr := string(readDat)
	lineSplitted := strings.Split(readDatStr, "\n")

	// map via number and map via name
	mapNumber := make(map[string]string)

	lineStrSplitted := make([]string, 2)
	lbmNumber, lbmDir := "", ""
	lastLineNumberStr := ""

	if lineSplitted[0] != "" {
		for _, lineStr := range lineSplitted {
			if lineStr != "" {
				lineStrSplitted = strings.Split(lineStr, `,`)

				lbmNumber = lineStrSplitted[0]
				lbmDir = lineStrSplitted[1]

				mapNumber[lbmNumber] = lbmDir
				lastLineNumberStr = lbmNumber
			}
		}
	}

	lastLineNumber, _ := strconv.Atoi(lastLineNumberStr)
	currentNumberStr := strconv.Itoa(lastLineNumber + 1)

	if command == "add" {
		var buffer bytes.Buffer
		buffer.WriteString(currentNumberStr)
		buffer.WriteString(",")
		argument = processDir(argument)
		buffer.WriteString(argument)
		buffer.WriteString("\n")
		added := []byte(buffer.String())

		f, _ := os.OpenFile(dumpName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.Write(added)
		f.Close()
	} else if command == "get" {
		// cmd := exec.Command("bash", "-c", "cd", mapNumber[argument])
		// err := cmd.Run()
		// os.Chdir(mapNumber[argument])
		// sh.Command(sh.Dir(mapNumber[argument])).Run()
		fmt.Println(mapNumber[argument])
	} else if command == "rm" {
		var buffer bytes.Buffer

		if _, ok := mapNumber[argument]; !ok {
			return
		}

		delete(mapNumber, argument)
		os.Remove(dumpName)
		for tempNumber, tempDir := range mapNumber {
			buffer.WriteString(tempNumber)
			buffer.WriteString(",")
			buffer.WriteString(tempDir)
			buffer.WriteString("\n")
		}

		added := []byte(buffer.String())

		f, _ := os.OpenFile(dumpName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.Write(added)
		f.Close()
	} else if command == "ls" {

		if lineSplitted[0] == "" {
			fmt.Println("Empty")
			return
		}

		tw.Init(os.Stdout, 0, 0, 5, ' ', tabwriter.AlignRight)
		fmt.Fprintln(tw, "id\tdirectory\t")

		for tempNumber, tempDir := range mapNumber {
			toJoin := []string{tempNumber, tempDir, ""}
			fmt.Fprintln(tw, strings.Join(toJoin, "\t"))
		}

		tw.Flush()
	}
}

func processDir(dir string) string {
	firstChar := string(dir[0])

	if firstChar == "." && len(dir) == 1 {
		return currentPath()
	}

	if firstChar != "/" && firstChar != "~" {
		return processAbsolutePath(dir)
	}

	return dir
}

func currentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}

func processAbsolutePath(dir string) string {
	var buffer bytes.Buffer
	buffer.WriteString(currentPath())
	buffer.WriteString("/")
	buffer.WriteString(dir)
	return buffer.String()
}
