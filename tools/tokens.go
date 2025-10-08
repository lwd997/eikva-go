package tools

import (
	"github.com/pkoukk/tiktoken-go"
	"github.com/pkoukk/tiktoken-go-loader"
)

func CountTokens(input string) int {
	tiktoken.SetBpeLoader(tiktoken_loader.NewOfflineLoader())
	enc, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		panic(err)
	}

	tokens := enc.Encode(input, nil, nil)
	return len(tokens)
}

