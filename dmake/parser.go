package main

import (
  "bufio"
  "os"
  "regexp"
  "strings"
)

// struct containing a target parsed after the first pass
type RawTarget struct {
  deps, cmds string
}

// Map containing the targets parsed after the first pass
type RawTargetMap map[string]*RawTarget

// Map containing the targets contained in the final dependency tree
type TargetMap map[string]*Target

/* Parse the makefile at the given path into a tree of Targets.
* path: Path to the makefile to parse.
*
* return: The root of the dependency tree.
* return: error, if any.
*/
func Parse(path string) (*Target, error) {
  lines, err := readFile(path)
  if err != nil {
    return nil, err
  }
  // first pass
  raws, root := parseRawTargets(lines)
  targetBuilt := make(TargetMap)
  // second pass
  targetTree := builDependencyTree(raws,root, targetBuilt)
  return targetTree, err
}

/* Read the file at the given path and return its content line by line as an array of string.
* path: Path to the file to read
*
* return: the file's content as an array of string.
* return: error, if any.
*/
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

/* parse the given string array and return a map with all the targets and the name of the root target
* lines: le contenu du Makefile en tableau de string
*
* return: A map of targets crudely parsed
* return: The name of the root target.
*/
func parseRawTargets(lines []string)(*RawTargetMap,string) {
  var (
    // regexp for the target first line
    tar = regexp.MustCompile("([\\w\\.]+)\\s?:(?:\\s[\\w\\.]+)*")
    // regexp for the targets' command
    command = regexp.MustCompile("[[:print]]*")
    rawTargets = make(RawTargetMap)
    root = ""
    )
    for i := range lines {
      // check if the current line match the target pattern
      if tar.MatchString(lines[i]) {
        targets:= strings.Split(lines[i],":")
        cur := strings.TrimSpace(targets[0])
        // create a new RawTarget element mapped to the target's name
        rawTargets[cur] = &RawTarget{}
        if root == "" {
          root = cur
        }
        // Add dependencies if there's any
        if targets[1] != "" {
          rawTargets[cur].deps = strings.TrimSpace(targets[1])
        }
        // Add command if there's any (following line)
        if command.MatchString(lines[i+1]) {
          rawTargets[cur].cmds = strings.TrimSpace(lines[i+1])
        }
      }
    }
  return &rawTargets, root
}

/* Recursively build a dependency tree from the targets parsed in the first pass.
* raws: pointer to the map of targets parsed in the first pass
* root: name of the current target to build (root target in the first call)
* targetBuilt: Map containing the already built targets (to ensure node uniqueness)
*
* return: A pointer to the root of the currently built subtree
*/
func builDependencyTree(raws *RawTargetMap, root string, targetBuilt TargetMap)(*Target) {
  raw, ok := (*raws)[root]
  // some targets are referenced but not defined in the makefile. In that case we create a reference to the external target.
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
      // retrieve an already built target or build it
      if !ok {
        dep = builDependencyTree(raws,deps[i],targetBuilt)
        targetBuilt[deps[i]] = dep
      }
      // add the dependency
      target.dependencies = append(target.dependencies,dep, targetBuilt[deps[i]])
    }
  }
  // complete subtree (or tree)
  return target
}
