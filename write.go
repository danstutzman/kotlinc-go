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

func main() {
	constantPool := NewConstantPool()
	constantPool.Add(NewConstantUtf8("System"))
	MinimalGoString := constantPool.Add(NewConstantUtf8("MinimalGo"))
	MinimalGoClass := constantPool.Add(NewConstantClass(MinimalGoString))
	javaLangObjectString := constantPool.Add(NewConstantUtf8("java/lang/Object"))
	javaLangObjectClass :=
		constantPool.Add(NewConstantClass(javaLangObjectString))
	initString := constantPool.Add(NewConstantUtf8("<init>"))
	noArgsString := constantPool.Add(NewConstantUtf8("()V"))
	CodeString := constantPool.Add(NewConstantUtf8("Code"))
	//LineNumberTableString := constantPool.Add(NewConstantUtf8("LineNumberTable"))
	//SourceFileString := constantPool.Add(NewConstantUtf8("SourceFile"))
	//MinimalGoDotJavaString := constantPool.Add(NewConstantUtf8("Minimal.java"))
	initNoArgsNameAndType :=
		constantPool.Add(NewConstantNameAndType(initString, noArgsString))
	javaLangObjectInit :=
		constantPool.Add(NewConstantMethodRef(
			javaLangObjectClass, initNoArgsNameAndType))

	initCode := CodeAttribute{
		attribute_name_index: CodeString,
		max_stack:            1,
		max_locals:           1,
		instructionsSerialized: []byte{
			ALOAD_0,
			INVOKE_SPECIAL, 0, uint8(uint16(javaLangObjectInit) % 256),
			RETURN,
		},
		attributes: []Attribute{},
	}

	initMethod := Method{
		access_flags:     ACC_PUBLIC,
		name_index:       initString,
		descriptor_index: noArgsString,
		attributes:       []Attribute{initCode},
	}

	classFile := ClassFile{
		magic:         0xCAFEBABE,
		minor_version: 0,
		major_version: 52, // 1.8
		constantPool:  constantPool,
		access_flags:  ACC_PUBLIC,
		this_class:    MinimalGoClass,
		super_class:   javaLangObjectClass,
		methods:       []Method{initMethod},
	}

	path := "MinimalGo.class"
	out, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	classFile.Write(out)
	fmt.Println(path)
}
