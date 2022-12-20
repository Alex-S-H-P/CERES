package ceres

import (
    "fmt"
)
type CERES struct {
    ics ICS
    pcs PCS
    ucs UCS
    ctx *CTX

    root *EntityType

    grammar *grammar

    initialized bool

    sentence_analyser_workers int

}

func (c *CERES)Initialize(workers int){
    c.ics = ICS{}
    c.ics.Initialize(c)
    c.pcs = PCS{}
    c.pcs.Initialize()
    c.ucs = UCS{}
    c.ucs.Initialize()
    c.ctx = new(CTX)
    c.ctx.Initialize()

    if c.grammar == nil {
        c.grammar = new(grammar)
        c.grammar.groups = make(map[string]group)
        c.grammar.rules = make([]rule, 0, 1024)
    }

    c.sentence_analyser_workers = 1

    c.createEntityType("entity")
    c.initialized = true

    c.ucs.ceres_main = c.root
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
    var P float64 = 1
    var past_entities []RecognizedEntity = make([]RecognizedEntity,
        0, len(idp.analyseResult))
    for _, re := range idp.analyseResult {
        analyser := re.proposer
        if analyser == nil {
            fmt.Println("Recognized :", re,  "\nCERES :", c, "\nAnalyser :", analyser)
            panic("analyser is nil ?")
        }
        P *= analyser.computeP(re, c.ctx, past_entities...)
        past_entities = append(past_entities, re)
    }
    return P
}
