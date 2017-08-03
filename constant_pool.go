package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ConstantPool struct {
	poolEntries          []ConstantPoolEntry
	poolEntryToPoolIndex map[ConstantPoolEntry]PoolIndex
}

type PoolIndex uint16

func NewConstantPool() ConstantPool {
	return ConstantPool{
		poolEntries:          []ConstantPoolEntry{nil},
		poolEntryToPoolIndex: map[ConstantPoolEntry]PoolIndex{},
	}
}

func (self *ConstantPool) Add(entry ConstantPoolEntry) PoolIndex {
	self.poolEntries = append(self.poolEntries, entry)
	poolIndex := PoolIndex(uint16(len(self.poolEntries) - 1))
	self.poolEntryToPoolIndex[entry] = poolIndex
	return poolIndex
}

func (self *ConstantPool) Write(out io.Writer) {
	var constant_pool_count uint16 = uint16(len(self.poolEntries))
	binary.Write(out, binary.BigEndian, constant_pool_count)

	for i, entry := range self.poolEntries {
		if i > 0 {
			entry.Write(out)
		}
	}
}

///////////

type ConstantPoolEntry interface {
	IsConstantPoolEntry()
	Write(out io.Writer)
}

///////////////

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

///////////

type ConstantClass struct {
	nameIndex PoolIndex
}

func NewConstantClass(nameIndex PoolIndex) ConstantClass {
	return ConstantClass{
		nameIndex: nameIndex,
	}
}

func (self ConstantClass) Write(out io.Writer) {
	var type_ uint8 = CONSTANT_Class
	binary.Write(out, binary.BigEndian, type_)

	binary.Write(out, binary.BigEndian, self.nameIndex)
}

func (self ConstantClass) IsConstantPoolEntry() {}

//////

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
