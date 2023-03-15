package main

import "fmt"
// O(n)

func checkIfBalanced(str string) bool {
	result := true
	temp := make([]rune, 0, len(str))
	for _, val := range str {
		if val == '(' {
			temp = append(temp, val)
		} else if val == ')' {
			if len(temp) == 0 {
				result = false
				break
			}
			temp = temp[:len(temp)-1]
		}
	}

	if len(temp) != 0 {
		result = false
	}

	return result
}

func main() {
	fmt.Println(checkIfBalanced(")()()("))
}
