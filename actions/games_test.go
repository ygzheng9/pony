package actions

import (
	"testing"
)

func Test_gamesCalcMaxEig(t *testing.T) {
	input := `[["1","3","6"],["1/3","1","9"],["1/6","1/9","1"]]`

	gamesCalcMaxEig(input)
}
