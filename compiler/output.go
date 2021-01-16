package compiler

import (
	"encoding/binary"
	"fmt"

	"github.com/takeru56/tcompiler/obj"
)

// output tarto IR bytecode Format
// ***************************************

// struct {
// 	u4 magic
// 	u2 constant_pool_count
// 	cp constant_pool[constant_pool_count]
// 	u2 instruction_count
// 	ins instructions[instruction_count]
// }

// struct constant_pool {
// 	u1 constant type
// 	u2 constant size
// 	c [const size]constants
// }
// ***************************************
type ConstantType byte

// Define Constant type
const (
	CONST_INT  ConstantType = iota
	CONST_FUNC ConstantType = iota
)

// TODO: 32bitに拡張+エラー処理

func (c *Compiler) Bytecode() string {
	b := ""
	// u4 magic（特に意味無し）
	b += fmt.Sprintf("%02x", []byte{255, 255, 255, 255})
	// u2 constant_pool_count
	b += fmt.Sprintf("%02x", toUint16(len(c.constantPool)))
	// const pool
	for _, constant := range c.constantPool {
		switch constant := constant.(type) {
		case *obj.Integer:
			// u1
			b += fmt.Sprintf("%02x", CONST_INT)
			// u2
			b += fmt.Sprintf("%02x", toUint16(constant.Size()))
			// u2
			b += fmt.Sprintf("%02x", toUint16(constant.Value))
		case *obj.Function:
			// u1
			b += fmt.Sprintf("%02x", CONST_FUNC)
			// u2
			b += fmt.Sprintf("%02x", toUint16(constant.Size()))
			// u2
			for _, bytecode := range constant.Instructions {
				b += fmt.Sprintf("%02x", bytecode)
			}
		}
	}
	// u2 instruction_count
	b += fmt.Sprintf("%02x", toUint16(len(c.scopes[c.scopeIndex].instructions)))

	// instruction
	for _, bytecode := range c.scopes[c.scopeIndex].instructions {
		b += fmt.Sprintf("%02x", bytecode)
	}
	return b
}

func (c *Compiler) Output() {
	fmt.Print(c.Bytecode())
}

func (c *Compiler) Dump() {
	b := c.Bytecode()
	p := 0
	size := 0
	for p < len(b) {
		if size%16 == 0 {
			if size != 0 {
				fmt.Print("\n")
			}
			fmt.Printf("%02X: ", size)
		}
		if size%16 != 0 && size%8 == 0 {
			fmt.Print(" ")
		}
		fmt.Print(b[p : p+2])
		p += 2
		size++
	}
}

func toUint16(num int) [2]byte {
	b := [2]byte{}
	binary.BigEndian.PutUint16(b[0:], uint16(num))
	return b
}
