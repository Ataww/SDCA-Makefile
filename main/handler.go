package main

import (
	"fmt"
	//"io"
	"SDCA-Makefile/compilationInterface"
	"os/exec"
)

type CompilationHandler struct {
}

func NewCompilationHandler() *CompilationHandler {
	return &CompilationHandler{}
}

func (p *CompilationHandler) ExecuteCommand(command *compilationInterface.Command) (status compilationInterface.Int, err error) {
	fmt.Print("Execute command ->", command.Program, " ", command.Arguments, "\n")

	// Create command
	cmd := exec.Command(command.Program, command.Arguments)
	/*stdErrPipe, error := cmd.StderrPipe()
	if error != nil {
		fmt.Print("ExecuteCommand() an error occureds\n")
		return 1, error
	}
	stdOutPipe, error := cmd.StdoutPipe()
	if error != nil {
		fmt.Print("ExecuteCommand() an error occureds\n")
		return 1, error
	}*/

	// Run command
	error := cmd.Start()
	if error != nil {
		fmt.Print("ExecuteCommand() an error occureds.\n")
		return 1, error
	}

	// Command success
	fmt.Print("Command xexcecute without errors \n")
	return 0, error
}
