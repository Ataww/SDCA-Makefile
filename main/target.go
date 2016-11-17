package main

import (
	"fmt"
)

type Target struct {
	id           string
	program      string
	args         string
	serverId     int
	done         bool
	computing    bool
	dependencies []*Target
}

func NewTarget(id string, program string, args string) *Target {
	target := new(Target)
	target.id = id
	target.program = program
	target.args = args
	target.done = false
	target.computing = false
	return target
}

func (t *Target) Add_Dependency(dependency *Target) {
	t.dependencies = append(t.dependencies, dependency)
}

func (t Target) Print(level int) {
	fmt.Print("\n")
	for i := 0; i < level; i++ {
		fmt.Print("\t")
	}
	fmt.Print(t.id + " Args : ",t.args,"\n")
	fmt.Print(" Dependencies : \n")
	for i := 0; i < len(t.dependencies); i++ {
		t.dependencies[i].Print(level + 1)
	}
}

func (t *Target) Is_Computable() bool {
	if t.computing == true {
		return false
	}

	if t.done == true {
		return false
	} else {
		for i := 0; i < len(t.dependencies); i++ {
			if t.dependencies[i].done != true {
				return false
			}
		}
		return true
	}
}

func (t *Target) Get_Leaf() *Target {
	if t.Is_Computable() {
		return t
	} else {
		for i := 0; i < len(t.dependencies); i++ {
			if t.dependencies[i].done == false && t.dependencies[i].computing == false {
				return t.dependencies[i].Get_Leaf()
			}
		}

		return nil
	}
}
