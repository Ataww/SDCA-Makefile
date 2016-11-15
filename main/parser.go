package main

import (
  "fmt"
  "bufio"
  "os"
)

type Target struct {
  name string
  commands []string
  dependencies []string
}

func parse(path string) {
  lines, err = readFile(path)
}

func readFile(path string)([]string, error) {
  file, err := os.Open(path)
  if err != nil {
        log.Fatal(err)
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lines = make([]string,0,8)
    for i := 0; scanner.Scan(); i++ {
        lines[i] = scanner.Text()
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return lines, err
}
