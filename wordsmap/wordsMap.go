package wordsmap

import (
    "embed"
    "strings"
)

// Эта строка говорит Go: "возьми файл words.txt и положи его в переменную wordsFile"

//go:embed words.txt
var wordsFile embed.FS

func LoadEmbeddedDictionary() (map[string]bool, []string) {
	dictMap := make(map[string]bool)
	dictSlice := []string{}

	data, _ := wordsFile.ReadFile("words.txt")
	lines := strings.Split(string(data), "\n")


	for _, w := range lines {
		word := strings.ToUpper(strings.TrimSpace(w))

		if len(word) == 5 { 
				dictMap[word] = true
				dictSlice = append(dictSlice, word)
			}
	}
    return dictMap, dictSlice
}