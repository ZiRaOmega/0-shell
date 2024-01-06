package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var PATH string = getCurrentDirectory()

func main() {

	/* //alternative of Scanln

	reader := bufio.NewReader(os.Stdin)

	test, _ := reader.ReadString('\n')

	fmt.Println(test)
	return */
	for {
		handle_user_input()
	}
}

func handle_user_input() {
	//print a simple $ and wait for user input
	fmt.Print(PATH + " $ ")
	var input string
	reader := bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')

	//fmt.Println(input)
	/*echo
	cd
	ls
	pwd
	cat
	cp
	rm
	mv
	mkdir
	exit*/
	commands := []string{"echo", "cd", "ls", "pwd", "cat", "cp", "rm", "mv", "mkdir", "exit"}
	input_command := ""
	if strings.Contains(input, " ") {
		input_command = strings.Split(input, " ")[0]
	} else {
		input_command = input

	}
	//trim newline
	input_command = strings.Trim(input_command, "\n")
	//fmt.Println(input_command)
	for i := 0; i < len(commands); i++ {
		if input_command == commands[i] {
			args := split_input(input)
			switch input_command {
			case "echo":
				//fmt.Println("echo")
				echo(strings.Join(args, " "))
				break
			case "cd":
				//fmt.Println("cd")
				cd(args[0])
				break
			case "ls":
				//fmt.Println("ls")
				ls(args[0], args[1:])
				break
			case "pwd":
				//fmt.Println("pwd")
				pwd()
				break
			case "cat":
				//fmt.Println("cat")
				cat(args[0])
				break
			case "cp":
				//fmt.Println("cp")
				cp(args[0], args[1], args[2:])
				break
			case "rm":
				rm(args[0], args[1:])
				break
			case "mv":
				//fmt.Println("mv")
				mv(args[0], args[1], args[2:])
				break
			case "mkdir":
				//fmt.Println("mkdir")
				mkdir(args[0])
				break
			case "exit":
				//fmt.Println("exit")
				exit()
				break
			}
		}
	}
}
func echo(s string) {
	fmt.Println(s)
}

func cd(s string) {
	//check if the new path exist
	//if it does change the path
	//if it doesn't print an error
	var err error
	if s == "~" {
		err = os.Chdir(PATH)
		PATH = getCurrentDirectory()
	} else {
		err = os.Chdir(s)
		PATH = getCurrentDirectory()
	}
	if err != nil {
		Error_Handler(err)
	}

}

// ls function with flag handling
func ls(path string, flags []string) {
	// Default to current directory if no path is specified
	if path == "" {
		path = PATH
	} else {
		fullPath := PATH + "/" + path
		if IsPathExist(fullPath) {
			path = fullPath
		} else {
			Error_Handler(fmt.Errorf("Path '%s' does not exist", path))
			return
		}
	}

	// Append flags to the command arguments
	cmdArgs := append([]string{path}, flags...)
	cmd := exec.Command("ls", cmdArgs...)

	// Execute the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		Error_Handler(fmt.Errorf("ls command failed: %s", err))
		return
	}
	fmt.Println(string(output))
}

// cp function with flag handling
func cp(src string, dest string, flags []string) {
	// Verify source file existence
	srcFullPath := PATH + "/" + src
	if !IsPathExist(srcFullPath) {
		Error_Handler(fmt.Errorf("Source file '%s' does not exist", src))
		return
	}

	// Construct destination path
	destFullPath := PATH + "/" + dest

	// Append flags to the command arguments
	cmdArgs := append([]string{srcFullPath, destFullPath}, flags...)
	cmd := exec.Command("cp", cmdArgs...)

	// Execute the command
	err := cmd.Run()
	if err != nil {
		Error_Handler(fmt.Errorf("cp command failed: %s", err))
		return
	}
	fmt.Println("File copied successfully")
}

// rm function with flag handling
func rm(target string, flags []string) {
	// Verify file or directory existence
	targetFullPath := PATH + "/" + target
	if !IsPathExist(targetFullPath) {
		Error_Handler(fmt.Errorf("Target '%s' does not exist", target))
		return
	}

	// Append flags to the command arguments
	cmdArgs := append([]string{targetFullPath}, flags...)
	cmd := exec.Command("rm", cmdArgs...)

	// Execute the command
	err := cmd.Run()
	if err != nil {
		Error_Handler(fmt.Errorf("rm command failed: %s", err))
		return
	}
	fmt.Println("Item removed successfully")
}

// mv function with flag handling
func mv(src string, dest string, flags []string) {
	// Verify source file or directory existence
	srcFullPath := PATH + "/" + src
	if !IsPathExist(srcFullPath) {
		Error_Handler(fmt.Errorf("Source '%s' does not exist", src))
		return
	}

	// Construct destination path
	destFullPath := PATH + "/" + dest

	// Append flags to the command arguments
	cmdArgs := append([]string{srcFullPath, destFullPath}, flags...)
	cmd := exec.Command("mv", cmdArgs...)

	// Execute the command
	err := cmd.Run()
	if err != nil {
		Error_Handler(fmt.Errorf("mv command failed: %s", err))
		return
	}
	fmt.Println("Item moved successfully")
}

// [Other functions like mkdir, cat, pwd, etc., remain unchanged]

func pwd() {
	fmt.Println(PATH)
}

func cat(s string) {
	//check if the file exist
	//if it does print the content of the file
	//if it doesn't print an error

	if IsPathExist(PATH + "/" + s) {
		s = PATH + "/" + s
		file, _ := os.OpenFile(s, os.O_RDONLY, 0644)
		var content []byte
		_, _ = file.Read(content)
		fmt.Println(content)
	} else {
		Error_Handler(fmt.Errorf("File doesn't exist"))
	}
}

func mkdir(dir string) {
	//check if the directory exist
	//if it does print an error
	//if it doesn't create the directory
	if IsPathExist(PATH + "/" + dir) {
		Error_Handler(fmt.Errorf("Directory already exist"))
	} else {
		dir = PATH + "/" + dir
		createDir(dir)
	}
}
func exit() {
	os.Exit(0)
}

func createDir(dirName string) bool {
	//check if the directory exist
	//if it does print an error
	//if it doesn't create the directory
	if IsPathExist(PATH + "/" + dirName) {
		Error_Handler(fmt.Errorf("Directory already exist"))
		return false
	}
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		Error_Handler(err)
		return false
	}
	return true
}
func getCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		Error_Handler(err)
	}
	return dir
}

func IsPathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func Error_Handler(err_message error) {
	//if type of err_message

	//Print error has error
	fmt.Println("Error has occured: " + err_message.Error())
}

func split_input(input string) []string {
	input = strings.Trim(input, "\n")
	if strings.Contains(input, " ") {
		return strings.Split(input, " ")[1:]
	} else {
		return []string{""}

	}
}
