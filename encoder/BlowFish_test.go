package encoder

import (
	"fmt"
	"testing"
)

func TestBlowFish(t *testing.T) {
	encoder, _ := NewBlowFishEncoder("")
	text := "acerwhdhdhdhdhdhdhgagdgdgagei hahah 我的"
	et := encoder.Encode([]byte(text))
	fmt.Println(et)
	pt := encoder.Decode(et)
	fmt.Println(string(pt))
}
