package main

import (
  "flag"
  "log"
  "bufio"
  "os"
  "fmt"
  "regexp"
)

func parse(path string) ([]string) {
  lines, _ := readFile(path)
  rawTargets := parseRawTargets(lines)
  return lines
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

func parseRawTargets(lines []string)(*map[string][]string) {
  rawTargets := make(map[string][]string)
  var (
    tar = regexp.MustCompile("(?P<target>[[:print:]]+): (?P<deps>[[:print:]]*)")
    command = regexp.MustCompile("(?P<cmds>)")
    )
    for i:= range lines {
      if(tar.MatchString(lines[i])) {
        if(command.MatchString(lines[i+1])) {

        }
      }
    }
  return &rawTargets
}

func main() {
  file := flag.String("path", "Makefile", "Specify the Makefile path.")
  flag.Parse()
  fmt.Println("Opening Makefile",*file)
  lines := parse(*file)
  for i := range lines {
    fmt.Println(i, lines[i])
  }

  re := regexp.MustCompile("(?P<first>[a-zA-Z]+) (?P<last>[a-zA-Z]+)")
	fmt.Println(re.MatchString("Alan Turing"))
	fmt.Printf("%q\n", re.SubexpNames())
	reversed := fmt.Sprintf("${%s} ${%s}", re.SubexpNames()[2], re.SubexpNames()[1])
	fmt.Println(reversed)
	fmt.Println(re.ReplaceAllString("Alan Turing", reversed))
}
