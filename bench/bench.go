package main

import (
	"fmt"
	"testing"
)

func main() {
	items := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'w', 'x', 'y', 'z'}

	char := '#'

	r1 := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range items {
				if char == v {
					break
				}
			}
		}
	})

	r2 := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			switch char {
			case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'w', 'x', 'y', 'z':
			}
		}
	})

	fmt.Println(r1)
	fmt.Println(r2)
}

// 182712844                6.50 ns/op
// 1000000000               0.429 ns/op
