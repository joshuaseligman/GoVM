package cpu

import (
	"github.com/joshuaseligman/GoVM/pkg/hardware"
	"github.com/joshuaseligman/GoVM/pkg/hardware/memory"
)

type FetchUnit struct {
	hw *hardware.Hardware
	mmu *memory.Mmu
}

func NewFetchUnit(mem *memory.Memory) *FetchUnit {
	fetchUnit := FetchUnit {
		hw: hardware.NewHardware("IFU", 0),
		mmu: memory.NewMmu(mem),
	}
	return &fetchUnit
}

// Fetches the instruction from memory
func (ifu *FetchUnit) FetchInstruction(addr *uint) *IFIDReg {
	ifu.mmu.SetMar(*addr)
	ifu.mmu.CallRead()
	*addr++
	return NewIFIDReg(ifu.mmu.GetMdr(), *addr)
}

// Logs a message
func (ifu *FetchUnit) Log(msg string) {
	ifu.hw.Log(msg)
}