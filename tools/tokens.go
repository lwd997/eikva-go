package tools

import (
	"github.com/pkoukk/tiktoken-go"
)

func CountTokens(input string) int {
	enc, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		panic(err)
	}

	tokens := enc.Encode(input, nil, nil)
	return len(tokens)
}

