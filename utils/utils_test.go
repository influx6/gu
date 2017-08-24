package utils_test

import (
	"testing"

	"github.com/gu-io/gu/utils"
	"github.com/influx6/faux/tests"
)

func TestGenerateEMScale(t *testing.T) {
	scales := utils.GenerateScale(utils.AugmentedFourth, 2, 2)

	if len(scales) == 0 {
		tests.Failed("Should have generated a lists of scales")
	}
	tests.Passed("Should have generated a lists of scales")

	if !floatEquals(scales[0], 0.5002) {
		tests.Failed("Should have generated a lists with scale at index %d to be %.4f", 0, scales[0])
	}
	tests.Passed("Should have generated a lists with scale at index %d to be %.4f", 0, scales[0])

	if floatEquals(scales[2], 1.000) {
		tests.Failed("Should have generated a lists with scale at index %d to be %.4f", 2, scales[2])
	}
	tests.Passed("Should have generated a lists with scale at index %d to be %.4f", 2, scales[2])

	if floatEquals(scales[3], 1.4140) {
		tests.Failed("Should have generated a lists with scale at index %d to be %.4f", 3, scales[3])
	}
	tests.Passed("Should have generated a lists with scale at index %d to be %.4f", 3, scales[3])
}

// EPSILON in mathematics (particularly calculus), an arbitrarily small positive quantity is commonly
// denoted ε.
var EPSILON float64 = 0.00000001

func floatEquals(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}
