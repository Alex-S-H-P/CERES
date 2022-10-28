package ceres


type CERES struct {
    ics ICS
    pcs PCS
    ctx *CTX

    root *EntityType

    initialized bool

    sentence_analyser_workers int

}

func (c *CERES)Initialize(workers int){
    c.ics = ICS{}
    c.ics.Initialize(c)
    c.pcs = PCS{}
    c.pcs.Initialize()
    c.ctx = new(CTX)
    c.ctx.Initialize()

    c.sentence_analyser_workers = 1

    c.createEntityType("entity")

    c.initialized = true
}

/*
Creates a valid entityType.

Adds it to the CERES dictionary.

*/
func (c *CERES)  createEntityType(w Word) {
    et := c.ics.createEntityType(w)

    if c.root == nil {
        c.root = et
    } else {
        c.root.addChild(et)
    }
}


func (c *CERES) createEntityInstance(w Word, et *EntityType) {
    c.ics.createEntityInstance(w, et)
}


func (c* CERES) computeP(idp *dijkstraPossibility) float64 {
    // TODO: Do this
    return 0.
}
