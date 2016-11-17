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
	fmt.Print("Executing target ", command.ID, " : ", command.Program, " ", command.Arguments, "\n")

	bash_command := command.Program+" "+ command.Arguments

	// Create command
	cmd := exec.Command("bash", "-c",bash_command)
	out, err := cmd.Output()

	if (err != nil){
		fmt.Print("Command executed with errors :",err.Error()," \n")
		fmt.Print("output :",string(out[:])," \n")
		return -1, err
	}else{
		fmt.Print("Command executed without errors :",string(out[:])," \n")
		return 0, nil
	}
}
