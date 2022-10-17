package ceres

import (
)

type PCS struct {
    pronounDictionary map[Word]Entity

    initialized bool
}

func (pcs *PCS)Initialize() {
    pcs.pronounDictionary = make(map[Word]Entity)

    pcs.initialized = true
}
