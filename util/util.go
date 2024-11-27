package util

import "github.com/gdanko/clearsky/globals"

func SliceChunker(input []globals.BlockingUser, chunkSize int) (output [][]globals.BlockingUser) {
	for i := 0; i < len(input); i += chunkSize {
		end := i + chunkSize
		if end > len(input) {
			end = len(input)
		}
		output = append(output, input[i:end])
	}
	return output
}
