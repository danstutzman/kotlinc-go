package assembler

import (
	"encoding/binary"
	"io"
)

// See https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html
type ClassFile struct {
	magic         uint32
	minor_version uint16
	major_version uint16
	constantPool  ConstantPool
	access_flags  uint16
	this_class    PoolIndex
	super_class   PoolIndex
	methods       []Method
}

func (self *ClassFile) Write(out io.Writer) {
	binary.Write(out, binary.BigEndian, self.magic)
	binary.Write(out, binary.BigEndian, self.minor_version)
	binary.Write(out, binary.BigEndian, self.major_version)

	self.constantPool.Write(out)

	binary.Write(out, binary.BigEndian, self.access_flags)
	binary.Write(out, binary.BigEndian, self.this_class)
	binary.Write(out, binary.BigEndian, self.super_class)

	var numInterfaces uint16 = 0
	binary.Write(out, binary.BigEndian, numInterfaces)

	var numFields uint16 = 0
	binary.Write(out, binary.BigEndian, numFields)

	numMethods := uint16(len(self.methods))
	binary.Write(out, binary.BigEndian, numMethods)
	for _, method := range self.methods {
		method.Write(out)
	}

	var numAttributes uint16 = 0
	binary.Write(out, binary.BigEndian, numAttributes)
}
