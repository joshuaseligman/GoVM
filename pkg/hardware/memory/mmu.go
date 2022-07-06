package memory

import (
	"github.com/joshuaseligman/GoVM/pkg/hardware"
)

// The struct for the memory management unit
type Mmu struct {
	hw *hardware.Hardware // The hardware struct
	mar uint // The memory address register
	mdr uint32 // The memory data register
	memory *Memory
}

// Function that creates the MMU
func NewMmu() *Mmu {
	mmu := Mmu {
		hw: hardware.NewHardware("MMU", 0),
		mar: 0,
		mdr: 0,
		memory: NewMemory(0x10000),
	}
	return &mmu
}

// Sends the signal to memory to read the value in the address of the MAR
func (mmu *Mmu) CallRead() {
	mmu.memory.Read()
}

// Sends the signal to memory to write the value in the MDR to the address of the MAR
func (mmu *Mmu) CallWrite() {
	mmu.memory.Write()
}

// Sets the MAR of the MMU
func (mmu *Mmu) SetMar(newMar uint) {
	mmu.mar = newMar
}

// Sets the MDR of the MMU
func (mmu *Mmu) SetMdr(newMdr uint32) {
	mmu.mdr = newMdr
}

// Logs a message
func (mmu *Mmu) Log(msg string) {
	mmu.hw.Log(msg)
}