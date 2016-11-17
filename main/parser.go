package main

import (
  "log"
  "bufio"
  "os"
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
    target := NewTarget(root,"")
    return target
  }
  target := NewTarget(root, red.cmds)
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
