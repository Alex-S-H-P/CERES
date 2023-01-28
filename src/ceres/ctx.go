package ceres

import (
	"CERES/src/utils"
)

type SOMEONE = *EntityInstance // in particular, SOMEONE is either human or self

type PreExpressedCTXEntry struct {
	e Entity
	s float64
}

func (pece PreExpressedCTXEntry) Entity() Entity {
	return pece.e
}

func (pece PreExpressedCTXEntry) Equal(other utils.Equalable) bool {
	if pece2, ok := other.(PreExpressedCTXEntry); ok {
		if pece.e.Equal(pece2.e) && pece.s != pece2.s {
			return true
		}
	}

	return false
}

/*
The conversationnal context.
*/
type CTX struct {
	SPEAKER    SOMEONE
	DESTINATOR SOMEONE

	expressed_buffer *utils.Buffer[PreExpressedCTXEntry]
}

func (ctx *CTX) Initialize() {
	// TODO: set SPEAKER & DESTINATOR to interlocutor & self
	ctx.expressed_buffer = utils.NewBuffer[PreExpressedCTXEntry](256)
}

func (ctx *CTX) P_ctx() float64 {
	var P float64 = 0.

	for i := 0; i < ctx.expressed_buffer.Len(); i++ {
		entry := ctx.expressed_buffer.Get(i)
		P *= entry.s
	}

	return P
}
