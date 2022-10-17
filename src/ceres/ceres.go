package ceres


type CERES struct {
    ics ICS
    pcs PCS

    root *EntityType

    initialized bool
}

func (c *CERES)Initialize(){
    c.ics = ICS{}
    c.ics.Initialize()
    c.pcs = PCS{}
    c.pcs.Initialize()

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
