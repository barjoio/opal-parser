package main

import (
	"fmt"
	"testing"
)

func main() {

	char := '#'

	r1 := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'w', 'x', 'y', 'z', '#'} {
				if char == v {
					break
				}
			}
		}
	})

	r2 := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if char == 'a' || char == 'b' || char == 'c' || char == 'd' || char == 'e' || char == 'f' || char == 'g' || char == 'h' || char == 'i' || char == 'j' || char == 'k' || char == 'l' || char == 'm' || char == 'n' || char == 'o' || char == 'p' || char == 'q' || char == 'r' || char == 's' || char == 't' || char == 'u' || char == 'w' || char == 'x' || char == 'y' || char == 'z' || char == '#' {
			}
		}
	})

	r3 := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			switch char {
			case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'w', 'x', 'y', 'z', '#':
			}
		}
	})

	fmt.Println(r1)
	fmt.Println(r2)
	fmt.Println(r3)
}

// 100000000               10.1 ns/op
// 393648298                3.08 ns/op
// 1000000000               0.218 ns/op
