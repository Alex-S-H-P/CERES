package ceres

import (
    "testing"
    "fmt"
)

func TestCERES_Initialization(t*testing.T){
    ceres := new(CERES)
    ceres.Initialize(1)

    if !ceres.initialized {
        t.Errorf("CERES should have been marked as initialized")
    }

    if !ceres.ics.initialized {
        t.Errorf("ICS should have been initialized")
    }

    if !ceres.pcs.initialized {
        t.Errorf("PCS should have been initialized")
    }

    if ceres.root == nil {
        t.Errorf("CERES should have  ")
    }
}

func TestCeres_TypeAdding(t *testing.T){
    c := new(CERES)
    c.Initialize(1)
    c.createEntityType("caracteristic")
    if len(c.root.links) != 1 {
        t.Errorf("Created a child for the root, but the root did not care")
    } else if (c.root.links)[0].GetB().directTypeOf() == nil {
        t.Errorf("Created a child that did not hook to its parent correctly")
    }

}

func TestICS_LevelsOfAbstractions(t *testing.T) {
    // setup
    c := new(CERES)
    c.Initialize(1)
    c.createEntityType("caracteristic")
    c.createEntityType("action")
    c.createEntityType("thing")
    thing := (c.root.links)[2].GetB().(*EntityType)
    living_being := c.ics.createEntityType("living being")
    thing.addChild(living_being)
    box := c.ics.createEntityType("box")
    thing.addChild(box)
    shmoop := c.ics.createEntityInstance("shmoop", thing)
    John := c.ics.createEntityInstance("john", living_being)

    if i, e := LevelsOfAbstractionBetween(thing, c.root); i != -1 || e != nil {
        t.Errorf("Unexpected result : (%v instead of %v, %v instead of %v)",
            i, -1, e, nil)
    }

    if i, e := LevelsOfAbstractionBetween(c.root, thing); i != 1 || e != nil {
        t.Errorf("Unexpected result : (%v instead of %v, %v instead of %v)",
            i, 1, e, nil)
    }

    if i, e := LevelsOfAbstractionBetween(living_being, c.root); i != -2 || e != nil {
        t.Errorf("Unexpected result : (%v instead of %v, %v instead of %v)",
            i, -2, e, nil)
    }

    if i, e := LevelsOfAbstractionBetween(c.root, living_being); i != 2 || e != nil {
        t.Errorf("Unexpected result : (%v instead of %v, %v instead of %v)",
            i, 2, e, nil)
    }

    if i, e := LevelsOfAbstractionBetween(living_being, box);  e == nil {
        t.Errorf("Unexpected result : (%v instead of %v, %v instead of %v)",
            i, 4, e, "non-nil error")
    }

    if i, e := LevelsOfAbstractionBetween(shmoop, John); e == nil {
        t.Errorf("Unexpected result : (%v instead of %v, %v instead of %v)",
            i, 0, e, "non-nil error")
    }
}

func TestIsTypeOf(t *testing.T) {
    // setup
    c := new(CERES)
    c.Initialize(1)
    c.createEntityType("caracteristic")
    c.createEntityType("action")
    c.createEntityType("thing")
    thing := (c.root.links)[2].GetB().(*EntityType)
    living_being := c.ics.createEntityType("living being")
    thing.addChild(living_being)
    box := c.ics.createEntityType("box")
    thing.addChild(box)
    John := c.ics.createEntityInstance("John", living_being)

    if !IsTypeOf(John, living_being) {
        t.Errorf("IsTypeOf failed")
    }
    if IsTypeOf(living_being, living_being) {
        t.Errorf("Entities should not match with themselves")
    }
}

const fileSpace = "../../var/saveFolder/test.ceres"

func TestCERES_Saving(t *testing.T) {
    c := new(CERES)
    c.Initialize(1)

    c.grammar = &grammar{rules:[]rule{ruleString("NominGroup -> ADJ NOUN"),
        ruleString("NominGroup -> ADJ NominGroup"),
        ruleString("VerbalGroup -> NOUN VERB"),
        ruleString("VerbalGroup -> NominGroup VERB")},
        groups:make(map[string]group)}

    words := []Word{"caracteristic", "action", "thing"}
    c.createEntityType(words[0])
    c.createEntityType(words[1])
    c.createEntityType(words[2])


    thing := (c.root.links)[2].GetB().(*EntityType)
    thing.grammar_group = group{"NOUN", thing}
    caracteristic := (c.root.links)[1].GetB().(*EntityType)
    caracteristic.grammar_group = group{"ADJ", caracteristic}
    action := (c.root.links)[0].GetB().(*EntityType)
    action.grammar_group = group{"VERB", action}

    e := c.save(fileSpace)
    if e != nil {
        t.Error(e)
    }

    c2 := new(CERES)
    c2.Initialize(1)
    e = c2.load(fileSpace)
    if e != nil {
        t.Error(e)
    }
    if len(c.root.links) != len(c2.root.links) {
        t.Errorf("The copy does not have as many links as the original (%v instead of %v)",
                    len(c2.root.links), len(c.root.links))
    }
    for i, w := range words{
        de := c2.ics.entityDictionary[w]
        if len(de.entities) != 1 {
            t.Errorf("The word %s was not saved correctly into the new CERES instance (%v)", w, de)
        } else {
            if de.entities[0].(*EntityType).grammar_group.name != (c.root.links)[i].GetB().(*EntityType).grammar_group.name {
                t.Errorf("Entity for \"%s\" was found in both original and copies, but they did not have the same grammar_groups [%s!=%s]",
                        w, de.entities[0].(*EntityType).grammar_group.name,
                        (c.root.links)[i].GetB().(*EntityType).grammar_group.name)
            }
        }
    }

    c2.grammar.RefreshAllGroups(&(c2.ics))
    for _, rootC := range c2.root.links {
        if !rootC.GetB().(*EntityType).grammar_group.instanceSolver.Equal(rootC.GetB()) {
            t.Error("Ancestors are not given their own grammar_group")
        }
    }
}

func TestRuleMatcher(t*testing.T) {
    r := ruleString("A -> A B")
    if g, b := r.matches("A", "B"); b {
        if g != "A" {
            t.Errorf("Did not match to the right group. Matched to \"%s\" instead of \"A\"", string(g))
        }
    } else {
        t.Errorf("Did not match where there should be a match")
    }
}

func TestCYK(t*testing.T){
    var At, Bt, Ct *EntityType = new(EntityType),
        new(EntityType), new(EntityType)
    var Ar, Br, Cr *RecognizedEntity = new(RecognizedEntity), new(RecognizedEntity), new(RecognizedEntity)
    *Ar, *Br, *Cr = MakeRecognizedEntity(At, false, false, nil, "A"),
        MakeRecognizedEntity(Bt, false, false, nil, "B"),
        MakeRecognizedEntity(Ct, false, false, nil, "C")


    var g *grammar = &grammar{rules:[]rule{ruleString("E -> B C"),
                                         ruleString("D -> A E"),
                                         ruleString("B -> B"),
                                         ruleString("C -> C"),
                                         ruleString("A -> A")},
                            groups:make(map[string]group)}

    table := CYK_PARSE([]RecognizedEntity{*Ar, *Br, *Cr}, g)

    fmt.Println("                ")
    fmt.Println(table.rslt())
    for j := 0; j<3; j++ {
        if len(table[j][j].assignments) != 1 {
            line := "["
            sep := ""
            for _, assignment := range table[j][j].assignments{
                line += fmt.Sprintf("%s \"%s\"", sep, assignment)
                sep = ","
            }
            fmt.Println(line, "]")
            t.Errorf("Too many assignments (%v instead of 1) in the case at location (%v, %v)",
                        len(table[j][j].assignments), j, j)
        }
    }
    if len(table.rslt().assignments) < 1 {
        t.Errorf("Could not get an answer")
    } else if len(table.rslt().assignments) > 1 {
        t.Errorf("Got too many answers (%v instead of 1)", len(table.rslt().assignments))
    }

    table.display_tree(0)
}

func TestClosestCommonAncestor(t *testing.T) {
    A := new(EntityType)
    A.Initialize()
    if ClosestCommonAncestor(A, A) != A {
        t.Errorf("Equal types aren't their own closest ancestor (they should be)")
    }
    ancestorsOfA := make([]*EntityType, 3)
    for i := range(ancestorsOfA) {
        ancestorsOfA[i] = new(EntityType)
        ancestorsOfA[i].Initialize()
        if i == 0 {
            ancestorsOfA[0].addChild(A)
        } else {
            ancestorsOfA[i].addChild(A)
        }
        if ClosestCommonAncestor(A, ancestorsOfA[i]) != ancestorsOfA[i] {
            t.Errorf("Ancestors should be the closest ancestor between themselves and their descendant")
        }
        if ClosestCommonAncestor(ancestorsOfA[i], A) != ancestorsOfA[i] {
            t.Errorf("Ancestors should be the closest ancestor between themselves and their descendant")
        }
    }
}

func TestSentenceParser(t *testing.T){
    c := new(CERES)
    c.Initialize(1)
    words := []Word{"caracteristic", "action", "thing"}
    c.createEntityType(words[0])
    c.createEntityType(words[1])
    c.createEntityType(words[2])


    thing := (c.root.links)[2].GetB().(*EntityType)
    thing.grammar_group = group{"NOUN", thing}
    caracteristic := (c.root.links)[1].GetB().(*EntityType)
    caracteristic.grammar_group = group{"ADJ", caracteristic}
    action := (c.root.links)[0].GetB().(*EntityType)
    action.grammar_group = group{"VERB", action}

    c.createEntityType("hello")
    c.createEntityType("am")
    c.createEntityType("alexandre")
    c.pcs.pronounDictionary["i"] = Pronoun{GNP:12,
        Posessive:false,
        Adjective:false}


    recog, P := c.ParseSentence("hello I am alexandre")
    fmt.Println(P)
    for _, el := range recog {
        if el.proposer == &(c.ucs) {
            t.Errorf("\"%s\" was interpreted by UCS, which is not what it should have been",
                el.s)
        }
    }
}
