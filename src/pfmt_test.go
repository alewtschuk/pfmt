package pfmt

import (
	"testing"
)

func TestInitColorMap(t *testing.T) {
	var keys []int = pfmt.initColorMap()
	if len(keys) != 463 {
		t.Errorf("Error: colorMap not correctly initialized. Need key slice of length %d, have length %d.", 463, len(keys))
	}
}

func TestShowAvailableColors(t *testing.T) {
	pfmt.AvailableColors()
}
