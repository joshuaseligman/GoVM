package assembler

import (
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
	"fmt"
)

// Assembles a program into instructions for the computer to read
func AssembleProgram(filePath string, maxSize int) []uint32 {
	// Open the file
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Create the array that will become the memory
	program := make([]uint32, maxSize)
	instrIndex := 0

	// Read the file line by line
	for scanner.Scan() {
		instr := scanner.Text()

		// Get the end of the opcode
		opcodeSplit := strings.Index(instr, " ")
		if opcodeSplit == -1 {
			log.Fatal("Invalid instruction ", instr)
		}

		// Get the opcode
		opcode := instr[:opcodeSplit]
		fmt.Println(opcode)

		// Get a list of operands
		operands := strings.Split(instr[opcodeSplit + 1:], ", ")
		fmt.Println(operands)

		instrBin := uint32(0)

		switch opcode {
		// IM instructions
		case "MOVZ", "MOVK":
			instrBin = instrIM(opcode, operands, filePath, instrIndex + 1)
		}

		// Add the instruction to the program
		program[instrIndex] = instrBin
		instrIndex++
	}
	
	return program
}

// Generates the binary for IM instructions
func instrIM(opcode string, operands []string, fileName string, lineNumber int) uint32 {
	// Make sure we have the right number of operands
	if len(operands) != 3 {
		errMsg := fmt.Sprintf("Invalid instruction format: Expected 3 operands but got %d; File: %s; Line: %d", len(operands), fileName, lineNumber)
		log.Fatal(errMsg)
	}

	outBin := uint32(0)

	// Generate initial binary
	switch opcode {
	case "MOVZ":
		outBin = 0b110100101
	case "MOVK":
		outBin = 0b111100101
	}

	// Get the shift amount
	shiftStr := strings.Split(operands[2], " ")[1]
	shiftInt, errConv := strconv.ParseInt(shiftStr, 10, 0)
	if errConv == nil {
		// Add the shift to the binary
		outBin  = outBin << 2 | uint32(shiftInt / 16)
	} else {
		// Bad shift value error
		errMsg := fmt.Sprintf("Bad shift value; File: %s; Line: %d", fileName, lineNumber)
		log.Fatal(errMsg)
	}

	// Get the value to move into the register
	base := 0
	cut := 0
	if len(operands[1]) >= 4 && operands[1][:3] == "#0x" {
		// Base 16
		base = 16
		cut = 3
	} else {
		// Base 10
		base = 10
		cut = 1
	}
	// Get the value based on the base that was decided earlier
	val, errConv := strconv.ParseUint(operands[1][cut:], base, 16) // FIXME Figure out why it wont take a 4 digit hex number and why the value may change
	if errConv == nil {
		// Add the value to the binary
		outBin = outBin << 16 | uint32(val)
	} else {
		errMsg := ""
		if strings.Contains(errConv.Error(), "value out of range") {
			// Out of range errors
			if base == 10 {
				errMsg = fmt.Sprintf("Bad move immediate value: Value must be between 0 and 65,535 (16 bits) but got %s; File: %s; Line: %d", operands[1][1:], fileName, lineNumber)
			} else {
				errMsg = fmt.Sprintf("Bad move immediate value: Value must be between 0x0000 and 0xFFFF (16 bits) but got %s; File: %s; Line: %d", operands[1][1:], fileName, lineNumber)
			}
		} else {
			// Bad value error
			errMsg = fmt.Sprintf("Bad move immediate value; File: %s; Line: %d", fileName, lineNumber)
		}
		log.Fatal(errMsg)
	}

	// Get the register to move the value to
	reg, errConv := strconv.ParseInt(operands[0][1:], 10, 0)
	if errConv != nil {
		// Bad value error
		errMsg := fmt.Sprintf("Bad register value; File: %s; Line: %d", fileName, lineNumber)
		log.Fatal(errMsg)
	} else if reg < 0 || reg > 30 {
		// Invalid register error
		errMsg := fmt.Sprintf("Bad register value: Register must be between 0 and 30 (inclusive); File: %s; Line: %d", fileName, lineNumber)
		log.Fatal(errMsg)
	} else {
		// Add the register to the binary
		outBin = outBin << 5 | uint32(reg)
	}

	// Return the instruction binary
	return outBin
}