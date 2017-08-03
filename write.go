package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
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
}

const ACC_PUBLIC = 0x0001    // 	Declared public; may be accessed from outside its package.
const ACC_PRIVATE = 0x0002   // 	Declared private; usable only within the defining class.
const ACC_PROTECTED = 0x0004 // 	Declared protected; may be accessed within subclasses.
const ACC_STATIC = 0x0008    // 	Declared static.
const ACC_FINAL = 0x0010     // 	Declared final; never directly assigned to after object construction (JLS §17.5).
const ACC_VOLATILE = 0x0040  // 	Declared volatile; cannot be cached.
const ACC_TRANSIENT = 0x0080 // 	Declared transient; not written or read by a persistent object manager.
const ACC_SYNTHETIC = 0x1000 // 	Declared synthetic; not present in the source code.
const ACC_ENUM = 0x4000      // 	Declared as an element of an enum.

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

	var numMethods uint16 = 0
	binary.Write(out, binary.BigEndian, numMethods)

	var numAttributes uint16 = 0
	binary.Write(out, binary.BigEndian, numAttributes)

}

func main() {
	constantPool := NewConstantPool()
	constantPool.Add(NewConstantUtf8("System"))
	MinimalString := constantPool.Add(NewConstantUtf8("Minimal"))
	MinimalClass := constantPool.Add(NewConstantClass(MinimalString))
	javaLangObjectString := constantPool.Add(NewConstantUtf8("java/lang/Object"))
	javaLangObjectClass :=
		constantPool.Add(NewConstantClass(javaLangObjectString))

	classFile := ClassFile{
		magic:         0xCAFEBABE,
		minor_version: 0,
		major_version: 52, // 1.8
		constantPool:  constantPool,
		access_flags:  ACC_PUBLIC,
		this_class:    MinimalClass,
		super_class:   javaLangObjectClass,
	}

	path := "Out.class"
	out, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	classFile.Write(out)
	fmt.Println(path)
}
