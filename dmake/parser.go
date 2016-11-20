package main

import (
  "bufio"
  "os"
  "regexp"
  "strings"
)

type RawTarget struct {
  deps, cmds string
}

type RawTargetMap map[string]*RawTarget

type TargetMap map[string]*Target

func Parse(path string) (*Target, error) {
  lines, err := readFile(path)
  if err != nil {
    return nil, err
  }
  raws, root := parseRawTargets(lines)
  targetBuilt := make(TargetMap)
  targetTree := builDependencyTree(raws,root, targetBuilt)
  return targetTree, err
}

func readFile(path string)([]string, error) {
  file, err := os.Open(path)
  if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lines := make([]string,0,8)
    for  scanner.Scan() {
        lines = append(lines,scanner.Text())
    }
    return lines, err
}

func parseRawTargets(lines []string)(*RawTargetMap,string) {
  var (
    tar = regexp.MustCompile("([\\w\\.]+)\\s?:(?:\\s[\\w\\.]+)*")
    command = regexp.MustCompile("[[:print]]*")
    rawTargets = make(RawTargetMap)
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

func builDependencyTree(raws *RawTargetMap, root string, targetBuilt TargetMap)(*Target) {
  // cmds
  raw, ok := (*raws)[root]
  if !ok {
    target := NewTarget(root,"")
    targetBuilt[root] = target
    return target
  }
  target := NewTarget(root, raw.cmds)
  // Deps
  deps := strings.Split(raw.deps, " ")
  for i := range deps {
    if deps[i] != "" {
      dep, ok := targetBuilt[deps[i]]
      if !ok {
        dep = builDependencyTree(raws,deps[i],targetBuilt)
        targetBuilt[deps[i]] = dep
      }
      target.dependencies = append(target.dependencies,dep, targetBuilt[deps[i]])
    }
  }
  // tree
  return target
}
