package main

import (
	"flag"
	"fmt"
)

func main() {
	firstString := flag.String("original", "not", "The original string")
	secondString := flag.String("compare", "working", "The string to compare with")
	recursive := flag.Bool("recursive", true, "recursive computation of levensthein distance")
	damerau := flag.Bool("damerau", false, "Enables Damerau-Levensthein Distance")
	flag.Parse()
	if *damerau {
		if *recursive {
			fmt.Printf("The Damerau-Leventhein Edit Distance between %[1]s and %[2]s is %[3]d operations.", *firstString, *secondString, DamerauLevenstheinDistanceRecursive(*firstString, *secondString))
		}
	}
	if *recursive {
		fmt.Printf("The Leventhein Edit Distance between %[1]s and %[2]s is %[3]d operations.", *firstString, *secondString, LevenstheinDistanceRecursive(*firstString, *secondString))
	} else {
		fmt.Printf("The Leventhein Edit Distance between %[1]s and %[2]s is %[3]d operations.", *firstString, *secondString, LevenstheinDistanceRecursive(*firstString, *secondString))
	}
}

func LevenstheinDistanceDynamic(first, second string) int {
	// y-dimension of the table
	table := make([][]int, len(first)+1)
	// x-dimension of the table
	for x := range table {
		table[x] = make([]int, len(second)+1)
	}
	// initialiyze the column representing the original string
	for i := 0; i <= len(first); i++ {
		table[i][0] = i
	}
	// initialize the row representing the compare string
	for j := 0; j < len(second); j++ {
		table[0][j] = j
	}
	// iterate through the table row by row
	for j := 1; j <= len(second); j++ {
		for i := 1; i <= len(first); i++ {
			if first[i-1] == second[j-1] {
				table[i][j] = minimum(table[i-1][j]+1, table[i][j-1]+1, table[i-1][j-1])
			} else {
				table[i][j] = minimum(table[i-1][j]+1, table[i][j-1]+1, table[i-1][j-1]+1)
			}
		}
	}
	return table[len(first)][len(second)]
}

func DamerauLevenstheinDistanceDynamic(first, second string) int {
	// y-dimension of the table
	table := make([][]int, len(first)+1)
	// x-dimension of the table
	for x := range table {
		table[x] = make([]int, len(second)+1)
	}
	// initialiyze the column representing the original string
	for i := 0; i <= len(first); i++ {
		table[i][0] = i
	}
	// initialize the row representing the compare string
	for j := 0; j < len(second); j++ {
		table[0][j] = j
	}
	// iterate through the table row by row
	for j := 1; j <= len(second); j++ {
		for i := 1; i <= len(first); i++ {
			if i > 1 && j > 1 && first[i] == second[j-1] && first[i-1] == second[j] {
				if first[i-1] == second[j-1] {
					table[i][j] = minimum(table[i-1][j]+1, table[i][j-1]+1, table[i-1][j-1]+1, table[i-2][j-2]+1)
				} else {
					table[i][j] = minimum(table[i-1][j]+1, table[i][j-1]+1, table[i-1][j-1], table[i-2][j-2]+1)
				}
			} else {
				if first[i-1] == second[j-1] {
					table[i][j] = minimum(table[i-1][j]+1, table[i][j-1]+1, table[i-1][j-1])
				} else {
					table[i][j] = minimum(table[i-1][j]+1, table[i][j-1]+1, table[i-1][j-1]+1)
				}
			}
		}
	}
	return table[len(first)][len(second)]
}

func LevenstheinDistanceRecursive(first, second string) int {
	if first == "" {
		return len(second)
	}
	if second == "" {
		return len(first)
	}
	if first[0] == second[0] {
		return LevenstheinDistanceRecursive(first[1:], second[1:])
	}
	return minimum(LevenstheinDistanceRecursive(first[1:], second),
		LevenstheinDistanceRecursive(first, second[1:]),
		LevenstheinDistanceRecursive(first[1:], second[1:])) + 1

}

func DamerauLevenstheinDistanceRecursive(first, second string) int {
	if first == "" {
		return len(second)
	}
	if second == "" {
		return len(first)
	}
	if first[0] == second[0] {
		return LevenstheinDistanceRecursive(first[1:], second[1:])
	}
	return minimum(LevenstheinDistanceRecursive(first[1:], second),
		LevenstheinDistanceRecursive(first, second[1:]),
		LevenstheinDistanceRecursive(first[1:], second[1:])) + 1
}

func minimum(numbers ...int) int {
	min := numbers[0]
	for _, num := range numbers {
		if num < min {
			min = num
		}
	}
	return min
}

func maximum(numbers ...int) int {
	max := numbers[0]
	for _, num := range numbers {
		if num > max {
			max = num
		}
	}
	return max
}
