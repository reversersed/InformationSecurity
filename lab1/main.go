package main

import "fmt"

func main() {
	random := NewRandom()

	count := make([]int, 0)
	fmt.Println("Числа:")
	for i := 0; i < m; i++ {
		count = append(count, random.Next())
		fmt.Printf("%d\t", count[i])
	}

	fmt.Println("\n\nПропущенные числа:")
	for i := 0; i < m; i++ {
		if !Search(count, i) {
			fmt.Printf("%d\t", i)
		}
	}

	fmt.Println("\n\nПовторяющиеся числа:")
	repeated := make([]bool, len(count))
	for i := 0; i < m; i++ {
		if repeated[count[i]] {
			fmt.Printf("%d\t", count[i])
		}
		repeated[count[i]] = true
	}
}
func Search(slice []int, value int) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == value {
			return true
		}
	}
	return false
}
