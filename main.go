package main

import (
	"bufio"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("sample.txt")
	if err != nil {
		os.Exit(1)
	}

	output, err := os.Create("result.txt")
	if err != nil {
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		modified := processLine(line)
		fmt.Fprintln(output, modified)
	}
}

func processLine(line string) string {
	line = brackets(line)
	line = article(line)
	line = punctuation(line)
	return line
}

func brackets(line string) string {
	caser := cases.Title(language.Und, cases.NoLower)
	words := strings.Fields(line)
	for index, word := range words {
		switch word {
		case "(cap)":
			words[index-1] = caser.String(words[index-1])
			remove(&words, index, 1)
		case "(up)":
			words[index-1] = strings.ToUpper(words[index-1])
			remove(&words, index, 1)
		case "(low)":
			words[index-1] = strings.ToLower(words[index-1])
			remove(&words, index, 1)
		case "(hex)":
			decimal, _ := strconv.ParseInt(words[index-1], 16, 64)
			words[index-1] = strconv.Itoa(int(decimal))
			remove(&words, index, 1)
		case "(bin)":
			decimal, _ := strconv.ParseInt(words[index-1], 2, 64)
			words[index-1] = strconv.Itoa(int(decimal))
			remove(&words, index, 1)
		case "(cap,":
			convertString(words, index, caser.String)
			remove(&words, index, 2)
		case "(low,":
			convertString(words, index, strings.ToLower)
			remove(&words, index, 2)
		case "(up,":
			convertString(words, index, strings.ToUpper)
			remove(&words, index, 2)
		}
	}
	return strings.Join(words, " ")
}

func article(line string) string {
	vowels := "aeiouh"
	words := strings.Fields(line)
	for index, word := range words {
		if len(word) == 1 && string(word[0]) == "a" || string(word[0]) == "A" {
			for _, vowel := range vowels {
				if string(words[index+1][0]) == string(vowel) {
					words[index] = words[index] + "n"
					break
				}
			}
		}
	}
	return strings.Join(words, " ")
}

func punctuation(line string) string {
	punctuationRegex := regexp.MustCompile(`\s*([.,!?:;]+)\s*`)
	modified := punctuationRegex.ReplaceAllString(line, "$1 ")
	modified = strings.TrimRight(modified, " ")
	return modified
}

func apostrophes(line string) string {
	// TODO: Implement
	return ""
}

func convertString(words []string, index int, fn func(string) string) {
	strNum := words[index+1][:len(words[index+1])-1]
	num, _ := strconv.Atoi(strNum)
	for i := index - 1; i >= index-num; i-- {
		words[i] = fn(words[i])
	}
}

func remove(slice *[]string, index, count int) {
	*slice = append((*slice)[:index], (*slice)[index+count:]...)
}
