package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	var tokens []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		tokens = append(tokens, scanner.Text())
	}
	generate(tokens)
	generate(tokens)
	generate(tokens)
	generate(tokens)
	generate(tokens)
	generate(tokens)
	generate(tokens)
	generate(tokens)
	generate(tokens)
	generate(tokens)
	generate(tokens)
}

func out(size int, tokens []string) {
	i := -1
	for {
		i++
		n := i + size
		if n == len(tokens) {
			break
		}
		line := strings.Join(tokens[i:n], ",")
		line = fmt.Sprintf("%s|%s", mrn(7), line)
		fmt.Println(line)
	}
}

func generate(tokens []string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(tokens), func(i, j int) { tokens[i], tokens[j] = tokens[j], tokens[j]})
	out(1, tokens)
	out(2, tokens)
	out(3, tokens)
	out(4, tokens)
	out(5, tokens)
}

func mrn(size int) string {
	runes := make([]rune, size)
	alphabet := []rune("1234567890")
	n := len(alphabet)
	rand.Seed(time.Now().UnixNano())
	for i := range runes {
		runes[i] = alphabet[rand.Intn(n)]
	}
	return fmt.Sprintf("M%s", string(runes))
}