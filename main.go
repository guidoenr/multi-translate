package main

import (
	"fmt"
	"main/translator"
)

type input struct {
	text   string
	source string
	target string
}

func (in input) translate() (string, error) {
	return translator.Translate(in.text, in.source, in.target)
}

func main() {
	const text string = `Hello, World!`
	// you can use "auto" for source language
	// so, translator will detect language
	result, _ := translator.Translate(text, "en", "es")
	fmt.Println(result)
	// Output: "Hola, Mundo!"

}
