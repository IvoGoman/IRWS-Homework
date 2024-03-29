package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

var weightMap = loadWeights("./weights.csv")

func main() {
	firstString := flag.String("original", "not", "The original string")
	secondString := flag.String("compare", "working", "The string to compare with")
	recursive := flag.Bool("recursive", false, "recursive computation of levensthein distance")
	damerau := flag.Bool("damerau", true, "Enables Damerau-Levensthein Distance")
	weights := flag.Bool("weights", true, "Enables Damerau-Levensthein Distance with Weights")
	flag.Parse()
	if *damerau {
		if *recursive {
			fmt.Printf("The Damerau-Leventhein Edit Distance between %[1]s and %[2]s is %[3]d operations.", *firstString, *secondString, DamerauLevenstheinDistanceRecursive(*firstString, *secondString))
		} else if *weights {
			fmt.Printf("The Damerau-Leventhein Edit Distance between %[1]s and %[2]s is %[3]f operations with custom weight.", *firstString, *secondString, DamerauLevenstheinDistanceDynamicWeighted(*firstString, *secondString))
		} else {
			fmt.Printf("The Damerau-Leventhein Edit Distance between %[1]s and %[2]s is %[3]d operationss.", *firstString, *secondString, DamerauLevenstheinDistanceDynamic(*firstString, *secondString))
		}
	} else {
		if *recursive {
			fmt.Printf("The Leventhein Edit Distance between %[1]s and %[2]s is %[3]d operations.", *firstString, *secondString, LevenstheinDistanceRecursive(*firstString, *secondString))
		} else {
			fmt.Printf("The Leventhein Edit Distance between %[1]s and %[2]s is %[3]d operations.", *firstString, *secondString, LevenstheinDistanceRecursive(*firstString, *secondString))
		}
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
	if minimum(len(first), len(second)) == 0 {
		return maximum(len(first), len(second))
	}
	// y-dimension of the table
	table := make([][]int, len(first)+1)
	// x-dimension of the table
	for x := range table {
		table[x] = make([]int, len(second)+1)
	}
	// initialiyze the column representing the original string
	for i := 1; i <= len(first); i++ {
		table[i][0] = i
	}
	// initialize the row representing the compare string
	for j := 1; j <= len(second); j++ {
		table[0][j] = j
	}
	// iterate through the table row by row
	for i := 1; i <= len(first); i++ {
		for j := 1; j <= len(second); j++ {
			if i > 1 && j > 1 && first[i-1] == second[j-2] && first[i-2] == second[j-1] {
				if first[i-1] == second[j-1] {
					table[i][j] = minimum(table[i-1][j]+1, table[i][j-1]+1, table[i-1][j-1], table[i-2][j-2]+1)
				} else {
					table[i][j] = minimum(table[i-1][j]+1, table[i][j-1]+1, table[i-1][j-1]+1, table[i-2][j-2]+1)
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

func DamerauLevenstheinDistanceDynamicWeighted(first, second string) float64 {
	if minimum(len(first), len(second)) == 0 {
		return float64(maximum(len(first), len(second)))
	}
	// y-dimension of the table
	table := make([][]float64, len(first)+1)
	// x-dimension of the table
	for x := range table {
		table[x] = make([]float64, len(second)+1)
	}
	// initialize the column representing the original string
	for i := 1; i <= len(first); i++ {
		table[i][0] = float64(i)
	}
	// initialize the row representing the compare string
	for j := 1; j <= len(second); j++ {
		table[0][j] = float64(j)
	}
	// iterate through the table row by row
	for i := 1; i <= len(first); i++ {
		for j := 1; j <= len(second); j++ {
			if i > 1 && j > 1 && first[i-1] == second[j-2] && first[i-2] == second[j-1] {
				weight := getWeight(string(first[i-1]), string(second[j-1]))
				if first[i-1] == second[j-1] {
					table[i][j] = minimumFloat(table[i-1][j]+1, table[i][j-1]+1, table[i-1][j-1], table[i-2][j-2]+weight)
				} else {
					table[i][j] = minimumFloat(table[i-1][j]+1, table[i][j-1]+1, table[i-1][j-1]+weight, table[i-2][j-2]+weight)
				}
			} else {
				weight := getWeight(string(first[i-1]), string(second[j-1]))
				if first[i-1] == second[j-1] {
					table[i][j] = minimumFloat(table[i-1][j]+1, table[i][j-1]+1, table[i-1][j-1])
				} else {
					table[i][j] = minimumFloat(table[i-1][j]+1, table[i][j-1]+1, table[i-1][j-1]+weight)
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
	var result int
	if minimum(len(first), len(second)) == 0 {
		return maximum(len(first), len(second))
	}
	if len(first) > 1 && len(second) > 1 && first[0] == second[1] && first[1] == second[0] {
		if first[0] == second[0] {
			result += minimum(DamerauLevenstheinDistanceRecursive(first[1:], second)+1,
				DamerauLevenstheinDistanceRecursive(first, second[1:])+1,
				DamerauLevenstheinDistanceRecursive(first[1:], second[1:]),
				DamerauLevenstheinDistanceRecursive(first[2:], second[2:])+1)
		} else {
			result = minimum(DamerauLevenstheinDistanceRecursive(first[1:], second)+1,
				DamerauLevenstheinDistanceRecursive(first, second[1:])+1,
				DamerauLevenstheinDistanceRecursive(first[1:], second[1:])+1,
				DamerauLevenstheinDistanceRecursive(first[2:], second[2:])+1)
		}
	} else if first[0] == second[0] {
		result = minimum(DamerauLevenstheinDistanceRecursive(first[1:], second)+1,
			DamerauLevenstheinDistanceRecursive(first, second[1:])+1,
			DamerauLevenstheinDistanceRecursive(first[1:], second[1:]))
	} else {
		result = minimum(DamerauLevenstheinDistanceRecursive(first[1:], second),
			DamerauLevenstheinDistanceRecursive(first, second[1:]),
			DamerauLevenstheinDistanceRecursive(first[1:], second[1:])) + 1
	}
	// println("Comparing %s with %s and returning %d", first, second, result)
	return result
}

func DamerauLevenstheinDistanceRecursiveWeighted(first, second string) float64 {
	var result float64
	if minimum(len(first), len(second)) == 0 {
		return float64(maximum(len(first), len(second)))
	}
	if len(first) > 1 && len(second) > 1 && first[0] == second[1] && first[1] == second[0] {
		weight := getWeight(string(first[0]), string(second[0]))
		if first[0] == second[0] {
			result += minimumFloat(DamerauLevenstheinDistanceRecursiveWeighted(first[1:], second)+1,
				DamerauLevenstheinDistanceRecursiveWeighted(first, second[1:])+1,
				DamerauLevenstheinDistanceRecursiveWeighted(first[1:], second[1:]),
				DamerauLevenstheinDistanceRecursiveWeighted(first[2:], second[2:])+weight)
		} else {
			result += minimumFloat(DamerauLevenstheinDistanceRecursiveWeighted(first[1:], second)+1,
				DamerauLevenstheinDistanceRecursiveWeighted(first, second[1:])+1,
				DamerauLevenstheinDistanceRecursiveWeighted(first[1:], second[1:])+weight,
				DamerauLevenstheinDistanceRecursiveWeighted(first[2:], second[2:])+weight)
		}
	} else {
		weight := getWeight(string(first[0]), string(second[0]))
		if first[0] == second[0] {
			result = minimumFloat(DamerauLevenstheinDistanceRecursiveWeighted(first[1:], second)+1,
				DamerauLevenstheinDistanceRecursiveWeighted(first, second[1:])+1,
				DamerauLevenstheinDistanceRecursiveWeighted(first[1:], second[1:]))
		} else {
			result = minimumFloat(DamerauLevenstheinDistanceRecursiveWeighted(first[1:], second)+1,
				DamerauLevenstheinDistanceRecursiveWeighted(first, second[1:])+1,
				DamerauLevenstheinDistanceRecursiveWeighted(first[1:], second[1:])+weight)
		}
	}
	return result
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

func minimumFloat(numbers ...float64) float64 {
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

// LoadWeights loads config from file
func loadWeights(weights string) map[string]map[string]float64 {
	csvfile, err := os.Open(weights)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	weightMap := make(map[string]map[string]float64)
	defer csvfile.Close()
	csvreader := csv.NewReader(csvfile)
	csvreader.Comma = ','
	for {
		record, err := csvreader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		innerMap := make(map[string]float64)
		weight, _ := strconv.ParseFloat(string(record[2]), 64)
		innerMap[string(record[1])] = weight
		weightMap[string(record[0])] = innerMap
	}
	return weightMap
}

func getWeight(first, second string) float64 {
	if m, ok := weightMap[first]; ok {
		if val, ok := m[second]; ok {
			return val
		}
	}
	if m, ok := weightMap[second]; ok {
		if val, ok := m[first]; ok {
			return val
		}
	}
	return 1
}
