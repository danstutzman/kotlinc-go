package assembler

import (
	"fmt"
	"os"
)

func Demo(outPath string) {
	constantPool := NewConstantPool()
	constantPool.Add(NewConstantUtf8("System"))
	MinimalGoUtf8 := constantPool.Add(NewConstantUtf8("MinimalGo"))
	MinimalGoClass := constantPool.Add(NewConstantClass(MinimalGoUtf8))
	javaLangObjectUtf8 := constantPool.Add(NewConstantUtf8("java/lang/Object"))
	javaLangObjectClass :=
		constantPool.Add(NewConstantClass(javaLangObjectUtf8))
	initUtf8 := constantPool.Add(NewConstantUtf8("<init>"))
	noArgsUtf8 := constantPool.Add(NewConstantUtf8("()V"))
	CodeUtf8 := constantPool.Add(NewConstantUtf8("Code"))
	//LineNumberTableUtf8 := constantPool.Add(NewConstantUtf8("LineNumberTable"))
	//SourceFileUtf8 := constantPool.Add(NewConstantUtf8("SourceFile"))
	//MinimalGoDotJavaUtf8 := constantPool.Add(NewConstantUtf8("Minimal.java"))
	initNoArgsNameAndType :=
		constantPool.Add(NewConstantNameAndType(initUtf8, noArgsUtf8))
	javaLangObjectInit :=
		constantPool.Add(NewConstantMethodRef(
			javaLangObjectClass, initNoArgsNameAndType))
	mainUtf8 := constantPool.Add(NewConstantUtf8("main"))
	javaLangSystemUtf8 := constantPool.Add(NewConstantUtf8("java/lang/System"))
	outUtf8 := constantPool.Add(NewConstantUtf8("out"))
	returnsPrintStreamUtf8 := constantPool.Add(NewConstantUtf8("Ljava/io/PrintStream;"))
	javaLangSystemClass :=
		constantPool.Add(NewConstantClass(javaLangSystemUtf8))
	outReturnsPrintStreamNameAndType := constantPool.Add(
		NewConstantNameAndType(outUtf8, returnsPrintStreamUtf8))
	javaLangSystemOutFieldRef := constantPool.Add(NewConstantFieldRef(
		javaLangSystemClass, outReturnsPrintStreamNameAndType,
	))
	helloUtf8 := constantPool.Add(NewConstantUtf8("hello"))
	helloString := constantPool.Add(NewConstantString(helloUtf8))
	stringArgNoReturnUtf8 :=
		constantPool.Add(NewConstantUtf8("(Ljava/lang/String;)V"))
	stringArrayArgNoReturnUtf8 :=
		constantPool.Add(NewConstantUtf8("([Ljava/lang/String;)V"))

	javaIoPrintStreamUtf8 :=
		constantPool.Add(NewConstantUtf8("java/io/PrintStream"))
	javaIoPrintStreamClass :=
		constantPool.Add(NewConstantClass(javaIoPrintStreamUtf8))
	printlnUtf8 := constantPool.Add(NewConstantUtf8("println"))
	printlnTakesStringArg := constantPool.Add(
		NewConstantNameAndType(printlnUtf8, stringArgNoReturnUtf8),
	)
	javaIoPrintStreamPrintlnMethodRef := constantPool.Add(NewConstantMethodRef(
		javaIoPrintStreamClass,
		printlnTakesStringArg,
	))

	initCode := CodeAttribute{
		attribute_name_index: CodeUtf8,
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
		name_index:       initUtf8,
		descriptor_index: noArgsUtf8,
		attributes:       []Attribute{initCode},
	}

	mainCode := CodeAttribute{
		attribute_name_index: CodeUtf8,
		max_stack:            2, // TODO
		max_locals:           1, // TODO
		instructionsSerialized: []byte{
			GETSTATIC, 0, uint8(uint16(javaLangSystemOutFieldRef) % 256),
			LDC, uint8(uint16(helloString) % 256),
			INVOKEVIRTUAL, 0, uint8(uint16(javaIoPrintStreamPrintlnMethodRef) % 256),
			RETURN,
		},
		attributes: []Attribute{},
	}
	/*
	   public static void main();
	     descriptor: ()V
	     Code:
	        0: getstatic     #2                  // Field java/lang/System.out:Ljava/io/PrintStream;
	        3: ldc           #3                  // String Hello
	        5: invokevirtual #4                  // Method java/io/PrintStream.println:(Ljava/lang/String;)V
	        8: return
	*/

	mainMethod := Method{
		access_flags:     ACC_PUBLIC | ACC_STATIC,
		name_index:       mainUtf8,
		descriptor_index: stringArrayArgNoReturnUtf8,
		attributes:       []Attribute{mainCode},
	}

	classFile := ClassFile{
		magic:         0xCAFEBABE,
		minor_version: 0,
		major_version: 52, // 1.8
		constantPool:  constantPool,
		access_flags:  ACC_PUBLIC,
		this_class:    MinimalGoClass,
		super_class:   javaLangObjectClass,
		methods:       []Method{initMethod, mainMethod},
	}

	out, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	classFile.Write(out)
	fmt.Println(outPath)
}
