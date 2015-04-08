package main

import "fmt"

const (
	Negative   byte = 1 << 7 // N
	Overflow   byte = 1 << 6 // V
	Unused     byte = 1 << 5
	Break      byte = 1 << 4 // B
	BCD        byte = 1 << 3 // D
	Interrupts byte = 1 << 2 // I
	Zero       byte = 1 << 1 // Z
	Carry      byte = 1 << 0 // C
)

type mode int

const (
	none mode = iota
	immediate
	implied
	indirect
	accumulator
	relative
	absolute
	absoluteX
	absoluteY
	zeroPage
	zeroPageX
	zeroPageY
	zeroPageIndirectX
	zeroPageIndirectY
)

type opcode struct {
	mnemonic    string
	noBytes     byte
	noCycles    uint64
	extraCycles uint64
	mode        mode
}

var (
	invalidOpcode = opcode{
		mnemonic: "???",
	}

	opcodes = []opcode{
		// 0x00
		{
			mnemonic: "BRK",
			noBytes:  1,
			noCycles: 7,
			mode:     implied,
		},
		// 0x01
		{
			mnemonic:    "ORA",
			mode:        zeroPageIndirectX,
			noBytes:     2,
			noCycles:    5,
			extraCycles: 1,
		},
		// 0x02
		invalidOpcode,
		// 0x03
		invalidOpcode,
		// 0x04
		invalidOpcode,
		// 0x05
		{
			mnemonic: "ORA",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 2,
		},
		// 0x06
		{
			mnemonic: "ASL",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 5,
		},
		// 0x07
		invalidOpcode,
		// 0x08
		{
			mnemonic: "PHP",
			mode:     implied,
			noBytes:  1,
			noCycles: 3,
		},
		// 0x09
		{
			mnemonic: "ORA",
			mode:     immediate,
			noBytes:  2,
			noCycles: 2,
		},
		// 0x0a
		{
			mnemonic: "ASL",
			mode:     accumulator,
			noBytes:  1,
			noCycles: 2,
		},
		// 0x0b
		invalidOpcode,
		// 0x0c
		invalidOpcode,
		// 0x0d
		{
			mnemonic: "ORA",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0x0e
		{
			mnemonic: "ASL",
			mode:     absolute,
			noBytes:  3,
			noCycles: 6,
		},
		// 0x0f
		invalidOpcode,
		// 0x10
		{
			mnemonic:    "BPL",
			mode:        relative,
			noBytes:     2,
			noCycles:    2,
			extraCycles: 1, // XXX or 2
		},
		// 0x11
		{
			mnemonic:    "ORA",
			mode:        zeroPageIndirectY,
			noBytes:     2,
			noCycles:    5,
			extraCycles: 1,
		},
		// 0x12
		invalidOpcode,
		// 0x13
		invalidOpcode,
		// 0x14
		invalidOpcode,
		// 0x15
		{
			mnemonic: "ORA",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 3,
		},
		// 0x16
		{
			mnemonic: "ASL",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 5,
		},
		// 0x17
		invalidOpcode,
		// 0x18
		{
			mnemonic: "CLC",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0x19
		{
			mnemonic:    "ORA",
			mode:        absoluteY,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0x1a
		invalidOpcode,
		// 0x1b
		invalidOpcode,
		// 0x1c
		invalidOpcode,
		// 0x1d
		{
			mnemonic:    "ORA",
			mode:        absoluteX,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0x1e
		{
			mnemonic: "ASL",
			mode:     absoluteX,
			noBytes:  3,
			noCycles: 7,
		},
		// 0x1f
		invalidOpcode,
		// 0x20
		{
			mnemonic: "JSR",
			mode:     absolute,
			noBytes:  3,
			noCycles: 6,
		},
		// 0x21
		{
			mnemonic: "AND",
			mode:     zeroPageIndirectX,
			noBytes:  2,
			noCycles: 6,
		},
		// 0x22
		invalidOpcode,
		// 0x23
		invalidOpcode,
		// 0x24
		{
			mnemonic: "BIT",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 3,
		},
		// 0x25
		{
			mnemonic: "AND",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 2,
		},
		// 0x26
		{
			mnemonic: "ROL",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 5,
		},
		// 0x27
		invalidOpcode,
		// 0x28
		{
			mnemonic: "PLP",
			mode:     implied,
			noBytes:  1,
			noCycles: 4,
		},
		// 0x29
		{
			mnemonic: "AND",
			mode:     immediate,
			noBytes:  2,
			noCycles: 2,
		},
		// 0x2a
		{
			mnemonic: "ROL",
			mode:     accumulator,
			noBytes:  1,
			noCycles: 2,
		},
		// 0x2b
		invalidOpcode,
		// 0x2c
		{
			mnemonic: "BIT",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0x2d
		{
			mnemonic: "AND",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0x2e
		{
			mnemonic: "ROL",
			mode:     absolute,
			noBytes:  3,
			noCycles: 6,
		},
		// 0x2f
		invalidOpcode,
		// 0x30
		{
			mnemonic:    "BMI",
			mode:        relative,
			noBytes:     2,
			noCycles:    2,
			extraCycles: 1, // XXX or 2
		},
		// 0x31
		{
			mnemonic:    "AND",
			mode:        zeroPageIndirectY,
			noBytes:     2,
			noCycles:    5,
			extraCycles: 1,
		},
		// 0x32
		invalidOpcode,
		// 0x33
		invalidOpcode,
		// 0x34
		invalidOpcode,
		// 0x35
		{
			mnemonic: "AND",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 3,
		},
		// 0x36
		{
			mnemonic: "ROL",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 6,
		},
		// 0x37
		invalidOpcode,
		// 0x38
		{
			mnemonic: "SEC",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0x39
		{
			mnemonic:    "AND",
			mode:        absoluteY,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0x3a
		invalidOpcode,
		// 0x3b
		invalidOpcode,
		// 0x3c
		invalidOpcode,
		// 0x3d
		{
			mnemonic:    "AND",
			mode:        absoluteX,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0x3e
		{
			mnemonic: "ROL",
			mode:     absoluteX,
			noBytes:  3,
			noCycles: 7,
		},
		// 0x3f
		invalidOpcode,
		// 0x40
		{
			mnemonic: "RTI",
			mode:     implied,
			noBytes:  1,
			noCycles: 6,
		},
		// 0x41
		{
			mnemonic: "EOR",
			mode:     zeroPageIndirectX,
			noBytes:  2,
			noCycles: 6,
		},
		// 0x42
		invalidOpcode,
		// 0x43
		invalidOpcode,
		// 0x44
		invalidOpcode,
		// 0x45
		{
			mnemonic: "EOR",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 3,
		},
		// 0x46
		{
			mnemonic: "LSR",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 5,
		},
		// 0x47
		invalidOpcode,
		// 0x48
		{
			mnemonic: "PHA",
			mode:     implied,
			noBytes:  1,
			noCycles: 3,
		},
		// 0x49
		{
			mnemonic: "EOR",
			mode:     immediate,
			noBytes:  2,
			noCycles: 2,
		},
		// 0x4a
		{
			mnemonic: "LSR",
			mode:     accumulator,
			noBytes:  1,
			noCycles: 2,
		},
		// 0x4b
		invalidOpcode,
		// 0x4c
		{
			mnemonic: "JMP",
			mode:     absolute,
			noBytes:  3,
			noCycles: 3,
		},
		// 0x4d
		{
			mnemonic: "EOR",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0x4e
		{
			mnemonic: "LSR",
			mode:     absolute,
			noBytes:  3,
			noCycles: 6,
		},
		// 0x4f
		invalidOpcode,
		// 0x50
		{
			mnemonic:    "BVC",
			mode:        relative,
			noBytes:     2,
			noCycles:    2,
			extraCycles: 1, // XXX or 2
		},
		// 0x51
		{
			mnemonic:    "EOR",
			mode:        zeroPageIndirectY,
			noBytes:     2,
			noCycles:    5,
			extraCycles: 1,
		},
		// 0x52
		invalidOpcode,
		// 0x53
		invalidOpcode,
		// 0x54
		invalidOpcode,
		// 0x55
		{
			mnemonic: "EOR",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 4,
		},
		// 0x56
		{
			mnemonic: "EOR",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 6,
		},
		// 0x57
		invalidOpcode,
		// 0x58
		{
			mnemonic: "CLI",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0x59
		{
			mnemonic:    "EOR",
			mode:        absoluteY,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0x5a
		invalidOpcode,
		// 0x5b
		invalidOpcode,
		// 0x5c
		invalidOpcode,
		// 0x5d
		{
			mnemonic:    "EOR",
			mode:        absoluteX,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0x5e
		{
			mnemonic: "LSR",
			mode:     absoluteX,
			noBytes:  3,
			noCycles: 7,
		},
		// 0x5f
		invalidOpcode,
		// 0x60
		{
			mnemonic: "RTS",
			mode:     implied,
			noBytes:  1,
			noCycles: 6,
		},
		// 0x61
		{
			mnemonic: "ADC",
			mode:     zeroPageIndirectX,
			noBytes:  2,
			noCycles: 6,
		},
		// 0x62
		invalidOpcode,
		// 0x63
		invalidOpcode,
		// 0x64
		invalidOpcode,
		// 0x65
		{
			mnemonic: "ADC",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 3,
		},
		// 0x66
		{
			mnemonic: "ROR",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 5,
		},
		// 0x67
		invalidOpcode,
		// 0x68
		{
			mnemonic: "PLA",
			mode:     implied,
			noBytes:  1,
			noCycles: 4,
		},
		// 0x69
		{
			mnemonic: "ADC",
			mode:     immediate,
			noBytes:  2,
			noCycles: 2,
		},
		// 0x6a
		{
			mnemonic: "ROR",
			mode:     accumulator,
			noBytes:  1,
			noCycles: 2,
		},
		// 0x6b
		invalidOpcode,
		// 0x6c
		{
			mnemonic: "JMP",
			mode:     indirect,
			noBytes:  3,
			noCycles: 5,
		},
		// 0x6d
		{
			mnemonic: "ADC",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0x6e
		{
			mnemonic: "ROR",
			mode:     absolute,
			noBytes:  3,
			noCycles: 6,
		},
		// 0x6f
		invalidOpcode,
		// 0x70
		{
			mnemonic:    "BVS",
			mode:        relative,
			noBytes:     2,
			noCycles:    2,
			extraCycles: 1, // XXX or 2
		},
		// 0x71
		{
			mnemonic:    "ADC",
			mode:        zeroPageIndirectY,
			noBytes:     2,
			noCycles:    5,
			extraCycles: 1,
		},
		// 0x72
		invalidOpcode,
		// 0x73
		invalidOpcode,
		// 0x74
		invalidOpcode,
		// 0x75
		{
			mnemonic: "ADC",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 4,
		},
		// 0x76
		{
			mnemonic: "ROR",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 6,
		},
		// 0x77
		invalidOpcode,
		// 0x78
		{
			mnemonic: "SEI",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0x79
		{
			mnemonic:    "ADC",
			mode:        absoluteY,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0x7a
		invalidOpcode,
		// 0x7a
		invalidOpcode,
		// 0x7c
		invalidOpcode,
		// 0x7d
		{
			mnemonic:    "ADC",
			mode:        absoluteX,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0x7e
		{
			mnemonic: "ROR",
			mode:     absoluteX,
			noBytes:  3,
			noCycles: 7,
		},
		// 0x7f
		invalidOpcode,
		// 0x80
		invalidOpcode,
		// 0x81
		{
			mnemonic: "STA",
			mode:     zeroPageIndirectX,
			noBytes:  2,
			noCycles: 6,
		},
		// 0x82
		invalidOpcode,
		// 0x83
		invalidOpcode,
		// 0x84
		{
			mnemonic: "STY",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 3,
		},
		// 0x85
		{
			mnemonic: "STA",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 3,
		},
		// 0x86
		{
			mnemonic: "STX",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 3,
		},
		// 0x87
		invalidOpcode,
		// 0x88
		{
			mnemonic: "DEY",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0x89
		invalidOpcode,
		// 0x8a
		{
			mnemonic: "TXA",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0x8b
		invalidOpcode,
		// 0x8c
		{
			mnemonic: "STY",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0x8d
		{
			mnemonic: "STA",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0x8e
		{
			mnemonic: "STX",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0x8f
		invalidOpcode,
		// 0x90
		{
			mnemonic: "BCC",
			mode:     relative,
			noBytes:  2,
			noCycles: 2, // XXX or 2
		},
		// 0x91
		{
			mnemonic: "STA",
			mode:     zeroPageIndirectY,
			noBytes:  2,
			noCycles: 6,
		},
		// 0x92
		invalidOpcode,
		// 0x93
		invalidOpcode,
		// 0x94
		{
			mnemonic: "STY",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 4,
		},
		// 0x95
		{
			mnemonic: "STA",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 4,
		},
		// 0x96
		{
			mnemonic: "STX",
			mode:     zeroPageY,
			noBytes:  2,
			noCycles: 4,
		},
		// 0x97
		invalidOpcode,
		// 0x98
		{
			mnemonic: "TYA",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0x99
		{
			mnemonic: "STA",
			mode:     absoluteY,
			noBytes:  3,
			noCycles: 5,
		},
		// 0x9a
		{
			mnemonic: "TXS",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0x9b
		invalidOpcode,
		// 0x9c
		invalidOpcode,
		// 0x9d
		{
			mnemonic: "STA",
			mode:     absoluteX,
			noBytes:  3,
			noCycles: 5,
		},
		// 0x9e
		invalidOpcode,
		// 0x9f
		invalidOpcode,
		// 0xa0
		{
			mnemonic: "LDY",
			mode:     immediate,
			noBytes:  2,
			noCycles: 2,
		},
		// 0xa1
		{
			mnemonic: "LDA",
			mode:     zeroPageIndirectX,
			noBytes:  2,
			noCycles: 6,
		},
		// 0xa2
		{
			mnemonic: "LDX",
			mode:     immediate,
			noBytes:  2,
			noCycles: 2,
		},
		// 0xa3
		invalidOpcode,
		// 0xa4
		{
			mnemonic: "LDY",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 3,
		},
		// 0xa5
		{
			mnemonic: "LDA",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 3,
		},
		// 0xa6
		{
			mnemonic: "LDX",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 3,
		},
		// 0xa7
		invalidOpcode,
		// 0xa8
		{
			mnemonic: "TAY",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0xa9
		{
			mnemonic: "LDA",
			mode:     immediate,
			noBytes:  2,
			noCycles: 2,
		},
		// 0xaa
		{
			mnemonic: "TAX",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0xab
		invalidOpcode,
		// 0xac
		{
			mnemonic: "LDY",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0xad
		{
			mnemonic: "LDA",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0xae
		{
			mnemonic: "LDX",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0xaf
		invalidOpcode,
		// 0xb0
		{
			mnemonic:    "BCS",
			mode:        relative,
			noBytes:     2,
			noCycles:    2,
			extraCycles: 1, // XXX or 2
		},
		// 0xb1
		{
			mnemonic:    "LDA",
			mode:        zeroPageIndirectY,
			noBytes:     2,
			noCycles:    5,
			extraCycles: 1,
		},
		// 0xb2
		invalidOpcode,
		// 0xb3
		invalidOpcode,
		// 0xb4
		{
			mnemonic: "LDY",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 4,
		},
		// 0xb5
		{
			mnemonic: "LDA",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 4,
		},
		// 0xb6
		{
			mnemonic: "LDX",
			mode:     zeroPageY,
			noBytes:  2,
			noCycles: 4,
		},
		// 0xb7
		invalidOpcode,
		// 0xb8
		{
			mnemonic: "CLV",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0xb9
		{
			mnemonic:    "LDA",
			mode:        absoluteY,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0xba
		{
			mnemonic: "TSX",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0xbb
		invalidOpcode,
		// 0xbc
		{
			mnemonic:    "LDY",
			mode:        absoluteX,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0xbd
		{
			mnemonic:    "LDA",
			mode:        absoluteX,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0xbe
		{
			mnemonic:    "LDX",
			mode:        absoluteY,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0xbf
		invalidOpcode,
		// 0xc0
		{
			mnemonic: "CPY",
			mode:     immediate,
			noBytes:  2,
			noCycles: 2,
		},
		// 0xc1
		{
			mnemonic: "CMP",
			mode:     zeroPageIndirectX,
			noBytes:  2,
			noCycles: 6,
		},
		// 0xc2
		invalidOpcode,
		// 0xc3
		invalidOpcode,
		// 0xc4
		{
			mnemonic: "CPY",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 3,
		},
		// 0xc5
		{
			mnemonic: "CMP",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 3,
		},
		// 0xc6
		{
			mnemonic: "DEC",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 5,
		},
		// 0xc7
		invalidOpcode,
		// 0xc8
		{
			mnemonic: "INY",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0xc9
		{
			mnemonic: "CMP",
			mode:     immediate,
			noBytes:  2,
			noCycles: 2,
		},
		// 0xca
		{
			mnemonic: "DEX",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0xcb
		invalidOpcode,
		// 0xcc
		{
			mnemonic: "CPY",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0xcd
		{
			mnemonic: "CMP",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0xce
		{
			mnemonic: "DEC",
			mode:     absolute,
			noBytes:  3,
			noCycles: 6,
		},
		// 0xcf
		invalidOpcode,
		// 0xd0
		{
			mnemonic:    "BNE",
			mode:        relative,
			noBytes:     2,
			noCycles:    2,
			extraCycles: 1, // XXX or 2
		},
		// 0xd1
		{
			mnemonic:    "CMP",
			mode:        zeroPageIndirectY,
			noBytes:     2,
			noCycles:    5,
			extraCycles: 1,
		},
		// 0xd2
		invalidOpcode,
		// 0xd3
		invalidOpcode,
		// 0xd4
		invalidOpcode,
		// 0xd5
		{
			mnemonic: "CMP",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 4,
		},
		// 0xd6
		{
			mnemonic: "DEC",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 6,
		},
		// 0xd7
		invalidOpcode,
		// 0xd8
		{
			mnemonic: "CLD",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0xd9
		{
			mnemonic:    "CMP",
			mode:        absoluteY,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0xda
		invalidOpcode,
		// 0xdb
		invalidOpcode,
		// 0xdc
		invalidOpcode,
		// 0xdd
		{
			mnemonic:    "CMP",
			mode:        absoluteX,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0xde
		{
			mnemonic: "DEC",
			mode:     absoluteX,
			noBytes:  3,
			noCycles: 7,
		},
		// 0xdf
		invalidOpcode,
		// 0xe0
		{
			mnemonic: "CPX",
			mode:     immediate,
			noBytes:  2,
			noCycles: 2,
		},
		// 0xe1
		{
			mnemonic: "SBC",
			mode:     zeroPageIndirectX,
			noBytes:  2,
			noCycles: 6,
		},
		// 0xe2
		invalidOpcode,
		// 0xe3
		invalidOpcode,
		// 0xe4
		{
			mnemonic: "CPX",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 3,
		},
		// 0xe5
		{
			mnemonic: "SBC",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 3,
		},
		// 0xe6
		{
			mnemonic: "INC",
			mode:     zeroPage,
			noBytes:  2,
			noCycles: 5,
		},
		// 0xe7
		invalidOpcode,
		// 0xe8
		{
			mnemonic: "INX",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0xe9
		{
			mnemonic: "SBC",
			mode:     immediate,
			noBytes:  2,
			noCycles: 2,
		},
		// 0xea
		{
			mnemonic: "NOP",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0xeb
		invalidOpcode,
		// 0xec
		{
			mnemonic: "CPX",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0xed
		{
			mnemonic: "SBC",
			mode:     absolute,
			noBytes:  3,
			noCycles: 4,
		},
		// 0xee
		{
			mnemonic: "INC",
			mode:     absolute,
			noBytes:  3,
			noCycles: 6,
		},
		// 0xef
		invalidOpcode,
		// 0xf0
		{
			mnemonic:    "BEQ",
			mode:        relative,
			noBytes:     2,
			noCycles:    2,
			extraCycles: 1, // XXX or 2
		},
		// 0xf1
		{
			mnemonic: "SBC",
			mode:     zeroPageIndirectY,
			noBytes:  2,
			noCycles: 5,
		},
		// 0xf2
		invalidOpcode,
		// 0xf3
		invalidOpcode,
		// 0xf4
		invalidOpcode,
		// 0xf5
		{
			mnemonic: "SBC",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 4,
		},
		// 0xf6
		{
			mnemonic: "INC",
			mode:     zeroPageX,
			noBytes:  2,
			noCycles: 6,
		},
		// 0xf7
		invalidOpcode,
		// 0xf8
		{
			mnemonic: "SED",
			mode:     implied,
			noBytes:  1,
			noCycles: 2,
		},
		// 0xf9
		{
			mnemonic:    "SBC",
			mode:        absoluteY,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0xfa
		invalidOpcode,
		// 0xfb
		invalidOpcode,
		// 0xfc
		invalidOpcode,
		// 0xfd
		{
			mnemonic:    "SBC",
			mode:        absoluteX,
			noBytes:     3,
			noCycles:    4,
			extraCycles: 1,
		},
		// 0xfe
		{
			mnemonic: "INC",
			mode:     absoluteX,
			noBytes:  3,
			noCycles: 7,
		},
		// 0xff
		invalidOpcode,
	}
)

type CPU struct {
	pc     uint16 // program counter
	sp     byte   // stack pointer
	sr     byte   // status register
	a      byte   // accumulator
	x      byte   // X index register
	y      byte   // Y index register
	cycles uint64

	// 0000-00FF  - RAM for Zero-Page & Indirect-Memory Addressing
	// 0100-01FF  - RAM for Stack Space & Absolute Addressing
	// 0200-3FFF  - RAM for programmer use
	// 4000-7FFF  - Memory mapped I/O
	// 8000-FFF9  - ROM for programmer useage
	// FFFA       - Vector address for NMI (low byte)
	// FFFB       - Vector address for NMI (high byte)
	// FFFC       - Vector address for RESET (low byte)
	// FFFD       - Vector address for RESET (high byte)
	// FFFE       - Vector address for IRQ & BRK (low byte)
	// FFFF       - Vector address for IRQ & BRK (high byte)
	memory []byte // memory
}

func New() *CPU {
	c := CPU{
		memory: make([]byte, 65536),
		sp:     0xff, // 0x01ff by convention
		sr:     0x34,
	}

	return &c
}

func (c *CPU) evalZ(src byte) {
	if src == 0x00 {
		c.sr |= Zero
	} else {
		c.sr &^= Zero
	}
}

func (c *CPU) evalN(src byte) {
	if src&0x80 == 0x80 {
		c.sr |= Negative
	} else {
		c.sr &^= Negative
	}
}

func (c *CPU) evalV(src byte) {
	if src&0x40 == 0x40 {
		c.sr |= Overflow
	} else {
		c.sr &^= Overflow
	}
}

func (c *CPU) lda(src byte) {
	c.a = src
	c.evalN(c.a)
	c.evalZ(c.a)
}

func (c *CPU) sta(addr uint16) {
	c.memory[addr] = c.a
}

func (c *CPU) ldy(src byte) {
	c.y = src
	c.evalN(c.y)
	c.evalZ(c.y)
}

func (c *CPU) sty(addr uint16) {
	c.memory[addr] = c.y
}

func (c *CPU) ldx(src byte) {
	c.x = src
	c.evalN(c.x)
	c.evalZ(c.x)
}

func (c *CPU) stx(addr uint16) {
	c.memory[addr] = c.x
}

func (c *CPU) txa() {
	c.a = c.x
	c.evalN(c.a)
	c.evalZ(c.a)
}

func (c *CPU) txs() {
	c.sp = c.x
}

func (c *CPU) tsx() {
	c.x = c.sp
	c.evalN(c.x)
	c.evalZ(c.x)
}

func (c *CPU) tya() {
	c.a = c.y
	c.evalN(c.a)
	c.evalZ(c.a)
}

func (c *CPU) tay() {
	c.y = c.a
	c.evalN(c.y)
	c.evalZ(c.y)
}

func (c *CPU) tax() {
	c.x = c.a
	c.evalN(c.x)
	c.evalZ(c.x)
}

func (c *CPU) dey() {
	c.y--
	c.evalN(c.y)
	c.evalZ(c.y)
}

func (c *CPU) dex() {
	c.x--
	c.evalN(c.x)
	c.evalZ(c.x)
}

func (c *CPU) inx() {
	c.x++
	c.evalN(c.x)
	c.evalZ(c.x)
}

func (c *CPU) iny() {
	c.y++
	c.evalN(c.y)
	c.evalZ(c.y)
}

func (c *CPU) inc(src *byte) {
	*src++
	c.evalN(*src)
	c.evalZ(*src)
}

func (c *CPU) dec(src *byte) {
	*src--
	c.evalN(*src)
	c.evalZ(*src)
}

func (c *CPU) sec() {
	c.sr |= Carry
}

func (c *CPU) ora(src byte) {
	c.a |= src
	c.evalN(c.a)
	c.evalZ(c.a)
}

func (c *CPU) and(src byte) {
	c.a &= src
	c.evalN(c.a)
	c.evalZ(c.a)
}

func (c *CPU) eor(src byte) {
	c.a ^= src
	c.evalN(c.a)
	c.evalZ(c.a)
}

func (c *CPU) lsr(src *byte) {
	if *src&0x01 == 0x01 {
		c.sr |= Carry
	} else {
		c.sr &^= Carry
	}
	*src >>= 1
	c.sr &^= Negative // clear N
	c.evalZ(*src)
}

func (c *CPU) cpy(src byte) {
	t := c.y - src
	c.evalN(t)
	c.evalZ(t)
	if c.y >= src {
		c.sr |= Carry
	} else {
		c.sr &^= Carry
	}
}

func (c *CPU) cpx(src byte) {
	t := c.x - src
	c.evalN(t)
	c.evalZ(t)
	if c.x >= src {
		c.sr |= Carry
	} else {
		c.sr &^= Carry
	}
}

func (c *CPU) cmp(src byte) {
	t := c.a - src
	c.evalN(t)
	c.evalZ(t)
	if c.a >= src {
		c.sr |= Carry
	} else {
		c.sr &^= Carry
	}
}

func (c *CPU) adcNormal(i byte, j byte, carryIn byte) {
	result16 := uint16(i) + uint16(j) + uint16(carryIn)
	r := byte(result16)
	carryOut := (result16 & 0x100) != 0
	overflow := (i^r)&(j^r)&0x80 != 0

	// set carry
	if carryOut {
		c.sr |= Carry
	} else {
		c.sr &^= Carry
	}
	// set overflow
	if overflow {
		c.sr |= Overflow
	} else {
		c.sr &^= Overflow
	}
	c.evalN(r)
	c.evalZ(r)

	c.a = r
	//fmt.Printf("adcNormal a %02x, sr %02x\n", c.a, c.sr)
}

func (c *CPU) adcDecimal(i byte, j byte, carryIn byte) {
	var carryB byte = 0

	low := (i & 0x0f) + (j & 0x0f) + carryIn
	if (low & 0xff) > 9 {
		low += 6
	}
	if low > 15 {
		carryB = 1
	}

	high := (i >> 4) + (j >> 4) + carryB
	if (high & 0xff) > 9 {
		high += 6
	}

	r := (low & 0x0f) | (high<<4)&0xf0

	if high > 15 {
		c.sr |= Carry
	} else {
		c.sr &^= Carry
	}
	c.evalZ(r)
	c.sr &^= Negative
	c.sr &^= Overflow

	c.a = r
}

func (c *CPU) sbcDecimal(i byte, j byte, carryIn byte) {
	var carryB byte = 0

	if carryIn == 0 {
		carryIn = 1
	} else {
		carryIn = 0
	}

	low := (i & 0x0f) - (j & 0x0f) - carryIn
	if (low & 0x10) != 0 {
		low -= 6
	}
	if (low & 0x10) != 0 {
		carryB = 1
	}

	high := (i >> 4) - (j >> 4) - carryB
	if (high & 0x10) != 0 {
		high -= 6
	}

	r := (low & 0x0f) | (high << 4)

	if high&0xff < 15 {
		c.sr |= Carry
	} else {
		c.sr &^= Carry
	}
	c.evalZ(r)
	c.sr &^= Negative
	c.sr &^= Overflow

	c.a = r
}

func (c *CPU) adc(src byte) {
	carry := byte(0)
	if c.sr&Carry == Carry {
		carry = 1
	}
	if c.sr&BCD == BCD {
		c.adcDecimal(c.a, src, carry)
	} else {
		c.adcNormal(c.a, src, carry)
	}
}

func (c *CPU) sbc(src byte) {
	carry := byte(0)
	if c.sr&Carry == Carry {
		carry = 1
	}
	if c.sr&BCD == BCD {
		c.sbcDecimal(c.a, src, carry)
	} else {
		c.adcNormal(c.a, ^src, carry)
	}
}

func (c *CPU) bit(src byte) {
	// XXX this needs to be optimized to not have garbage

	r := c.a & src
	c.evalZ(r)
	// There seems to be some confusion about this.
	// Some sources say look at src others say look at r.
	c.evalN(src)
	c.evalV(src)
}

func (c *CPU) rol(src *byte) {
	// XXX this needs to be optimized to not have garbage

	// save carry bit
	carry := c.sr & Carry

	// rol
	r := uint16(*src) << 1
	// set carry
	if r&0x100 == 0x100 {
		c.sr |= Carry
	} else {
		c.sr &^= Carry
	}
	// set bit 0 to saved carry
	if carry != 0 {
		r |= 0x01
	}
	*src = byte(r)
	c.evalN(*src)
	c.evalZ(*src)
}

func (c *CPU) ror(src *byte) {
	// XXX this needs to be optimized to not have garbage

	// save carry bit
	carry := c.sr & Carry

	// set carry
	if *src&0x01 == 0x01 {
		c.sr |= Carry
	} else {
		c.sr &^= Carry
	}

	// ror
	r := *src >> 1

	// set bit 7 to saved carry
	if carry != 0 {
		r |= 0x80
	}
	*src = byte(r)
	c.evalN(*src)
	c.evalZ(*src)
}

func (c *CPU) pha() {
	c.memory[0x0100+uint16(c.sp)] = c.a
	c.sp--
}

func (c *CPU) pla() {
	c.sp++
	c.a = c.memory[0x0100+uint16(c.sp)]
	c.evalN(c.a)
	c.evalZ(c.a)
}

func (c *CPU) php() {
	c.memory[0x0100+uint16(c.sp)] = c.sr | Unused | Break
	c.sp--
}

func (c *CPU) plp() {
	c.sp++
	c.sr = c.memory[0x0100+uint16(c.sp)] | Unused
}

func (c *CPU) sed() {
	c.sr |= BCD
}

func (c *CPU) sei() {
	c.sr |= Interrupts
}

func (c *CPU) cld() {
	c.sr &^= BCD
}

func (c *CPU) clc() {
	c.sr &^= Carry
}

func (c *CPU) cli() {
	c.sr &^= Interrupts
}

func (c *CPU) clv() {
	c.sr &^= Overflow
}

func (c *CPU) asl(src *byte) {
	// carry = bit 7
	if *src&0x80 == 0x80 {
		c.sr |= Carry
	} else {
		c.sr &^= Carry
	}
	*src = *src << 1
	c.evalN(*src)
	c.evalZ(*src)
}

func (c *CPU) bpl(src byte) bool {
	if c.sr&Negative == 0x00 {
		c.pc = c.relative(c.pc)
		return true
	}
	return false
}

func (c *CPU) bcc(src byte) bool {
	if c.sr&Carry == 0x00 {
		c.pc = c.relative(c.pc)
		return true
	}
	return false
}

func (c *CPU) beq(src byte) bool {
	if c.sr&Zero == Zero {
		c.pc = c.relative(c.pc)
		return true
	}
	return false
}

func (c *CPU) bne(src byte) bool {
	if c.sr&Zero == 0x00 {
		c.pc = c.relative(c.pc)
		return true
	}
	return false
}

func (c *CPU) bcs(src byte) bool {
	if c.sr&Carry == Carry {
		c.pc = c.relative(c.pc)
		return true
	}
	return false
}

func (c *CPU) bvc(src byte) bool {
	if c.sr&Overflow == 0x00 {
		c.pc = c.relative(c.pc)
		return true
	}
	return false
}

func (c *CPU) bvs(src byte) bool {
	if c.sr&Overflow == Overflow {
		c.pc = c.relative(c.pc)
		return true
	}
	return false
}

func (c *CPU) bmi(src byte) bool {
	if c.sr&Negative == Negative {
		c.pc = c.relative(c.pc)
		return true
	}
	return false
}

func (c *CPU) jsr(addr uint16) {
	c.pc += uint16(opcodes[0x20].noBytes) - 1
	c.memory[0x0100+uint16(c.sp)] = byte(c.pc >> 8)
	c.sp--
	c.memory[0x0100+uint16(c.sp)] = byte(c.pc)
	c.sp--
	c.pc = addr
}

func (c *CPU) jmp(addr uint16) {
	c.pc = addr
}

func (c *CPU) brk() {
	c.sr |= Break

	// note that brk has a quirk that it skips 1 byte past pc
	pc := c.pc + uint16(opcodes[0x00].noBytes) + 1
	// high byte
	c.memory[0x0100+uint16(c.sp)] = byte(pc >> 8)
	c.sp--

	// low byte
	c.memory[0x0100+uint16(c.sp)] = byte(pc)
	c.sp--

	// status register
	c.memory[0x0100+uint16(c.sp)] = c.sr
	c.sp--

	c.sr |= Interrupts

	// set pc to interrupt vector
	c.pc = uint16(c.memory[0xffff])<<8 | uint16(c.memory[0xfffe])
}

func (c *CPU) rti() {
	// XXX this needs to be optimized to not have garbage

	// status register
	c.sp++
	c.sr = c.memory[0x0100+uint16(c.sp)] | Unused

	// low byte
	c.sp++
	l := c.memory[0x0100+uint16(c.sp)]

	// high byte
	c.sp++
	h := c.memory[0x0100+uint16(c.sp)]

	c.pc = uint16(l) | uint16(h)<<8
}

func (c *CPU) rts() {
	// XXX this needs to be optimized to not have garbage

	// low byte
	c.sp++
	l := c.memory[0x0100+uint16(c.sp)]

	// high byte
	c.sp++
	h := c.memory[0x0100+uint16(c.sp)]

	c.pc = uint16(l) | uint16(h)<<8 + 1
}

func (c *CPU) relative(addr uint16) uint16 {
	rel := int8(c.memory[addr+1])
	// Note that we post increment PC so we have to account for that here.
	// This may have to change in order to emulate hardware more correctly.
	addr += 2
	if rel >= 0 {
		addr += uint16(rel)
	} else {
		addr -= -uint16(rel)
	}
	return addr
}

func (c *CPU) immediate(address uint16) uint16 {
	return address + 1
}

// indirect returns (addr+1 | addr+2<<8)
func (c *CPU) indirect(addr uint16) uint16 {
	a := uint16(c.memory[addr+1]) | uint16(c.memory[addr+2])<<8
	return uint16(c.memory[a]) | uint16(c.memory[a+1])<<8
}

// zeroPage returns zp addr + offs
func (c *CPU) zeroPage(addr uint16, offs byte) uint16 {
	return uint16(c.memory[addr+1] + offs)
}

// absoluteX returns addr+1 | addr+2<<8 + ofs
func (c *CPU) absolute(addr uint16, offs byte) uint16 {
	return uint16(c.memory[addr+1]) | uint16(c.memory[addr+2])<<8 +
		uint16(offs)
}

// indexedIndirectX returns (zp,x)
func (c *CPU) indexedIndirectX(addr uint16, offs byte) uint16 {
	zpa := (uint16(c.memory[addr+1]) + uint16(offs)) & uint16(0xff)
	return uint16(c.memory[zpa]) | uint16(c.memory[zpa+1])<<8
}

// indexedIndirectY returns (zp),y
func (c *CPU) indexedIndirectY(addr uint16, offs byte) uint16 {
	return uint16(c.memory[c.memory[addr+1]]) |
		uint16(c.memory[c.memory[addr+1]+1])<<8 +
		uint16(offs)
}

func (c *CPU) executeInstruction() {
	// decode instruction
	opcode := c.memory[c.pc]
	c.cycles += opcodes[opcode].noCycles + opcodes[opcode].extraCycles
	switch opcode {
	case 0x00:
		c.brk()
		return
	case 0x01:
		c.ora(c.memory[c.indexedIndirectX(c.pc, c.x)])
	case 0x05:
		c.ora(c.memory[c.zeroPage(c.pc, 0)])
	case 0x06:
		c.asl(&c.memory[c.zeroPage(c.pc, 0)])
	case 0x08:
		c.php()
	case 0x09:
		c.ora(c.memory[c.immediate(c.pc)])
	case 0x0a:
		c.asl(&c.a)
	case 0x0d:
		c.ora(c.memory[c.absolute(c.pc, 0)])
	case 0x0e:
		c.asl(&c.memory[c.absolute(c.pc, 0)])
	case 0x10:
		if c.bpl(c.memory[c.pc+1]) {
			return
		}
	case 0x11:
		c.ora(c.memory[c.indexedIndirectY(c.pc, c.y)])
	case 0x15:
		c.ora(c.memory[c.zeroPage(c.pc, c.x)])
	case 0x16:
		c.asl(&c.memory[c.zeroPage(c.pc, c.x)])
	case 0x18:
		c.clc()
	case 0x19:
		c.ora(c.memory[c.absolute(c.pc, c.y)])
	case 0x1d:
		c.ora(c.memory[c.absolute(c.pc, c.x)])
	case 0x1e:
		c.asl(&c.memory[c.absolute(c.pc, c.x)])
	case 0x20:
		c.jsr(c.absolute(c.pc, 0))
		return
	case 0x21:
		c.and(c.memory[c.indexedIndirectX(c.pc, c.x)])
	case 0x24:
		c.bit(c.memory[c.zeroPage(c.pc, 0)])
	case 0x25:
		c.and(c.memory[c.zeroPage(c.pc, 0)])
	case 0x26:
		c.rol(&c.memory[c.zeroPage(c.pc, 0)])
	case 0x28:
		c.plp()
	case 0x29:
		c.and(c.memory[c.immediate(c.pc)])
	case 0x2a:
		c.rol(&c.a)
	case 0x2c:
		c.bit(c.memory[c.absolute(c.pc, 0)])
	case 0x2d:
		c.and(c.memory[c.absolute(c.pc, 0)])
	case 0x2e:
		c.rol(&c.memory[c.absolute(c.pc, 0)])
	case 0x30:
		if c.bmi(c.memory[c.pc+1]) {
			return
		}
	case 0x31:
		c.and(c.memory[c.indexedIndirectY(c.pc, c.y)])
	case 0x35:
		c.and(c.memory[c.zeroPage(c.pc, c.x)])
	case 0x36:
		c.rol(&c.memory[c.zeroPage(c.pc, c.x)])
	case 0x38:
		c.sec()
	case 0x39:
		c.and(c.memory[c.absolute(c.pc, c.y)])
	case 0x3d:
		c.and(c.memory[c.absolute(c.pc, c.x)])
	case 0x3e:
		c.rol(&c.memory[c.absolute(c.pc, c.x)])
	case 0x40:
		c.rti()
		return
	case 0x41:
		c.eor(c.memory[c.indexedIndirectX(c.pc, c.x)])
	case 0x45:
		c.eor(c.memory[c.zeroPage(c.pc, 0)])
	case 0x46:
		c.lsr(&c.memory[c.zeroPage(c.pc, 0)])
	case 0x48:
		c.pha()
	case 0x49:
		c.eor(c.memory[c.immediate(c.pc)])
	case 0x4a:
		c.lsr(&c.a)
	case 0x4c:
		c.jmp(c.absolute(c.pc, 0))
		return
	case 0x4d:
		c.eor(c.memory[c.absolute(c.pc, 0)])
	case 0x4e:
		c.lsr(&c.memory[c.absolute(c.pc, 0)])
	case 0x50:
		if c.bvc(c.memory[c.pc+1]) {
			return
		}
	case 0x51:
		c.eor(c.memory[c.indexedIndirectY(c.pc, c.y)])
	case 0x55:
		c.eor(c.memory[c.zeroPage(c.pc, c.x)])
	case 0x56:
		c.lsr(&c.memory[c.zeroPage(c.pc, c.x)])
	case 0x58:
		c.cli()
	case 0x59:
		c.eor(c.memory[c.absolute(c.pc, c.y)])
	case 0x5d:
		c.eor(c.memory[c.absolute(c.pc, c.x)])
	case 0x5e:
		c.lsr(&c.memory[c.absolute(c.pc, c.x)])
	case 0x60:
		c.rts()
		return
	case 0x61:
		c.adc(c.memory[c.indexedIndirectX(c.pc, c.x)])
	case 0x65:
		c.adc(c.memory[c.zeroPage(c.pc, 0)])
	case 0x66:
		c.ror(&c.memory[c.zeroPage(c.pc, 0)])
	case 0x68:
		c.pla()
	case 0x69:
		c.adc(c.memory[c.immediate(c.pc)])
	case 0x6a:
		c.ror(&c.a)
	case 0x6c:
		c.jmp(c.indirect(c.pc))
		return
	case 0x6d:
		c.adc(c.memory[c.absolute(c.pc, 0)])
	case 0x6e:
		c.ror(&c.memory[c.absolute(c.pc, 0)])
	case 0x70:
		if c.bvs(c.memory[c.pc+1]) {
			return
		}
	case 0x71:
		c.adc(c.memory[c.indexedIndirectY(c.pc, c.y)])
	case 0x75:
		c.adc(c.memory[c.zeroPage(c.pc, c.x)])
	case 0x76:
		c.ror(&c.memory[c.zeroPage(c.pc, c.x)])
	case 0x78:
		c.sei()
	case 0x79:
		c.adc(c.memory[c.absolute(c.pc, c.y)])
	case 0x7d:
		c.adc(c.memory[c.absolute(c.pc, c.x)])
	case 0x7e:
		c.ror(&c.memory[c.absolute(c.pc, c.x)])
	case 0x81:
		c.sta(c.indexedIndirectX(c.pc, c.x))
	case 0x84:
		c.sty(c.zeroPage(c.pc, 0))
	case 0x85:
		c.sta(c.zeroPage(c.pc, 0))
	case 0x86:
		c.stx(c.zeroPage(c.pc, 0))
	case 0x88:
		c.dey()
	case 0x8a:
		c.txa()
	case 0x8c:
		c.sty(c.absolute(c.pc, 0))
	case 0x8d:
		c.sta(c.absolute(c.pc, 0))
	case 0x8e:
		c.stx(c.absolute(c.pc, 0))
	case 0x90:
		if c.bcc(c.memory[c.pc+1]) {
			return
		}
	case 0x91:
		c.sta(c.indexedIndirectY(c.pc, c.y))
	case 0x94:
		c.sty(c.zeroPage(c.pc, c.x))
	case 0x95:
		c.sta(c.zeroPage(c.pc, c.x))
	case 0x96:
		c.stx(c.zeroPage(c.pc, c.y))
	case 0x98:
		c.tya()
	case 0x99:
		c.sta(c.absolute(c.pc, c.y))
	case 0x9a:
		c.txs()
	case 0x9d:
		c.sta(c.absolute(c.pc, c.x))
	case 0xa0:
		c.ldy(c.memory[c.immediate(c.pc)])
	case 0xa1:
		c.lda(c.memory[c.indexedIndirectX(c.pc, c.x)])
	case 0xa2:
		c.ldx(c.memory[c.immediate(c.pc)])
	case 0xa4:
		c.ldy(c.memory[c.zeroPage(c.pc, 0)])
	case 0xa5:
		c.lda(c.memory[c.zeroPage(c.pc, 0)])
	case 0xa6:
		c.ldx(c.memory[c.zeroPage(c.pc, 0)])
	case 0xa8:
		c.tay()
	case 0xa9:
		c.lda(c.memory[c.immediate(c.pc)])
	case 0xaa:
		c.tax()
	case 0xac:
		c.ldy(c.memory[c.absolute(c.pc, 0)])
	case 0xad:
		c.lda(c.memory[c.absolute(c.pc, 0)])
	case 0xae:
		c.ldx(c.memory[c.absolute(c.pc, 0)])
	case 0xb0:
		if c.bcs(c.memory[c.pc+1]) {
			return
		}
	case 0xb1:
		c.lda(c.memory[c.indexedIndirectY(c.pc, c.y)])
	case 0xb4:
		c.ldy(c.memory[c.zeroPage(c.pc, c.x)])
	case 0xb5:
		c.lda(c.memory[c.zeroPage(c.pc, c.x)])
	case 0xb6:
		c.ldx(c.memory[c.zeroPage(c.pc, c.y)])
	case 0xb8:
		c.clv()
	case 0xb9:
		c.lda(c.memory[c.absolute(c.pc, c.y)])
	case 0xba:
		c.tsx()
	case 0xbc:
		c.ldy(c.memory[c.absolute(c.pc, c.x)])
	case 0xbd:
		c.lda(c.memory[c.absolute(c.pc, c.x)])
	case 0xbe:
		c.ldx(c.memory[c.absolute(c.pc, c.y)])
	case 0xc0:
		c.cpy(c.memory[c.immediate(c.pc)])
	case 0xc1:
		c.cmp(c.memory[c.indexedIndirectX(c.pc, c.x)])
	case 0xc4:
		c.cpy(c.memory[c.zeroPage(c.pc, 0)])
	case 0xc5:
		c.cmp(c.memory[c.zeroPage(c.pc, 0)])
	case 0xc6:
		c.dec(&c.memory[c.zeroPage(c.pc, 0)])
	case 0xc8:
		c.iny()
	case 0xc9:
		c.cmp(c.memory[c.immediate(c.pc)])
	case 0xca:
		c.dex()
	case 0xcc:
		c.cpy(c.memory[c.absolute(c.pc, 0)])
	case 0xcd:
		c.cmp(c.memory[c.absolute(c.pc, 0)])
	case 0xce:
		c.dec(&c.memory[c.absolute(c.pc, 0)])
	case 0xd0:
		if c.bne(c.memory[c.pc+1]) {
			return
		}
	case 0xd1:
		c.cmp(c.memory[c.indexedIndirectY(c.pc, c.y)])
	case 0xd5:
		c.cmp(c.memory[c.zeroPage(c.pc, c.x)])
	case 0xd6:
		c.dec(&c.memory[c.zeroPage(c.pc, c.x)])
	case 0xd8:
		c.cld()
	case 0xd9:
		c.cmp(c.memory[c.absolute(c.pc, c.y)])
	case 0xdd:
		c.cmp(c.memory[c.absolute(c.pc, c.x)])
	case 0xde:
		c.dec(&c.memory[c.absolute(c.pc, c.x)])
	case 0xe0:
		c.cpx(c.memory[c.immediate(c.pc)])
	case 0xe1:
		c.sbc(c.memory[c.indexedIndirectX(c.pc, c.x)])
	case 0xe4:
		c.cpx(c.memory[c.zeroPage(c.pc, 0)])
	case 0xe5:
		c.sbc(c.memory[c.zeroPage(c.pc, 0)])
	case 0xe6:
		c.inc(&c.memory[c.zeroPage(c.pc, 0)])
	case 0xe8:
		c.inx()
	case 0xe9:
		c.sbc(c.memory[c.immediate(c.pc)])
	case 0xea:
		// nop
	case 0xec:
		c.cpx(c.memory[c.absolute(c.pc, 0)])
	case 0xed:
		c.sbc(c.memory[c.absolute(c.pc, 0)])
	case 0xee:
		c.inc(&c.memory[c.absolute(c.pc, 0)])
	case 0xf0:
		if c.beq(c.memory[c.pc+1]) {
			return
		}
	case 0xf1:
		c.sbc(c.memory[c.indexedIndirectY(c.pc, c.y)])
	case 0xf5:
		c.sbc(c.memory[c.zeroPage(c.pc, c.x)])
	case 0xf6:
		c.inc(&c.memory[c.zeroPage(c.pc, c.x)])
	case 0xf8:
		c.sed()
	case 0xf9:
		c.sbc(c.memory[c.absolute(c.pc, c.y)])
	case 0xfd:
		c.sbc(c.memory[c.absolute(c.pc, c.x)])
	case 0xfe:
		c.inc(&c.memory[c.absolute(c.pc, c.x)])
	default:
		// make this less drastic
		panic(fmt.Sprintf("invalid opcode: $%02x PC $%04x",
			opcode, c.pc))
	}

	c.pc += uint16(opcodes[opcode].noBytes)
}

func (c *CPU) snapshot() string {
	return fmt.Sprintf("A: $%02x X: $%02x Y: $%02x SR: $%02x PC: $%04x SP: $%02x",
		c.a,
		c.x,
		c.y,
		c.sr,
		c.pc,
		c.sp)
}

// disassemble disassembles an instruction at address and returns the
// instruction and bytes consumed.
func (c *CPU) disassemble(address uint16) (string, byte) {
	o := opcodes[c.memory[address]]
	switch o.mode {
	case accumulator:
		return fmt.Sprintf("%v", o.mnemonic), o.noBytes
	case implied:
		return fmt.Sprintf("%v", o.mnemonic), o.noBytes
	case immediate:
		return fmt.Sprintf("%v\t#$%02X", o.mnemonic,
			c.memory[c.immediate(address)]), o.noBytes
	case indirect:
		// absolute call here is intended
		return fmt.Sprintf("%v\t($%02X)", o.mnemonic,
			c.absolute(address, 0)), o.noBytes
	case relative:
		return fmt.Sprintf("%v\t$%04X", o.mnemonic,
			c.relative(address)), o.noBytes
	case zeroPage:
		return fmt.Sprintf("%v\t$%02X", o.mnemonic,
			c.zeroPage(address, 0)), o.noBytes
	case zeroPageX:
		return fmt.Sprintf("%v\t$%02X,X", o.mnemonic,
			c.zeroPage(address, 0)), o.noBytes
	case zeroPageY:
		return fmt.Sprintf("%v\t$%02X,Y", o.mnemonic,
			c.zeroPage(address, 0)), o.noBytes
	case absolute:
		return fmt.Sprintf("%v\t$%04X", o.mnemonic,
			c.absolute(address, 0)), o.noBytes
	case absoluteX:
		return fmt.Sprintf("%v\t$%04X,X", o.mnemonic,
			c.absolute(address, 0)), o.noBytes
	case absoluteY:
		return fmt.Sprintf("%v\t$%04X,Y", o.mnemonic,
			c.absolute(address, 0)), o.noBytes
	case zeroPageIndirectX:
		// zeroPage call here is intended
		return fmt.Sprintf("%v\t($%02X,X)", o.mnemonic,
			c.zeroPage(address, 0)), o.noBytes
	case zeroPageIndirectY:
		// zeroPage call here is intended
		return fmt.Sprintf("%v\t($%02X),Y", o.mnemonic,
			c.zeroPage(address, 0)), o.noBytes
	default:
		return fmt.Sprintf("INVALID MODE"), 0
	}
}
