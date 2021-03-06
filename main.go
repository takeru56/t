package main

import (
	"fmt"
	"log"
	"os"

	"github.com/takeru56/tcompiler/compiler"
	"github.com/takeru56/tcompiler/parser"
	"github.com/takeru56/tcompiler/token"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing argument error")
	}
	// from file
	// bytes, err := ioutil.ReadFile("in.txt")
	// if err != nil {
	// 	panic(err)
	// }
	source := os.Args[1]
	tok := token.New(source)
	parser, err := parser.New(tok)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	p, err := parser.Program()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c := compiler.Exec(p)
	c.Output()
	// c.Dump()
}
