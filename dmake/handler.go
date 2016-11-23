package main

import (
	"SDCA-Makefile/compilationInterface"
	"fmt"
	"os"
	"os/exec"
)

type CompilationHandler struct {
}

func NewCompilationHandler() *CompilationHandler {
	return &CompilationHandler{}
}

/*
Execute a command
*/
func (p *CompilationHandler) ExecuteCommand(command *compilationInterface.Command) (status compilationInterface.Int, err error) {
	fmt.Println("Going to execute ", command.ID, " : ", command.CommandLine)

	// Create command
	cmd := exec.Command("bash", "-c", command.CommandLine)
	cmd.Dir = command.WorkingDir
	out, err := cmd.Output()

	if err != nil {
		fmt.Print("Command executed with errors :", err.Error(), " \n")
		fmt.Print("output :", string(out[:]), " \n")
		return -1, err
	} else {
		fmt.Print("Command executed without errors :", string(out[:]), " \n")
		return 0, nil
	}
}

func (p *CompilationHandler) Stop() (err error) {
	fmt.Println("Server is going to shudown")
	defer os.Exit(0)
	return nil
}
