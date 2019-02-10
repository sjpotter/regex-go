package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sjpotter/regex-go/pkg/regex"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		r := scanner.Text()
		scanner.Scan()
		text := scanner.Text()

		t := regex.NewTokenizer(r)
		m, err := regex.NewMatcher(t)
		if err != nil {
			fmt.Printf("matcher failed: %v\n", err)
			continue
		}

		result, err := m.Match(text)
		if err != nil {
			fmt.Sprintf("regex library failed: %v", err)
		}
		if result {
			fmt.Println(text + " matched regex " + r)
			fmt.Printf("Groups:")
			for i, g := range m.GetGroups() {
				fmt.Printf("%v: %v\n", i, *g)
			}
		} else {
			fmt.Println(text + " didn't match regex " + r)

		}


	}
}
