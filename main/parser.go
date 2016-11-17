package main

import (
  "flag"
  "log"
  "bufio"
  "os"
  "fmt"
  "regexp"
  "strings"
)

type RawTarget struct {
  deps, cmds string
}

type TargetMap map[string]*RawTarget

func Parse(path string) (*Target, error) {
  lines, err := readFile(path)
  if err != nil {
    return nil, err
  }
  raws, root := parseRawTargets(lines)
  targetTree := builDependencyTree(raws,root)
  return targetTree, err
}

func readFile(path string)([]string, error) {
  file, err := os.Open(path)
  if err != nil {
        log.Fatal(err)
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lines := make([]string,0,8)
    for  scanner.Scan() {
        lines = append(lines,scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return lines, err
}

func parseRawTargets(lines []string)(*TargetMap,string) {
  var (
    tar = regexp.MustCompile("([\\w\\.]+)\\s?:(?:\\s[\\w\\.]+)*")
    command = regexp.MustCompile("[[:print]]*")
    rawTargets = make(TargetMap)
    root = ""
    )
    for i := range lines {
      if tar.MatchString(lines[i]) {
        targets:= strings.Split(lines[i],":")
        cur := strings.TrimSpace(targets[0])
        rawTargets[cur] = &RawTarget{}
        if root == "" {
          root = cur
        }
        if targets[1] != "" {
          rawTargets[cur].deps = strings.TrimSpace(targets[1])
        }
        if command.MatchString(lines[i+1]) {
          rawTargets[cur].cmds = strings.TrimSpace(lines[i+1])
        }
      }
    }
  return &rawTargets, root
}

func builDependencyTree(raws *TargetMap, root string)(*Target) {
  // cmds
  red, ok := (*raws)[root]
  if !ok {
    target := NewTarget(root,"","")
    return target
  }
  cmd, args := parseCmd(red.cmds)
  target := NewTarget(root, cmd,args)
  // Deps
  deps := strings.Split(red.deps, " ")
  for i := range deps {
    if deps[i] != "" {
      dep := builDependencyTree(raws,deps[i])
      target.dependencies = append(target.dependencies,dep)
    }
  }
  // tree
  return target
}

// split the command string between command and arguments
func parseCmd(cmd string)(string,string) {
  blocks := strings.Split(cmd," ")
  command := blocks[0]
  blocks = append(blocks[:0], blocks[1:]...)
  args := strings.Join(blocks," ")
  return command, args
}

func mainaze() {
  file := flag.String("makefile", "Makefile", "Specify the Makefile path.")
  flag.Parse()
  fmt.Println("Opening Makefile",*file)
  target, _ := Parse(*file)
  fmt.Println("parsed",*target)
}
