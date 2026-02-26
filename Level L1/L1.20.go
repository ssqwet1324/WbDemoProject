package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReverseWords - способ через стринг билдер
func ReverseWords(str string) string {
	var b strings.Builder
	b.Grow(len(str))

	i := len(str) - 1

	for i >= 0 {
		for i >= 0 && str[i] == ' ' {
			i--
		}
		if i < 0 {
			break
		}

		end := i

		// ищем начало слова
		for i >= 0 && str[i] != ' ' {
			i--
		}
		start := i + 1

		// добавляем слово
		b.WriteString(str[start : end+1])

		if i >= 0 {
			b.WriteByte(' ')
		}
	}

	return b.String()
}

func ReverseWordsRune(str string) {
	runes := []rune(str)

	// разворачиваем всю строку
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	//разворачиваем каждое слово отдельно
	start := 0
	for i := 0; i <= len(runes); i++ {
		if i == len(runes) || runes[i] == ' ' {
			for l, r := start, i-1; l < r; l, r = l+1, r-1 {
				runes[l], runes[r] = runes[r], runes[l]
			}
			start = i + 1
		}
	}

	fmt.Println(string(runes))
}

func main() {
	fmt.Println("Введите строку:")
	in, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	in = strings.TrimSpace(in)

	// 1й способ
	fmt.Println(ReverseWords(in))
	// 2й способ
	ReverseWordsRune(in)
}
