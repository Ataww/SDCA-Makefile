package main

import (
	"fmt"
)

type Target struct {
	id           string
	program      string
	args         string
	done         bool
	dependencies []*Target
}

func NewTarget(id string, program string, args string) *Target {
	target := new(Target)
	target.id = id
	target.program = program
	target.args = args
	target.done = false
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
	fmt.Print(t.id + " Dependencies : ")
	for i := 0; i < len(t.dependencies); i++ {
		t.dependencies[i].Print(level + 1)
	}
}

func (t *Target) Is_Computable() bool {
	for i := 0; i < len(t.dependencies); i++ {
		if t.dependencies[i].done != true {
			return false
		}
	}
	return true
}

func (t *Target) Get_Leaf() *Target {
	if t.Is_Computable() == true {
		//Compute
		t.done = true
		return t
	} else {
		//I'm not computable --> I have some "not-done" dependencies
		for i := 0; i < len(t.dependencies); i++ {
			//if t.dependencies[i].Is_Computable() == false {

			return t.dependencies[i].Get_Leaf()
			//}
		}
	}
	return nil
}

func main() {
	t1 := NewTarget("target1", "ls", "-l")
	t2 := NewTarget("target2", "ls", "-l")
	t3 := NewTarget("target3", "ls", "-l")
	t4 := NewTarget("target4", "ls", "-l")
	t5 := NewTarget("target5", "ls", "-l")
	t6 := NewTarget("target6", "ls", "-l")
	t7 := NewTarget("target7", "ls", "-l")

	t1.Add_Dependency(t2)
	t1.Add_Dependency(t3)

	t2.Add_Dependency(t4)
	t2.Add_Dependency(t5)

	t3.Add_Dependency(t6)
	t6.Add_Dependency(t7)
	//t1.Print(0)

	fmt.Printf("t1 : %t\n", t1.Is_Computable())
	fmt.Printf("t2 : %t\n", t2.Is_Computable())
	fmt.Printf("t3 : %t\n", t3.Is_Computable())
	fmt.Printf("t4 : %t\n", t4.Is_Computable())
	fmt.Printf("t5 : %t\n", t5.Is_Computable())
	fmt.Printf("t6 : %t\n", t6.Is_Computable())
	fmt.Printf("t7 : %t\n", t7.Is_Computable())

	for {
		var leaf = t1.Get_Leaf()
		if leaf == nil {
			fmt.Println("NULL")
		}
		leaf.Print(0)
	}
}
