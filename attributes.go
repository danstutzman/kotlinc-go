package main

import (
	"encoding/binary"
	"io"
)

type Attribute interface {
	Write(out io.Writer)
}

//////////

type CodeAttribute struct {
	attribute_name_index   PoolIndex
	max_stack              uint16
	max_locals             uint16
	instructionsSerialized []byte
	attributes             []Attribute
}

const ALOAD_0 = 42
const INVOKE_SPECIAL = 183
const RETURN = 177

func (self CodeAttribute) Write(out io.Writer) {
	binary.Write(out, binary.BigEndian, self.attribute_name_index)

	var thisLength uint32 = 2 + 2 + 4 + 5 + 2 + 2 // TODO
	binary.Write(out, binary.BigEndian, thisLength)

	binary.Write(out, binary.BigEndian, self.max_stack)
	binary.Write(out, binary.BigEndian, self.max_locals)

	codeLength := uint32(len(self.instructionsSerialized))
	binary.Write(out, binary.BigEndian, codeLength)
	out.Write(self.instructionsSerialized)

	var exceptionTableLength uint16 = 0
	binary.Write(out, binary.BigEndian, exceptionTableLength)

	numAttributes := uint16(len(self.attributes))
	binary.Write(out, binary.BigEndian, numAttributes)
	for _, attribute := range self.attributes {
		attribute.Write(out)
	}
}
