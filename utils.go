package main

import "crypto/sha256"

func ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

func hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}
