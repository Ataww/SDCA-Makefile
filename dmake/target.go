package main

import (
	"fmt"
)

/*
* Struct representing a Target in the makefile
* id: target's name
* lineCommand: Target's command
* serverId: Id of the server the target is being assigned to
* done: flag indicating whether the target has been executed
* computing: flag indicating whether the current target is being computed
* dependencies: slice containing pointers to the target's dependencies
*/
type Target struct {
	id           string
	lineCommand      string
	serverId     int
	done         bool
	computing    bool
	dependencies []*Target
}

/*
* Create a new instance of Target
* id: The target's name
* _lineCommand: The target's command
*
* return: a pointer to the newly created Target
*/
func NewTarget(id string, _lineCommand string) *Target {
	target := new(Target)
	target.id = id
	target.lineCommand = _lineCommand
	target.done = false
	target.computing = false
	return target
}

/*
* Add a dependency to the calling target.
* dependency: Pointer to the dependency to add.
*/
func (t *Target) Add_Dependency(dependency *Target) {
	t.dependencies = append(t.dependencies, dependency)
}

/*
* Pretty print the target at the given level (number of tabs)
* level: level in the tree
*/
func (t Target) Print(level int) {
	fmt.Print("\n")
	for i := 0; i < level; i++ {
		fmt.Print("\t")
	}
	//fmt.Print(t.id + " Args : ",t.args,"\n")
	fmt.Print(" Dependencies : \n")
	for i := 0; i < len(t.dependencies); i++ {
		t.dependencies[i].Print(level + 1)
	}
}

/*
* Check if the target is ready to be computed.
* A target is ready if it is neither being computed nor done and if all its dependencies have been computed already.
*/
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

/*
* Recursively check if the current leaf is ready for computation. Return the next available target, or nil if none can be found.
* return: a pointer to the next target available for computation.
*/
func (t *Target) Get_Leaf() *Target {
	if t.Is_Computable() {
		return t
	} else {
		for i := 0; i < len(t.dependencies); i++ {
			if t.dependencies[i].done == false && t.dependencies[i].computing == false {
				if leaf := t.dependencies[i].Get_Leaf(); leaf!= nil{
					return leaf
				}

			}
		}
		return nil
	}
}
