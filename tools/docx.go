package tools

import (
	"archive/zip"
	"bytes"
	"golang.org/x/net/html"
	"io"
)

func GetDocumentXmlReader(templateBytes []byte) (io.Reader, error) {
	templateReader := bytes.NewReader(templateBytes)
	zipReader, err := zip.NewReader(templateReader, int64(len(templateBytes)))
	if err != nil {
		return nil, err
	}

	file, err := zipReader.Open("word/document.xml")
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetAllXmlText(reader io.Reader) string {
	var output string
	tokenizer := html.NewTokenizer(reader)
	prevToken := tokenizer.Token()
loop:
	for {
		tok := tokenizer.Next()
		switch {
		case tok == html.ErrorToken:
			break loop
		case tok == html.StartTagToken:
			prevToken = tokenizer.Token()
		case tok == html.TextToken:
			if prevToken.Data == "script" {
				continue
			}
			TxtContent := html.UnescapeString(string(tokenizer.Text()))
			if len(TxtContent) > 0 {
				output += TxtContent
			}
		}
	}
	return output
}
