package main

import (
	"log"
	"strings"
	"testing"
	"time"
)

func containsSubstring(s, substr string) bool {
	return strings.Contains(s, substr)
}

func indexOfSubstring(s, substr string) int {
	return strings.Index(s, substr)
}

func TestContainsSubstring(t *testing.T) {
	start := time.Now()
	if !containsSubstring("hello world", "world") {
		t.Errorf("Expected true but got false")
	}
	if containsSubstring("hello world", "golang") {
		t.Errorf("Expected false but got true")
	}
	elapsed := time.Since(start)
	log.Printf("TestContainsSubstring took %s", elapsed)
}

func TestIndexOfSubstring(t *testing.T) {
	start := time.Now()
	if indexOfSubstring("hello world", "world") != 6 {
		t.Errorf("Expected 6 but got %d", indexOfSubstring("hello world", "world"))
	}
	if indexOfSubstring("hello world", "golang") != -1 {
		t.Errorf("Expected -1 but got %d", indexOfSubstring("hello world", "golang"))
	}
	elapsed := time.Since(start)
	log.Printf("TestIndexOfSubstring took %s", elapsed)
}
