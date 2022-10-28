package ceres

import (
    "CERES/src/utils"
)

type SOMEONE = *EntityInstance // in particular, SOMEONE is either human or self

/*
The conversationnal context.
*/
type CTX struct {
    SPEAKER    SOMEONE
    DESTINATOR SOMEONE

    expressed_buffer *utils.Buffer[Entity]
}

func (ctx *CTX) Initialize() {
    // TODO: set SPEAKER & DESTINATOR to interlocutor & self
    ctx.expressed_buffer = utils.NewBuffer[Entity](256)
}
