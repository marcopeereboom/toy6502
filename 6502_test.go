package main

import (
	"os"
	"testing"
)

func TestPha(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0xff
	c.memory[c.pc] = 0x48 // pha
	c.executeInstruction()
	if c.memory[0x01ff] != 0xff {
		t.Fatalf("unexpected memory %0x", c.memory[0x01ff])
	}
	if c.sp != 0xfe {
		t.Fatalf("unexpected sp %0x", c.sp)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test again
	c.a = 0xf0
	c.memory[0x1001] = 0x48 // pha
	c.executeInstruction()
	if c.memory[0x01fe] != 0xf0 {
		t.Fatalf("unexpected memory %0x", c.memory[0x01fe])
	}
}

func TestPhp(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.sr = 0xff
	c.memory[c.pc] = 0x08 // php
	c.executeInstruction()
	if c.memory[0x01ff] != 0xff {
		t.Fatalf("unexpected memory %0x", c.memory[0x01ff])
	}
	if c.sp != 0xfe {
		t.Fatalf("unexpected sp %0x", c.sp)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test again
	c.sr = 0xf0
	c.memory[0x1001] = 0x08 // php
	c.executeInstruction()
	if c.memory[0x01fe] != 0xf0 {
		t.Fatalf("unexpected memory %0x", c.memory[0x01fe])
	}
}

func TestPlp(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.sp = 0xfe
	c.sr = 0xff
	c.memory[c.pc] = 0x28 // plp
	c.memory[0x01ff] = 0x55
	c.executeInstruction()
	if c.sr != 0x75 {
		t.Fatalf("unexpected sr %0x", c.sr)
	}
	if c.sp != 0xff {
		t.Fatalf("unexpected sp %0x", c.sp)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test again with overflow
	c.memory[0x1001] = 0x28 // plp
	c.memory[0x0100] = 0xaa
	c.executeInstruction()
	if c.sr != 0xaa {
		t.Fatalf("unexpected sr %0x", c.sr)
	}
	if c.sp != 0x00 {
		t.Fatalf("unexpected sp %0x", c.sp)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestPla(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.sp = 0xfe
	c.a = 0xff
	c.memory[c.pc] = 0x68 // pla
	c.memory[0x01ff] = 0x55
	c.executeInstruction()
	if c.a != 0x55 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.sp != 0xff {
		t.Fatalf("unexpected sp %0x", c.sp)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test again with overflow
	c.memory[0x1001] = 0x68 // pla
	c.memory[0x0100] = 0xaa
	c.executeInstruction()
	if c.a != 0xaa {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.sp != 0x00 {
		t.Fatalf("unexpected sp %0x", c.sp)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestOraImmediate(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.memory[c.pc] = 0x09
	c.memory[c.pc+1] = 0xaa
	c.executeInstruction()
	if c.a != 0xff {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test zero flag
	c.pc = 0x1000
	c.a = 0x00
	c.memory[c.pc] = 0x09
	c.memory[c.pc+1] = 0x00
	c.executeInstruction()
	if c.a != 0x00 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestOraZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.memory[c.pc] = 0x05
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0xaa
	c.executeInstruction()
	if c.a != 0xff {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestOraZPX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.x = 0x08
	c.memory[c.pc] = 0x15
	c.memory[c.pc+1] = 0x80
	c.memory[0x88] = 0xaa
	c.executeInstruction()
	if c.a != 0xff {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestOraAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.memory[c.pc] = 0x0d
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0x80
	c.memory[0x8000] = 0xaa
	c.executeInstruction()
	if c.a != 0xff {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestOraAbsoluteX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.x = 0x08
	c.memory[c.pc] = 0x1d
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0x80
	c.memory[0x8008] = 0xaa
	c.executeInstruction()
	if c.a != 0xff {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestOraAbsoluteY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.y = 0xf0
	c.memory[c.pc] = 0x19
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0x80
	c.memory[0x80f0] = 0xaa
	c.executeInstruction()
	if c.a != 0xff {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestOraIndirectX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.x = 0x10

	c.memory[c.pc] = 0x01 // ora(0x20,x)
	c.memory[c.pc+1] = 0x20

	c.memory[0x30] = 0x04 // low byte
	c.memory[0x31] = 0xd0 // high byte
	c.memory[0xd004] = 0xaa
	c.executeInstruction()
	if c.a != 0xff {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestOraIndirectY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.y = 0x10

	c.memory[c.pc] = 0x11 // ora(0x20),y
	c.memory[c.pc+1] = 0x20

	c.memory[0x20] = 0x04 // low byte
	c.memory[0x21] = 0xd0 // high byte
	c.memory[0xd014] = 0xaa
	c.executeInstruction()
	if c.a != 0xff {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAsl(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x4e
	c.memory[c.pc] = 0x0a
	c.executeInstruction()
	if c.a != 0x9c {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAslZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x06
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0xee
	c.executeInstruction()
	if c.memory[0x80] != 0xdc {
		t.Fatalf("unexpected memory %0x", c.memory[0x80])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAslZPX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x01
	c.memory[c.pc] = 0x16
	c.memory[c.pc+1] = 0x80
	c.memory[0x81] = 0xee
	c.executeInstruction()
	if c.memory[0x81] != 0xdc {
		t.Fatalf("unexpected memory %0x", c.memory[0x81])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAslAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x0e
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd000] = 0xee
	c.executeInstruction()
	if c.memory[0xd000] != 0xdc {
		t.Fatalf("unexpected memory %0x", c.memory[0xd000])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAslAbsoluteX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x20
	c.memory[c.pc] = 0x1e
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd020] = 0xee
	c.executeInstruction()
	if c.memory[0xd020] != 0xdc {
		t.Fatalf("unexpected memory %0x", c.memory[0xd020])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestBpl(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x10   // bpl
	c.memory[c.pc+1] = 0x06 // pc+8
	c.executeInstruction()
	if c.pc != 0x1008 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// go backwards
	c.pc = 0x1000
	c.memory[c.pc] = 0x10   // bpl
	c.memory[c.pc+1] = 0xfa // pc-4
	c.executeInstruction()
	if c.pc != 0x1000-0x04 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// don't branch
	c.pc = 0x1000
	c.sr |= Negative
	c.memory[c.pc] = 0x10   // bpl
	c.memory[c.pc+1] = 0x06 // pc+8
	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestBmi(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.sr |= Negative
	c.memory[c.pc] = 0x30   // bmi
	c.memory[c.pc+1] = 0x06 // pc+8
	c.executeInstruction()
	if c.pc != 0x1008 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// go backwards
	c.pc = 0x1000
	c.memory[c.pc] = 0x30   // bmi
	c.memory[c.pc+1] = 0xfa // pc-4
	c.executeInstruction()
	if c.pc != 0x1000-0x04 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// don't branch
	c.pc = 0x1000
	c.sr = 0
	c.memory[c.pc] = 0x30   // bmi
	c.memory[c.pc+1] = 0x06 // pc+8
	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestBvc(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x50   // bvc
	c.memory[c.pc+1] = 0x06 // pc+8
	c.executeInstruction()
	if c.pc != 0x1008 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// go backwards
	c.pc = 0x1000
	c.memory[c.pc] = 0x50   // bvc
	c.memory[c.pc+1] = 0xfa // pc-4
	c.executeInstruction()
	if c.pc != 0x1000-0x04 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// don't branch
	c.pc = 0x1000
	c.sr |= Overflow
	c.memory[c.pc] = 0x50   // bvc
	c.memory[c.pc+1] = 0x06 // pc+8
	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestBcc(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x90   // bcc
	c.memory[c.pc+1] = 0x06 // pc+8
	c.executeInstruction()
	if c.pc != 0x1008 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// go backwards
	c.pc = 0x1000
	c.memory[c.pc] = 0x90   // bcc
	c.memory[c.pc+1] = 0xfa // pc-4
	c.executeInstruction()
	if c.pc != 0x1000-0x04 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// don't branch
	c.pc = 0x1000
	c.sr |= Carry
	c.memory[c.pc] = 0x90   // bvc
	c.memory[c.pc+1] = 0x06 // pc+8
	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestBvs(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x70   // bvs
	c.memory[c.pc+1] = 0x06 // pc+8
	c.sr |= Overflow
	c.executeInstruction()
	if c.pc != 0x1008 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// go backwards
	c.pc = 0x1000
	c.memory[c.pc] = 0x70   // bvs
	c.memory[c.pc+1] = 0xfa // pc-4
	c.sr |= Overflow
	c.executeInstruction()
	if c.pc != 0x1000-0x04 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// don't branch
	c.pc = 0x1000
	c.sr |= Overflow
	c.memory[c.pc] = 0x70   // bvs
	c.memory[c.pc+1] = 0x06 // pc+8
	c.sr = 0
	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestBcs(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xb0   // bcs
	c.memory[c.pc+1] = 0x06 // pc+8
	c.sr |= Carry
	c.executeInstruction()
	if c.pc != 0x1008 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// go backwards
	c.pc = 0x1000
	c.memory[c.pc] = 0xb0   // bcs
	c.memory[c.pc+1] = 0xfa // pc-4
	c.sr |= Carry
	c.executeInstruction()
	if c.pc != 0x1000-0x04 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// don't branch
	c.pc = 0x1000
	c.sr = 0
	c.memory[c.pc] = 0xb0   // bcs
	c.memory[c.pc+1] = 0x06 // pc+8
	c.sr = 0
	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestBeq2(t *testing.T) {
	c := New()

	c.pc = 0x0300
	c.memory[c.pc] = 0xf0 // beq
	c.memory[c.pc+1] = 0x05
	c.sr |= Zero
	c.executeInstruction()
	if c.pc != 0x0307 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ := c.disassemble(0x0300)
	t.Logf("%v\n", d)

	c.pc = 0x0300
	c.memory[c.pc] = 0xf0 // beq
	c.memory[c.pc+1] = 0x05
	c.sr = 0
	c.executeInstruction()
	if c.pc != 0x0302 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x0300)
	t.Logf("%v\n", d)

	c.pc = 0x0300
	c.memory[c.pc] = 0xf0 // beq
	c.memory[c.pc+1] = 0xfb
	c.sr |= Zero
	c.executeInstruction()
	if c.pc != 0x02fd {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x0300)
	t.Logf("%v\n", d)

	c.pc = 0x0300
	c.memory[c.pc] = 0xf0 // beq
	c.memory[c.pc+1] = 0xfb
	c.sr = 0
	c.executeInstruction()
	if c.pc != 0x0302 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x0300)
	t.Logf("%v\n", d)
}

func TestBeq(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xf0   // beq
	c.memory[c.pc+1] = 0x06 // pc+8
	c.sr |= Zero
	c.executeInstruction()
	if c.pc != 0x1008 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// go backwards
	c.pc = 0x1000
	c.memory[c.pc] = 0xf0   // beq
	c.memory[c.pc+1] = 0xfa // pc-4
	c.sr |= Zero
	c.executeInstruction()
	if c.pc != 0x1000-0x04 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// don't branch
	c.pc = 0x1000
	c.sr = 0
	c.memory[c.pc] = 0xf0   // beq
	c.memory[c.pc+1] = 0x06 // pc+8
	c.sr = 0
	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestBne(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xd0   // bne
	c.memory[c.pc+1] = 0x06 // pc+8
	c.executeInstruction()
	if c.pc != 0x1008 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// go backwards
	c.pc = 0x1000
	c.memory[c.pc] = 0xd0   // bne
	c.memory[c.pc+1] = 0xfa // pc-4
	c.sr |= Carry
	c.executeInstruction()
	if c.pc != 0x1000-0x04 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// don't branch
	c.pc = 0x1000
	c.sr = 0
	c.memory[c.pc] = 0xd0   // bne
	c.memory[c.pc+1] = 0x06 // pc+8
	c.sr |= Zero
	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestBne2(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xd0 // bne
	c.memory[c.pc+1] = 0x05
	c.sr |= Zero
	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// no zero
	c = New()
	c.pc = 0x1000
	c.memory[c.pc] = 0xd0 // bne
	c.memory[c.pc+1] = 0x05
	c.executeInstruction()
	if c.pc != 0x1007 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// backwards zero set
	c = New()
	c.pc = 0x1000
	c.memory[c.pc] = 0xd0 // bne
	c.memory[c.pc+1] = 0xfb
	c.sr |= Zero
	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// backwards zero not set
	c = New()
	c.pc = 0x1000
	c.memory[c.pc] = 0xd0 // bne
	c.memory[c.pc+1] = 0xfb
	c.executeInstruction()
	if c.pc != 0x0ffd {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestSed(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.sr |= Carry
	c.memory[c.pc] = 0xf8 // sed
	c.executeInstruction()
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	if c.sr&BCD != BCD {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestSei(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.sr |= Carry
	c.memory[c.pc] = 0x78 // sei
	c.executeInstruction()
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	if c.sr&Interrupts != Interrupts {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCld(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.sr |= BCD
	c.memory[c.pc] = 0xd8 // cld
	c.executeInstruction()
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	if c.sr&BCD == BCD {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestClc(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.sr |= Carry
	c.memory[c.pc] = 0x18 // clc
	c.executeInstruction()
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCli(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.sr |= Interrupts
	c.memory[c.pc] = 0x58 // cli
	c.executeInstruction()
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	if c.sr&Interrupts == Interrupts {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestClv(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.sr |= Overflow
	c.memory[c.pc] = 0xb8 // clv
	c.executeInstruction()
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	if c.sr&Overflow == Overflow {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestJsr(t *testing.T) {
	c := New()

	c.pc = 0x0300
	c.memory[c.pc] = 0x20
	c.memory[c.pc+1] = 0x34
	c.memory[c.pc+2] = 0x12
	c.executeInstruction()
	if c.pc != 0x1234 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ := c.disassemble(0x0300)
	t.Logf("%v\n", d)
	if c.sp != 0xfd {
		t.Fatalf("unexpected sp %02x", c.sp)
	}
	// check return address
	if c.memory[0x1ff] != 0x03 {
		t.Fatalf("invalid low byte %0x", c.memory[0x1ff])
	}
	if c.memory[0x1fe] != 0x02 {
		t.Fatal("invalid high byte %0x", c.memory[0x1fe])
	}
}

func TestJmp(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x4c // jmp $4030
	c.memory[c.pc+1] = 0x30
	c.memory[c.pc+2] = 0x40
	c.executeInstruction()
	if c.pc != 0x4030 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestJmpIndirect(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x6c // jmp ($4030)
	c.memory[c.pc+1] = 0x30
	c.memory[c.pc+2] = 0x40
	c.memory[0x4030] = 0xb0 // low byte
	c.memory[0x4031] = 0xf0 // high byte
	c.executeInstruction()
	if c.pc != 0xf0b0 {
		t.Fatalf("unexpected program counter %04x", c.pc)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAndIndirectX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0xf0
	c.x = 0x10

	c.memory[c.pc] = 0x21 // and(0x20,x)
	c.memory[c.pc+1] = 0x20

	c.memory[0x30] = 0x04 // low byte
	c.memory[0x31] = 0xd0 // high byte
	c.memory[0xd004] = 0x80
	c.executeInstruction()
	if c.a != 0x80 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAndZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.memory[c.pc] = 0x25 // and $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0xaa
	c.executeInstruction()
	if c.a != 0x00 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAndImmediate(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.memory[c.pc] = 0x29 // and #$aa
	c.memory[c.pc+1] = 0xaa
	c.executeInstruction()
	if c.a != 0x00 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAndAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.memory[c.pc] = 0x2d // and $2040
	c.memory[c.pc+1] = 0x40
	c.memory[c.pc+2] = 0x20
	c.memory[0x2040] = 0xaa
	c.executeInstruction()
	if c.a != 0x00 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAndIndirectY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.y = 0x10

	c.memory[c.pc] = 0x31 // and(0x20),y
	c.memory[c.pc+1] = 0x20

	c.memory[0x20] = 0x04 // low byte
	c.memory[0x21] = 0xd0 // high byte
	c.memory[0xd014] = 0xaa
	c.executeInstruction()
	if c.a != 0x00 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAndZPX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.x = 0x08
	c.memory[c.pc] = 0x35
	c.memory[c.pc+1] = 0x80
	c.memory[0x88] = 0xaa
	c.executeInstruction()
	if c.a != 0x00 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAndAbsoluteY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.y = 0xf0
	c.memory[c.pc] = 0x39
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0x80
	c.memory[0x80f0] = 0xaa
	c.executeInstruction()
	if c.a != 0x00 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAndAbsoluteX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.x = 0xf0
	c.memory[c.pc] = 0x3d
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0x80
	c.memory[0x80f0] = 0xaa
	c.executeInstruction()
	if c.a != 0x00 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestBitZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x00
	c.memory[c.pc] = 0x24 // bit $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0xff
	c.executeInstruction()
	if c.a != 0x00 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Overflow != Overflow {
		t.Fatalf("overflow unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test unset
	c.pc = 0x1000
	c.a = 0xf0
	c.memory[c.pc] = 0x24 // bit $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0x1f
	c.executeInstruction()
	if c.a != 0xf0 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Overflow == Overflow {
		t.Fatalf("overflow unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestBitAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x10            // test bit 4
	c.memory[c.pc] = 0x2c // bit $d020
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd020] = 0xff
	c.executeInstruction()
	if c.a != 0x10 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Overflow != Overflow {
		t.Fatalf("overflow unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestRolZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x26 // rol $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0x6e // 0b01101110
	c.sr |= Carry         // ends up in bit 0
	c.executeInstruction()
	if c.memory[0x80] != 0xdd {
		t.Fatalf("unexpected memory %0x", c.memory[0x80])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test 0
	c.pc = 0x1000
	c.memory[c.pc] = 0x26 // rol $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0x00
	c.sr = 0
	c.executeInstruction()
	if c.memory[0x80] != 0x00 {
		t.Fatalf("unexpected memory %0x", c.memory[0x80])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test carry
	c.pc = 0x1000
	c.memory[c.pc] = 0x26 // rol $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0x80
	c.sr = 0
	c.executeInstruction()
	if c.memory[0x80] != 0x00 {
		t.Fatalf("unexpected memory %0x", c.memory[0x80])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestRol(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x2a // rol $80
	c.a = 0x6e            // 0b01101110
	c.sr |= Carry         // ends up in bit 0
	c.executeInstruction()
	if c.a != 0xdd {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestRolAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x2e // rol $2112
	c.memory[c.pc+1] = 0x12
	c.memory[c.pc+2] = 0x21
	c.memory[0x2112] = 0x6e
	c.sr |= Carry // ends up in bit 0
	c.executeInstruction()
	if c.memory[0x2112] != 0xdd {
		t.Fatalf("unexpected memory %0x", c.memory[0x2112])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestRolZPX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x10
	c.memory[c.pc] = 0x36 // rol $70,x
	c.memory[c.pc+1] = 0x70
	c.memory[0x80] = 0x6e // 0b01101110
	c.sr |= Carry         // ends up in bit 0
	c.executeInstruction()
	if c.memory[0x80] != 0xdd {
		t.Fatalf("unexpected memory %0x", c.memory[0x80])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestRolAbsoluteX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x12
	c.memory[c.pc] = 0x3e // rol $2100
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0x21
	c.memory[0x2112] = 0x6e
	c.sr |= Carry // ends up in bit 0
	c.executeInstruction()
	if c.memory[0x2112] != 0xdd {
		t.Fatalf("unexpected memory %0x", c.memory[0x2112])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestSec(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x38 // sec
	c.executeInstruction()
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestRti(t *testing.T) {
	c := New()

	c.pc = 0x0300
	c.memory[0x1ff] = 0x12 // high byte
	c.memory[0x1fe] = 0x34 // low byte
	c.memory[0x1fd] = 0x5b // status register
	c.sp = 0xfc
	c.memory[c.pc] = 0x40 // rti
	c.executeInstruction()
	if c.pc != 0x1234 {
		t.Fatalf("unexpected program counter, %04x", c.pc)
	}
	if c.sr != 0x5b|0x20 {
		t.Fatalf("unexpected status register %0x", c.sr)
	}
	if c.sp != 0xff {
		t.Fatalf("unexpected sp %0x", c.sp)
	}
	d, _ := c.disassemble(0x0300)
	t.Logf("%v\n", d)
}

func TestRts(t *testing.T) {
	c := New()

	c.pc = 0x0300
	c.memory[0x1ff] = 0x12 // high byte
	c.memory[0x1fe] = 0x34 // low byte
	c.sp = 0xfd
	c.memory[c.pc] = 0x60 // rts
	c.executeInstruction()
	if c.pc != 0x1234+1 {
		t.Fatalf("unexpected program counter, %04x", c.pc)
	}
	// break flag is not set on cpu, only on stack
	if c.sr&Break != Break {
		t.Fatalf("unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x0300)
	t.Logf("%v\n", d)
}

func TestBrk(t *testing.T) {
	c := New()

	c.pc = 0x0300
	// set irq vector
	c.memory[0xffff] = 0x12 // high
	c.memory[0xfffe] = 0x34 // low
	c.memory[c.pc] = 0x00   // brk
	status := c.sr
	c.executeInstruction()
	if c.pc != 0x1234 {
		t.Fatalf("unexpected program counter, %04x", c.pc)
	}
	// break flag is not set on cpu, only on stack
	if c.sr&Break != Break {
		t.Fatalf("unexpected status register %0x", c.sr)
	}
	// check stack
	if c.memory[0x1ff] != 0x03 {
		t.Fatalf("unexpected high on stack %0x", c.memory[0x1ff])
	}
	// note the quirk of brk that adds one to the return address
	// meaning 0xc2 instead of 0xc1
	if c.memory[0x1fe] != 0x02 {
		t.Fatalf("unexpected low on stack %0x", c.memory[0x1fe])
	}
	if c.memory[0x1fd] != status {
		t.Fatalf("unexpected sr on stack %0x != %0x", c.memory[0x1fd], status)
	}
	d, _ := c.disassemble(0x300)
	t.Logf("%v\n", d)
}

func TestEorIndirectX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0xf5
	c.x = 0x10

	c.memory[c.pc] = 0x41 // eor(0x20,x)
	c.memory[c.pc+1] = 0x20

	c.memory[0x30] = 0x04 // low byte
	c.memory[0x31] = 0xd0 // high byte
	c.memory[0xd004] = 0xaa
	c.executeInstruction()
	if c.a != 0x5f {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestEorZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0xf5
	c.memory[c.pc] = 0x45 // eor $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0xaa
	c.executeInstruction()
	if c.a != 0x5f {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestEorImmediate(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0xf5
	c.memory[c.pc] = 0x49 // eor #$aa
	c.memory[c.pc+1] = 0xaa
	c.executeInstruction()
	if c.a != 0x5f {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestEorAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0xf5
	c.memory[c.pc] = 0x4d // eor $8000
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0x80
	c.memory[0x8000] = 0xaa
	c.executeInstruction()
	if c.a != 0x5f {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestEorIndirectY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0xf5
	c.y = 0x10

	c.memory[c.pc] = 0x51 // eor(0x20),y
	c.memory[c.pc+1] = 0x20

	c.memory[0x20] = 0x04 // low byte
	c.memory[0x21] = 0xd0 // high byte
	c.memory[0xd014] = 0xaa
	c.executeInstruction()
	if c.a != 0x5f {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestEorAbsoluteY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0xf5
	c.y = 0xf0
	c.memory[c.pc] = 0x59 // eor $8000,y
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0x80
	c.memory[0x80f0] = 0xaa
	c.executeInstruction()
	if c.a != 0x5f {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestEorAbsoluteX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0xf5
	c.x = 0xf0
	c.memory[c.pc] = 0x5d // eor $8000,x
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0x80
	c.memory[0x80f0] = 0xaa
	c.executeInstruction()
	if c.a != 0x5f {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestEorZPX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0xf5
	c.x = 0x08
	c.memory[c.pc] = 0x55 // eor $80,x
	c.memory[c.pc+1] = 0x80
	c.memory[0x88] = 0xaa
	c.executeInstruction()
	if c.a != 0x5f {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLsrZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x46 // lsr $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0x6e // 0b01101110
	c.executeInstruction()
	if c.memory[0x80] != 0x37 {
		t.Fatalf("unexpected memory %0x", c.memory[0x80])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test 0
	c.pc = 0x1000
	c.memory[c.pc] = 0x46 // lsr $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0x00
	c.sr = 0
	c.executeInstruction()
	if c.memory[0x80] != 0x00 {
		t.Fatalf("unexpected memory %0x", c.memory[0x80])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test carry
	c.pc = 0x1000
	c.memory[c.pc] = 0x46 // rol $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0x01
	c.sr = 0
	c.executeInstruction()
	if c.memory[0x80] != 0x00 {
		t.Fatalf("unexpected memory %0x", c.memory[0x80])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLsr(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x4a // lsr
	c.a = 0x6e            // 0b01101110
	c.executeInstruction()
	if c.a != 0x37 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLsrAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x4e // rol $2112
	c.memory[c.pc+1] = 0x12
	c.memory[c.pc+2] = 0x21
	c.memory[0x2112] = 0x6e
	c.executeInstruction()
	if c.memory[0x2112] != 0x37 {
		t.Fatalf("unexpected memory %0x", c.memory[0x2112])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLsrZPX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x10
	c.memory[c.pc] = 0x56 // lsr $70,x
	c.memory[c.pc+1] = 0x70
	c.memory[0x80] = 0x6e // 0b01101110
	c.executeInstruction()
	if c.memory[0x80] != 0x37 {
		t.Fatalf("unexpected memory %0x", c.memory[0x80])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLsrAbsoluteX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x12
	c.memory[c.pc] = 0x5e // lsr $2100
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0x21
	c.memory[0x2112] = 0x6e
	c.executeInstruction()
	if c.memory[0x2112] != 0x37 {
		t.Fatalf("unexpected memory %0x", c.memory[0x2112])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestRorZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x66 // ror $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0x6e // 0b01101110
	c.sr |= Carry         // ends up in bit 7
	c.executeInstruction()
	if c.memory[0x80] != 0xb7 {
		t.Fatalf("unexpected memory %0x", c.memory[0x80])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test 0
	c.pc = 0x1000
	c.memory[c.pc] = 0x66 // ror $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0x00
	c.sr = 0
	c.executeInstruction()
	if c.memory[0x80] != 0x00 {
		t.Fatalf("unexpected memory %0x", c.memory[0x80])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test carry
	c.pc = 0x1000
	c.memory[c.pc] = 0x66 // rol $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0x01
	c.sr = 0
	c.executeInstruction()
	if c.memory[0x80] != 0x00 {
		t.Fatalf("unexpected memory %0x", c.memory[0x80])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestRor(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x6a // ror $80
	c.a = 0x6e            // 0b01101110
	c.sr |= Carry         // ends up in bit 7
	c.executeInstruction()
	if c.a != 0xb7 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestRorAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x6e // ror $2112
	c.memory[c.pc+1] = 0x12
	c.memory[c.pc+2] = 0x21
	c.memory[0x2112] = 0x6e
	c.sr |= Carry // ends up in bit 7
	c.executeInstruction()
	if c.memory[0x2112] != 0xb7 {
		t.Fatalf("unexpected memory %0x", c.memory[0x2112])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestRorZPX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x10
	c.memory[c.pc] = 0x76 // ror $70,x
	c.memory[c.pc+1] = 0x70
	c.memory[0x80] = 0x6e // 0b01101110
	c.sr |= Carry         // ends up in bit 7
	c.executeInstruction()
	if c.memory[0x80] != 0xb7 {
		t.Fatalf("unexpected memory %0x", c.memory[0x80])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestRorAbsoluteX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x12
	c.memory[c.pc] = 0x7e // ror $2100
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0x21
	c.memory[0x2112] = 0x6e
	c.sr |= Carry // ends up in bit 7
	c.executeInstruction()
	if c.memory[0x2112] != 0xb7 {
		t.Fatalf("unexpected memory %0x", c.memory[0x2112])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestStaIndirectX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.x = 0x10

	c.memory[c.pc] = 0x81 // sta(0x20,x)
	c.memory[c.pc+1] = 0x20

	c.memory[0x30] = 0x04 // low byte
	c.memory[0x31] = 0xd0 // high byte
	c.executeInstruction()
	if c.memory[0xd004] != 0x55 {
		t.Fatalf("unexpected memory %0x", c.memory[0xd004])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestStaZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55

	c.memory[c.pc] = 0x85 // sta $20
	c.memory[c.pc+1] = 0x20

	c.executeInstruction()
	if c.memory[0x20] != 0x55 {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestStaAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55

	c.memory[c.pc] = 0x8d // sta $d020
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0

	c.executeInstruction()
	if c.memory[0xd020] != 0x55 {
		t.Fatalf("unexpected memory %0x", c.memory[0xd020])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestStaIndirectY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0xff
	c.y = 0x10

	c.memory[c.pc] = 0x91 // sta (0x20),y
	c.memory[c.pc+1] = 0x20

	c.memory[0x20] = 0x04 // low byte
	c.memory[0x21] = 0xd0 // high byte
	c.executeInstruction()
	if c.memory[0xd014] != 0xff {
		t.Fatalf("unexpected memory %0x", c.memory[0xd014])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestStaZPX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.x = 0x10

	c.memory[c.pc] = 0x95 // sta $10
	c.memory[c.pc+1] = 0x10

	c.executeInstruction()
	if c.memory[0x20] != 0x55 {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestStaAbsoluteY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.y = 0x10

	c.memory[c.pc] = 0x99 // sta $d020,y
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0

	c.executeInstruction()
	if c.memory[0xd030] != 0x55 {
		t.Fatalf("unexpected memory %0x", c.memory[0xd030])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestStaAbsoluteX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x55
	c.x = 0x10

	c.memory[c.pc] = 0x9d // sta $d020,x
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0

	c.executeInstruction()
	if c.memory[0xd030] != 0x55 {
		t.Fatalf("unexpected memory %0x", c.memory[0xd030])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestStyZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.y = 0x55

	c.memory[c.pc] = 0x84 // sty $20
	c.memory[c.pc+1] = 0x20

	c.executeInstruction()
	if c.memory[0x20] != 0x55 {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestStyAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.y = 0x55

	c.memory[c.pc] = 0x8c // sty $d020
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0

	c.executeInstruction()
	if c.memory[0xd020] != 0x55 {
		t.Fatalf("unexpected memory %0x", c.memory[0xd020])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestStyZPX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.y = 0x55
	c.x = 0x10

	c.memory[c.pc] = 0x94 // sty $10
	c.memory[c.pc+1] = 0x10

	c.executeInstruction()
	if c.memory[0x20] != 0x55 {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestStxZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x55

	c.memory[c.pc] = 0x86 // stx $20
	c.memory[c.pc+1] = 0x20

	c.executeInstruction()
	if c.memory[0x20] != 0x55 {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestStxAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x55

	c.memory[c.pc] = 0x8e // stx $d020
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0

	c.executeInstruction()
	if c.memory[0xd020] != 0x55 {
		t.Fatalf("unexpected memory %0x", c.memory[0xd020])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestStxZPY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x55
	c.y = 0x10

	c.memory[c.pc] = 0x96 // stx $10
	c.memory[c.pc+1] = 0x10

	c.executeInstruction()
	if c.memory[0x20] != 0x55 {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestDecZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xc6 // dec $20
	c.memory[c.pc+1] = 0x20
	c.memory[0x20] = 0x00
	c.executeInstruction()
	if c.memory[0x20] != 0xff {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// zero
	c.pc = 0x1000
	c.memory[c.pc] = 0xc6 // dec $20
	c.memory[c.pc+1] = 0x20
	c.memory[0x20] = 0x01
	c.executeInstruction()
	if c.memory[0x20] != 0x00 {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// non zero
	c.pc = 0x1000
	c.memory[c.pc] = 0xc6 // dec $20
	c.memory[c.pc+1] = 0x20
	c.memory[0x20] = 0x02
	c.executeInstruction()
	if c.memory[0x20] != 0x01 {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestDecAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xce // dec $d020
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd020] = 0x00
	c.executeInstruction()
	if c.memory[0xd020] != 0xff {
		t.Fatalf("unexpected memory %0x", c.memory[0xd020])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestDecZPX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x10

	c.memory[c.pc] = 0xd6 // dec $10,x
	c.memory[c.pc+1] = 0x10
	c.memory[0x20] = 0x30

	c.executeInstruction()
	if c.memory[0x20] != 0x2f {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestDecAbsoluteX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x10
	c.memory[c.pc] = 0xde // dec $d020,x
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd030] = 0x00
	c.executeInstruction()
	if c.memory[0xd030] != 0xff {
		t.Fatalf("unexpected memory %0x", c.memory[0xd030])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestIncZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xe6 // inc $20
	c.memory[c.pc+1] = 0x20
	c.memory[0x20] = 0xfe
	c.executeInstruction()
	if c.memory[0x20] != 0xff {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// zero
	c.pc = 0x1000
	c.memory[c.pc] = 0xe6 // inc $20
	c.memory[c.pc+1] = 0x20
	c.memory[0x20] = 0xff
	c.executeInstruction()
	if c.memory[0x20] != 0x00 {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// non zero
	c.pc = 0x1000
	c.memory[c.pc] = 0xe6 // inc $20
	c.memory[c.pc+1] = 0x20
	c.memory[0x20] = 0x00
	c.executeInstruction()
	if c.memory[0x20] != 0x01 {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestIncAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xee // inc $d020
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd020] = 0xfe
	c.executeInstruction()
	if c.memory[0xd020] != 0xff {
		t.Fatalf("unexpected memory %0x", c.memory[0xd020])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestIncZPX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x10
	c.memory[c.pc] = 0xf6 // inc $20,x
	c.memory[c.pc+1] = 0x20
	c.memory[0x30] = 0xfe
	c.executeInstruction()
	if c.memory[0x30] != 0xff {
		t.Fatalf("unexpected memory %0x", c.memory[0x30])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestIncAbsoluteX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x10
	c.memory[c.pc] = 0xfe // inc $d020,x
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd030] = 0xfe
	c.executeInstruction()
	if c.memory[0xd030] != 0xff {
		t.Fatalf("unexpected memory %0x", c.memory[0xd030])
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestDex(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x00
	c.memory[c.pc] = 0xca // dex
	c.executeInstruction()
	if c.x != 0xff {
		t.Fatalf("unexpected x %0x", c.x)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// zero
	c.pc = 0x1000
	c.x = 0x01
	c.memory[c.pc] = 0xca // dex
	c.executeInstruction()
	if c.x != 0x00 {
		t.Fatalf("unexpected x %0x", c.x)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// non zero
	c.pc = 0x1000
	c.x = 0x02
	c.memory[c.pc] = 0xca // dex
	c.executeInstruction()
	if c.x != 0x01 {
		t.Fatalf("unexpected x %0x", c.x)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestDey(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.y = 0x00
	c.memory[c.pc] = 0x88 // dey
	c.executeInstruction()
	if c.y != 0xff {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// zero
	c.pc = 0x1000
	c.y = 0x01
	c.memory[c.pc] = 0x88 // dey
	c.executeInstruction()
	if c.y != 0x00 {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// non zero
	c.pc = 0x1000
	c.y = 0x02
	c.memory[c.pc] = 0x88 // dey
	c.executeInstruction()
	if c.y != 0x01 {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestIny(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.y = 0xfe
	c.memory[c.pc] = 0xc8 // iny
	c.executeInstruction()
	if c.y != 0xff {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// zero
	c.pc = 0x1000
	c.y = 0xff
	c.memory[c.pc] = 0xc8 // iny
	c.executeInstruction()
	if c.y != 0x00 {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// non zero
	c.pc = 0x1000
	c.y = 0x01
	c.memory[c.pc] = 0xc8 // iny
	c.executeInstruction()
	if c.y != 0x02 {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestInx(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0xfe
	c.memory[c.pc] = 0xe8 // inx
	c.executeInstruction()
	if c.x != 0xff {
		t.Fatalf("unexpected x %0x", c.x)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// zero
	c.pc = 0x1000
	c.x = 0xff
	c.memory[c.pc] = 0xe8 // inx
	c.executeInstruction()
	if c.x != 0x00 {
		t.Fatalf("unexpected x %0x", c.x)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// non zero
	c.pc = 0x1000
	c.x = 0x01
	c.memory[c.pc] = 0xe8 // inx
	c.executeInstruction()
	if c.x != 0x02 {
		t.Fatalf("unexpected x %0x", c.x)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestTxa(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0xff
	c.memory[c.pc] = 0x8a // txa
	c.executeInstruction()
	if c.a != c.x {
		t.Fatalf("unexpected a %0x", c.a)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// zero
	c.pc = 0x1000
	c.x = 0x00
	c.memory[c.pc] = 0x8a // txa
	c.executeInstruction()
	if c.a != c.x {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// non zero
	c.pc = 0x1000
	c.x = 0x01
	c.memory[c.pc] = 0x8a // txa
	c.executeInstruction()
	if c.a != c.x {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestTya(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.y = 0xff
	c.memory[c.pc] = 0x98 // tya
	c.executeInstruction()
	if c.a != c.y {
		t.Fatalf("unexpected a %0x", c.a)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// zero
	c.pc = 0x1000
	c.y = 0x00
	c.memory[c.pc] = 0x98 // tya
	c.executeInstruction()
	if c.a != c.y {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// non zero
	c.pc = 0x1000
	c.y = 0x01
	c.memory[c.pc] = 0x98 // tya
	c.executeInstruction()
	if c.a != c.y {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestTay(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0xff
	c.memory[c.pc] = 0xa8 // tay
	c.executeInstruction()
	if c.a != c.y {
		t.Fatalf("unexpected a %0x", c.y)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// zero
	c.pc = 0x1000
	c.a = 0x00
	c.memory[c.pc] = 0xa8 // tay
	c.executeInstruction()
	if c.a != c.y {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// non zero
	c.pc = 0x1000
	c.a = 0x01
	c.memory[c.pc] = 0xa8 // tay
	c.executeInstruction()
	if c.a != c.y {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestTax(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0xff
	c.memory[c.pc] = 0xaa // tax
	c.executeInstruction()
	if c.a != c.x {
		t.Fatalf("unexpected a %0x", c.x)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// zero
	c.pc = 0x1000
	c.a = 0x00
	c.memory[c.pc] = 0xaa // tax
	c.executeInstruction()
	if c.a != c.x {
		t.Fatalf("unexpected y %0x", c.x)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// non zero
	c.pc = 0x1000
	c.a = 0x01
	c.memory[c.pc] = 0xaa // tax
	c.executeInstruction()
	if c.a != c.x {
		t.Fatalf("unexpected y %0x", c.x)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestTxs(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0xa0

	c.memory[c.pc] = 0x9a // txs

	c.executeInstruction()
	if c.sp != 0xa0 {
		t.Fatalf("unexpected stack pointer %0x", c.sp)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestTsx(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.sp = 0xa0

	c.memory[c.pc] = 0xba // tsx

	c.executeInstruction()
	if c.x != 0xa0 {
		t.Fatalf("unexpected x %0x", c.x)
	}
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdy(t *testing.T) {
	c := New()

	c.pc = 0x1000

	c.memory[c.pc] = 0xa0 // ldy #$8f
	c.memory[c.pc+1] = 0x8f

	c.executeInstruction()
	if c.y != 0x8f {
		t.Fatalf("unexpected memory %0x", c.y)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdyZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xa4 // ldy $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0xaa
	c.executeInstruction()
	if c.y != 0xaa {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdyAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xac // ldy $8000
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0x80
	c.memory[0x8000] = 0xaa
	c.executeInstruction()
	if c.y != 0xaa {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdyZPX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x10

	c.memory[c.pc] = 0xb4 // ldy $10,x
	c.memory[c.pc+1] = 0x10
	c.memory[0x20] = 0x01

	c.executeInstruction()
	if c.y != 0x01 {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdyAbsoluteX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x10

	c.memory[c.pc] = 0xbc // ldy $d020,x
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd030] = 0x55

	c.executeInstruction()
	if c.y != 0x55 {
		t.Fatalf("unexpected memory %0x", c.y)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdx(t *testing.T) {
	c := New()

	c.pc = 0x1000

	c.memory[c.pc] = 0xa2 // ldx #$8f
	c.memory[c.pc+1] = 0x8f

	c.executeInstruction()
	if c.x != 0x8f {
		t.Fatalf("unexpected memory %0x", c.y)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdxZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xa6 // ldx $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0xaa
	c.executeInstruction()
	if c.x != 0xaa {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdxAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xae // ldx $8000
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0x80
	c.memory[0x8000] = 0xaa
	c.executeInstruction()
	if c.x != 0xaa {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdxZPY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.y = 0x10

	c.memory[c.pc] = 0xb6 // ldx $10,y
	c.memory[c.pc+1] = 0x10
	c.memory[0x20] = 0x01

	c.executeInstruction()
	if c.x != 0x01 {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdxAbsoluteY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.y = 0x10

	c.memory[c.pc] = 0xbe // ldx $d020,y
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd030] = 0x55

	c.executeInstruction()
	if c.x != 0x55 {
		t.Fatalf("unexpected memory %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdaIndirectX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x10

	c.memory[c.pc] = 0xa1 // lda (0x20,x)
	c.memory[c.pc+1] = 0x20

	c.memory[0x30] = 0x04 // low byte
	c.memory[0x31] = 0xd0 // high byte
	c.memory[0xd004] = 0xaa
	c.executeInstruction()
	if c.a != 0xaa {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdaZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xa5 // lda $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0xaa
	c.executeInstruction()
	if c.a != 0xaa {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLda(t *testing.T) {
	c := New()

	c.pc = 0x1000

	c.memory[c.pc] = 0xa9 // lda #$8f
	c.memory[c.pc+1] = 0x8f

	c.executeInstruction()
	if c.a != 0x8f {
		t.Fatalf("unexpected memory %0x", c.y)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdaAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xad // lda $8000
	c.memory[c.pc+1] = 0x00
	c.memory[c.pc+2] = 0x80
	c.memory[0x8000] = 0xaa
	c.executeInstruction()
	if c.a != 0xaa {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdaIndirectY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.y = 0x10

	c.memory[c.pc] = 0xb1 // lda (0x20),y
	c.memory[c.pc+1] = 0x20

	c.memory[0x20] = 0x04 // low byte
	c.memory[0x21] = 0xd0 // high byte
	c.memory[0xd014] = 0xff

	c.executeInstruction()
	if c.a != 0xff {
		t.Fatalf("unexpected memory %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdaZPX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x10

	c.memory[c.pc] = 0xb5 // lda $10,x
	c.memory[c.pc+1] = 0x10
	c.memory[0x20] = 0x01

	c.executeInstruction()
	if c.a != 0x01 {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdaAbsoluteY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.y = 0x10

	c.memory[c.pc] = 0xb9 // lda $d020,y
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd030] = 0x55

	c.executeInstruction()
	if c.a != 0x55 {
		t.Fatalf("unexpected memory %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestLdaAbsoluteX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x10

	c.memory[c.pc] = 0xbd // lda $d020,x
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd030] = 0x55

	c.executeInstruction()
	if c.a != 0x55 {
		t.Fatalf("unexpected memory %0x", c.y)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCpy(t *testing.T) {
	c := New()

	// identical
	c.pc = 0x1000
	c.y = 0x40
	c.memory[c.pc] = 0xc0 // cpy #$40
	c.memory[c.pc+1] = 0x40
	c.executeInstruction()
	if c.y != 0x40 {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// register larger
	c.pc = 0x1000
	c.y = 0x40
	c.memory[c.pc] = 0xc0 // cpy #$41
	c.memory[c.pc+1] = 0x41
	c.executeInstruction()
	if c.y != 0x40 {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// register smaller
	c.pc = 0x1000
	c.y = 0x40
	c.memory[c.pc] = 0xc0 // cpy #$3f
	c.memory[c.pc+1] = 0x3f
	c.executeInstruction()
	if c.y != 0x40 {
		t.Fatalf("unexpected y %0x", c.y)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCpyZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xc4 // cpy $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0x40
	c.y = 0x40
	c.executeInstruction()
	if c.y != 0x40 {
		t.Fatalf("unexpected accumulator %0x", c.y)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCpyAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.y = 0x40

	c.memory[c.pc] = 0xcc // cpy $d020
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd020] = 0x40

	c.executeInstruction()
	if c.y != 0x40 {
		t.Fatalf("unexpected memory %0x", c.y)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCpx(t *testing.T) {
	c := New()

	// identical
	c.pc = 0x1000
	c.x = 0x40
	c.memory[c.pc] = 0xe0 // cpx #$40
	c.memory[c.pc+1] = 0x40
	c.executeInstruction()
	if c.x != 0x40 {
		t.Fatalf("unexpected x %0x", c.x)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// register larger
	c.pc = 0x1000
	c.x = 0x40
	c.memory[c.pc] = 0xe0 // cpx #$41
	c.memory[c.pc+1] = 0x41
	c.executeInstruction()
	if c.x != 0x40 {
		t.Fatalf("unexpected x %0x", c.x)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// register smaller
	c.pc = 0x1000
	c.x = 0x40
	c.memory[c.pc] = 0xe0 // cpx #$3f
	c.memory[c.pc+1] = 0x3f
	c.executeInstruction()
	if c.x != 0x40 {
		t.Fatalf("unexpected x %0x", c.x)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCpxZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xe4 // cpx $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0x40
	c.x = 0x40
	c.executeInstruction()
	if c.x != 0x40 {
		t.Fatalf("unexpected x %0x", c.y)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCpxAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x40

	c.memory[c.pc] = 0xec // cx $d020
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd020] = 0x40

	c.executeInstruction()
	if c.x != 0x40 {
		t.Fatalf("unexpected x %0x", c.x)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCmpIndirectX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x10

	c.memory[c.pc] = 0xc1 // cmp (0x20,x)
	c.memory[c.pc+1] = 0x20

	c.memory[0x30] = 0x04 // low byte
	c.memory[0x31] = 0xd0 // high byte
	c.a = 0x40
	c.memory[0xd004] = 0x40
	c.executeInstruction()
	if c.a != 0x40 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCmpZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xc5 // cpy $80
	c.memory[c.pc+1] = 0x80
	c.memory[0x80] = 0x40
	c.a = 0x40
	c.executeInstruction()
	if c.a != 0x40 {
		t.Fatalf("unexpected accumulator %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCmp(t *testing.T) {
	c := New()

	// identical
	c.pc = 0x1000
	c.a = 0x40
	c.memory[c.pc] = 0xc9 // cmp #$40
	c.memory[c.pc+1] = 0x40
	c.executeInstruction()
	if c.a != 0x40 {
		t.Fatalf("unexpected a %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// register larger
	c.pc = 0x1000
	c.a = 0x40
	c.memory[c.pc] = 0xc9 // cmp #$41
	c.memory[c.pc+1] = 0x41
	c.executeInstruction()
	if c.a != 0x40 {
		t.Fatalf("unexpected a %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative != Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// register smaller
	c.pc = 0x1000
	c.a = 0x40
	c.memory[c.pc] = 0xc9 // cmp #$3f
	c.memory[c.pc+1] = 0x3f
	c.executeInstruction()
	if c.a != 0x40 {
		t.Fatalf("unexpected memory %0x", c.a)
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCmpAbsolute(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x40

	c.memory[c.pc] = 0xcd // cmp $d020
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd020] = 0x40

	c.executeInstruction()
	if c.a != 0x40 {
		t.Fatalf("unexpected memory %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCmpIndirectY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x40
	c.y = 0x10

	c.memory[c.pc] = 0xd1 // cmp (0x20),y
	c.memory[c.pc+1] = 0x20

	c.memory[0x20] = 0x04 // low byte
	c.memory[0x21] = 0xd0 // high byte
	c.memory[0xd014] = 0x40
	c.executeInstruction()
	if c.memory[0xd014] != c.a {
		t.Fatalf("unexpected memory %0x", c.memory[0xd014])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCmpZPX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.x = 0x10
	c.a = 0x40

	c.memory[c.pc] = 0xd5 // cmp $10,x
	c.memory[c.pc+1] = 0x10
	c.memory[0x20] = 0x40

	c.executeInstruction()
	if c.a != c.memory[0x20] {
		t.Fatalf("unexpected memory %0x", c.memory[0x20])
	}
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCmpAbsoluteY(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x40
	c.y = 0x10

	c.memory[c.pc] = 0xd9 // cmp $d020,y
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd030] = 0x40

	c.executeInstruction()
	if c.a != c.memory[0xd030] {
		t.Fatalf("unexpected a %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestCmpAbsoluteX(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.a = 0x40
	c.x = 0x10

	c.memory[c.pc] = 0xdd // cmp $d020,x
	c.memory[c.pc+1] = 0x20
	c.memory[c.pc+2] = 0xd0
	c.memory[0xd030] = 0x40

	c.executeInstruction()
	if c.a != c.memory[0xd030] {
		t.Fatalf("unexpected a %0x", c.a)
	}
	if c.pc != 0x1003 {
		t.Fatalf("unexpected program counter")
	}
	if c.sr&Negative == Negative {
		t.Fatalf("negative unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("zero unexpected status register %0x", c.sr)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestNop(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xea // nop

	c.executeInstruction()
	if c.pc != 0x1001 {
		t.Fatalf("unexpected program counter")
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestSbcImmediate(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xe9 // sbc #$01
	c.memory[c.pc+1] = 0x01
	c.a = 0x42
	c.sr |= Carry

	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.a != 0x41 {
		t.Fatalf("unexpected a %x", c.a)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// without carry
	c = New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xe9 // sbc #$01
	c.memory[c.pc+1] = 0x01
	c.a = 0x42

	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.a != 0x40 {
		t.Fatalf("unexpected a %x", c.a)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// negative no carry
	c = New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xe9 // sbc #$43
	c.memory[c.pc+1] = 0x43
	c.a = 0x42
	c.sr |= Carry

	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.a != 0xff {
		t.Fatalf("unexpected a %x", c.a)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// zero
	c = New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xe9 // sbc #$42
	c.memory[c.pc+1] = 0x42
	c.a = 0x42
	c.sr |= Carry

	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.a != 0x00 {
		t.Fatalf("unexpected a %x", c.a)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

}

func TestSbcDecimal(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xe9 // sbc #$03
	c.memory[c.pc+1] = 0x03
	c.a = 0x32
	c.sr |= BCD

	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.a != 0x28 {
		t.Fatalf("unexpected a %x", c.a)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	if c.sr&Zero == Zero {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestSbcZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0xe5 // sbc $80
	c.memory[c.pc+1] = 0x80
	c.a = 0x42
	c.memory[0x80] = 0x12
	c.sr |= Carry

	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.a != 0x30 {
		t.Fatalf("unexpected a %x", c.a)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAdcImmediate(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x69 // adc #$53
	c.memory[c.pc+1] = 0x53
	c.a = 0x42

	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.a != 0x95 {
		t.Fatalf("unexpected a %x", c.a)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test carry
	c = New()
	c.pc = 0x1000
	c.memory[c.pc] = 0x69 // adc #$53
	c.memory[c.pc+1] = 0x53
	c.a = 0xc0

	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.a != 0x13 {
		t.Fatalf("unexpected a %x", c.a)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test carry over
	c = New()
	c.pc = 0x1000
	c.memory[c.pc] = 0x69 // adc #$04
	c.memory[c.pc+1] = 0x04
	c.sr |= Carry
	c.a = 0x05

	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.a != 0x0a {
		t.Fatalf("unexpected a %x", c.a)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test carry with overflow
	c = New()
	c.pc = 0x1000
	c.memory[c.pc] = 0x69 // adc #$d0
	c.memory[c.pc+1] = 0xd0
	c.a = 0x90

	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.a != 0x60 {
		t.Fatalf("unexpected a %x", c.a)
	}
	if c.sr&Carry != Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	if c.sr&Overflow != Overflow {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test zero
	c = New()
	c.pc = 0x1000
	c.memory[c.pc] = 0x69 // adc #$00
	c.memory[c.pc+1] = 0x00
	c.a = 0x00

	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.a != 0x00 {
		t.Fatalf("unexpected a %x", c.a)
	}
	if c.sr&Zero != Zero {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

	// test negative
	c = New()
	c.pc = 0x1000
	c.memory[c.pc] = 0x69 // adc #$f7
	c.memory[c.pc+1] = 0xf7
	c.a = 0x00

	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.a != 0xf7 {
		t.Fatalf("unexpected a %x", c.a)
	}
	if c.sr&Negative != Negative {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ = c.disassemble(0x1000)
	t.Logf("%v\n", d)

}

func TestAdcZP(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x65 // adc $80
	c.memory[c.pc+1] = 0x80
	c.a = 0x42
	c.memory[0x80] = 0x12

	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.a != 0x54 {
		t.Fatalf("unexpected a %x", c.a)
	}
	if c.sr&Carry == Carry {
		t.Fatalf("carry unexpected status register %0x", c.sr)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestAdcDecimal(t *testing.T) {
	c := New()

	c.pc = 0x1000
	c.memory[c.pc] = 0x69 // adc #$28
	c.memory[c.pc+1] = 0x28
	c.a = 0x19
	c.sr |= BCD

	c.executeInstruction()
	if c.pc != 0x1002 {
		t.Fatalf("unexpected program counter")
	}
	if c.a != 0x47 {
		t.Fatalf("unexpected a %x", c.a)
	}
	d, _ := c.disassemble(0x1000)
	t.Logf("%v\n", d)
}

func TestKlausDormann6502(t *testing.T) {
	c := New()
	f, err := os.Open("test/6502_functional_test.bin")
	if err != nil {
		t.Fatal(err)
	}
	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if fi.Size() > int64(len(c.memory)) {
		t.Fatal("invalid ram image size")
	}
	_, err = f.Read(c.memory)
	if err != nil {
		t.Fatal(err)
	}

	c.pc = 0x0400
	prevPC := uint16(0x0400)

	var instructions uint64
	for {
		//d, _ := c.disassemble(c.pc)
		//t.Logf("%v\n", d)
		c.executeInstruction()
		//fmt.Printf("%v\n", c.snapshot())
		//fmt.Printf("A: $%02x X: $%02x Y: $%02x SR: $%02x PC: $%04x SP: $%02x"+
		//	" %02x %02x %02x %02x %02x %02x %02x %02x\n",
		//	c.a,
		//	c.x,
		//	c.y,
		//	c.sr,
		//	c.pc,
		//	c.sp,
		//	c.memory[0x0a],
		//	c.memory[0x0b],
		//	c.memory[0x0c],
		//	c.memory[0x0d],
		//	c.memory[0x0e],
		//	c.memory[0x0f],
		//	c.memory[0x10],
		//	c.memory[0x11])
		instructions++
		if c.pc == prevPC {
			if c.pc != 0x3399 {
				t.Fatalf("loop detected at PC 0x%04X.", c.pc)
			}
			t.Logf("Klaus Dormann's 6502 functional tests passed.")
			t.Logf("instructions run: %v cycles: %v",
				instructions,
				c.cycles)
			return
		}

		prevPC = c.pc
	}
}
