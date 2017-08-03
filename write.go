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
	constant_pool []ConstantPoolEntry
	access_flags  uint16
}

type ConstantPoolEntry interface {
	IsConstantPoolEntry()
	Write(out io.Writer)
}

type ConstantUtf8 struct {
	string
}

func NewConstantUtf8(s string) ConstantUtf8 {
	return ConstantUtf8{
		string: s,
	}
}

func (self ConstantUtf8) Write(out io.Writer) {
	var type_ uint8 = CONSTANT_Utf8
	binary.Write(out, binary.BigEndian, type_)

	for _, c := range self.string {
		if c == 0 || c > 0x7f {
			panic(fmt.Errorf("Uh oh found bad byte %d in ConstantUtf8 %s",
				c, self.string))
		}
	}

	bytes := []byte(self.string)
	var length uint16 = uint16(len(bytes))
	binary.Write(out, binary.BigEndian, length)

	out.Write(bytes)
}

func (self ConstantUtf8) IsConstantPoolEntry() {}

const CONSTANT_Class = 7
const CONSTANT_Fieldref = 9
const CONSTANT_Methodref = 10
const CONSTANT_InterfaceMethodref = 11
const CONSTANT_String = 8
const CONSTANT_Integer = 3
const CONSTANT_Float = 4
const CONSTANT_Long = 5
const CONSTANT_Double = 6
const CONSTANT_NameAndType = 12
const CONSTANT_Utf8 = 1
const CONSTANT_MethodHandle = 15
const CONSTANT_MethodType = 16
const CONSTANT_InvokeDynamic = 18

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

	var constant_pool_count uint16 = uint16(len(self.constant_pool))
	binary.Write(out, binary.BigEndian, constant_pool_count)
	for _, entry := range self.constant_pool {
		entry.Write(out)
	}

	binary.Write(out, binary.BigEndian, self.access_flags)
}

func main() {
	classFile := ClassFile{
		magic:         0xCAFEBABE,
		minor_version: 0,
		major_version: 52, // 1.8
		constant_pool: []ConstantPoolEntry{
			NewConstantUtf8("System"),
		},
		access_flags: ACC_PUBLIC,
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
