// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

//go:build amd64 && !noasm
// +build amd64,!noasm

package p434

import (
	"reflect"
	"testing"
	"testing/quick"

	"github.com/JI-0/circl/dh/sidh/internal/common"
	"golang.org/x/sys/cpu"
)

type OptimFlag uint

const (
	// Indicates that optimisation which uses MUL instruction should be used
	kUse_MUL OptimFlag = 1 << 0
	// Indicates that optimisation which uses MULX, ADOX and ADCX instructions should be used
	kUse_MULXandADxX = 1 << 1
)

func resetCpuFeatures() {
	HasADXandBMI2 = cpu.X86.HasBMI2 && cpu.X86.HasADX
}

// Utility function used for testing Mul implementations. Tests caller provided
// mulFunc against mul()
func testMul(t *testing.T, f1, f2 OptimFlag) {
	doMulTest := func(multiplier, multiplicant common.Fp) bool {
		defer resetCpuFeatures()
		var resMulRef, resMulOptim common.FpX2

		// Compute multiplier*multiplicant with first implementation
		HasADXandBMI2 = (kUse_MULXandADxX & f1) == kUse_MULXandADxX
		mulP434(&resMulOptim, &multiplier, &multiplicant)

		// Compute multiplier*multiplicant with second implementation
		HasADXandBMI2 = (kUse_MULXandADxX & f2) == kUse_MULXandADxX
		mulP434(&resMulRef, &multiplier, &multiplicant)

		// Compare results
		return reflect.DeepEqual(resMulRef, resMulOptim)
	}

	if err := quick.Check(doMulTest, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

// Utility function used for testing REDC implementations. Tests caller provided
// redcFunc against redc()
func testRedc(t *testing.T, f1, f2 OptimFlag) {
	doRedcTest := func(aRR common.FpX2) bool {
		defer resetCpuFeatures()
		var resRedcF1, resRedcF2 common.Fp
		aRRcpy := aRR

		// Compute redc with first implementation
		HasADXandBMI2 = (kUse_MULXandADxX & f1) == kUse_MULXandADxX
		rdcP434(&resRedcF1, &aRR)

		// Compute redc with second implementation
		HasADXandBMI2 = (kUse_MULXandADxX & f2) == kUse_MULXandADxX
		rdcP434(&resRedcF2, &aRRcpy)

		// Compare results
		return reflect.DeepEqual(resRedcF2, resRedcF1)
	}

	if err := quick.Check(doRedcTest, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

// Ensures correctness of implementation of mul operation which uses MULX and ADOX/ADCX
func TestMulWithMULXADxX(t *testing.T) {
	defer resetCpuFeatures()
	if !HasADXandBMI2 {
		t.Skip("MULX, ADCX and ADOX not supported by the platform")
	}
	testMul(t, kUse_MULXandADxX, kUse_MUL)
}

// Ensures correctness of Montgomery reduction implementation which uses MULX
// and ADCX/ADOX.
func TestRedcWithMULXADxX(t *testing.T) {
	defer resetCpuFeatures()
	if !HasADXandBMI2 {
		t.Skip("MULX, ADCX and ADOX not supported by the platform")
	}
	testRedc(t, kUse_MULXandADxX, kUse_MUL)
}
