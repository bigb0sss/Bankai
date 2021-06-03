package Read

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func readFile(inputFile string) string {

	// Write hexdump file from binary file (.bin)
	dumpFile := "output/shellcode.hexdump"

	f, err := os.Create(dumpFile)
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
	}
	defer f.Close()

	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
	}

	binToHex := hex.Dump(content)
	f.WriteString(binToHex)

	// Read & Parse shellcode
	file, err := os.Open(dumpFile)

	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	shellcode := ""
	for _, eachline := range txtlines {
		column := eachline[10:58] // Stupid way to parse hexdump
		noSpace := strings.ReplaceAll(column, " ", "")
		noNewline := strings.TrimSuffix(noSpace, "\n")
		shellcode += noNewline
	}

	return shellcode
}
