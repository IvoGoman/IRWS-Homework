package main

import "testing"

func TestDamerauDynamicWithoutWeights1A(t *testing.T) {
	first := "amarillone"
	second := "smatillimr"
	resultFirst := DamerauLevenstheinDistanceDynamic(first, second)
	if resultFirst != 5 {
		t.Error("Expected 5, got ", resultFirst)
	}
}
func TestDamerauDynamicWithoutWeights1B(t *testing.T) {
	first := "maarillone"
	second := "amatilulimr"
	resultFirst := DamerauLevenstheinDistanceDynamic(first, second)
	if resultFirst != 6 {
		t.Error("Expected 6, got ", resultFirst)
	}
}
func TestDamerauDynamicWithoutWeights2A(t *testing.T) {
	first := "amarillone"
	second := "smtillimr"
	resultFirst := DamerauLevenstheinDistanceDynamicWeighted(first, second)
	if resultFirst != 2.75 {
		t.Error("Expected 2.75, got ", resultFirst)
	}
}
func TestDamerauDynamicWithoutWeights2B(t *testing.T) {
	first := "maarillone"
	second := "amatilulimr"
	resultFirst := DamerauLevenstheinDistanceDynamicWeighted(first, second)
	if resultFirst != 3.25 {
		t.Error("Expected 3.25, got ", resultFirst)
	}
}
func TestDamerauDynamicWithoutWeights2C(t *testing.T) {
	first := "saarillone"
	second := "asstilkimr"
	resultFirst := DamerauLevenstheinDistanceDynamicWeighted(first, second)
	if resultFirst != 3.25 {
		t.Error("Expected 3.25, got ", resultFirst)
	}
}
func TestDamerauRecursiveWithoutWeights1A(t *testing.T) {
	first := "amarillone"
	second := "smatillimr"
	resultFirst := DamerauLevenstheinDistanceRecursive(first, second)
	if resultFirst != 5 {
		t.Error("Expected 5, got ", resultFirst)
	}
}
func TestDamerauRecursiveWithoutWeights1B(t *testing.T) {
	first := "maarillone"
	second := "amatilulimr"
	resultFirst := DamerauLevenstheinDistanceRecursive(first, second)
	if resultFirst != 6 {
		t.Error("Expected 6, got ", resultFirst)
	}
}

func TestDamerauRecursiveWithWeights2A(t *testing.T) {
	first := "amarillone"
	second := "smtillimr"
	resultFirst := DamerauLevenstheinDistanceRecursiveWeighted(first, second)
	if resultFirst != 2.75 {
		t.Error("Expected 2.75, got ", resultFirst)
	}
}
func TestDamerauRecursiveWithWeights2B(t *testing.T) {
	first := "maarillone"
	second := "amatilulimr"
	resultFirst := DamerauLevenstheinDistanceRecursiveWeighted(first, second)
	if resultFirst != 3.25 {
		t.Error("Expected 3.25, got ", resultFirst)
	}
}
func TestDamerauRecursiveWithWeights2C(t *testing.T) {
	first := "saarillone"
	second := "asstilkimr"
	resultFirst := DamerauLevenstheinDistanceRecursiveWeighted(first, second)
	if resultFirst != 3.25 {
		t.Error("Expected 3.25, got ", resultFirst)
	}
}
