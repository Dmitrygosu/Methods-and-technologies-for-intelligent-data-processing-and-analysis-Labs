package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	m1 := readFile("olymp1.txt")
	m2 := readFile("olymp2.txt")
	all := make(map[string]bool)
	for k := range m1 {
		all[k] = true
	}
	for k := range m2 {
		all[k] = true
	}

	var names []string
	for k := range all {
		names = append(names, k)
	}
	sort.Strings(names)

	for _, name := range names {
		in1, ok1 := m1[name]
		in2, ok2 := m2[name]

		res := ""


		if (ok1 && !ok2) || (ok2 && !ok1) {
			res = "д"
		} else {
			score := 0
			if in1 { score++ }
			if in2 { score++ }

			if score == 0 {
				res = "3"
			} else if score == 1 {
				res = "4"
			} else {
				res = "5"
			}
		}
		
		fmt.Printf("%s %s\n", name, res)
	}
}

func readFile(name string) map[string]bool {
	m := make(map[string]bool)
	
	f, err := os.Open(name)
	if err != nil {
		fmt.Println("Err:", err)
		return m
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 { continue }

		fam := parts[0]
		prize := false

		if len(parts) > 1 {
			p := parts[1]
			if p == "1" || p == "2" || p == "3" {
				prize = true
			}
		}
		if val, ok := m[fam]; ok {
			if val { prize = true }
		}
		m[fam] = prize
	}
	return m
}
