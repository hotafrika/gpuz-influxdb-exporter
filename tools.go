package main

import (
	"math"
	"strings"
)

func simplifyString(s string) string {
	s = strings.ToLower(s)
	var result strings.Builder
	for _, char := range s {
		if ('a' <= char && char <= 'z') || ('0' <= char && char <= '9') {
			result.WriteRune(char)
		} else {
			result.WriteRune('_')
		}
	}
	return result.String()
}

func roundFloat(f float64) float64 {
	return math.Round(f*100) / 100
}
