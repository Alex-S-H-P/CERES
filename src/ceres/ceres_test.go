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
    if len(*c.root.children) != 1 {
        t.Errorf("Created a child for the root, but the root did not care")
    } else if (*c.root.children)[0].directTypeOf() == nil {
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
    thing := (*c.root.children)[2].(*EntityType)
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
    thing := (*c.root.children)[2].(*EntityType)
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

func TestSentenceParser(t *testing.T){
    c := new(CERES)
    c.Initialize(1)
    c.createEntityType("caracteristic")
    c.createEntityType("action")
    c.createEntityType("thing")

    c.ParseSentence("hello, I am alexandre")
}

const fileSpace = "test.ceres"

func TestCERES_Saving(t *testing.T) {
    c := new(CERES)
    c.Initialize(1)
    c.createEntityType("caracteristic")
    c.createEntityType("action")
    c.createEntityType("thing")
    c.save(fileSpace)

    c2 := new(CERES)
    c2.Initialize(1)
    e := c2.load(fileSpace)
    if e != nil {
        t.Error(e)
    }
    if len(*c.root.children) != len(*c2.root.children) {
        t.Fail()
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
    var Ar, Br, Cr RecognizedEntity = MakeRecognizedEntity(At, false, false, nil, "A"),
        MakeRecognizedEntity(Bt, false, false, nil, "B"),
        MakeRecognizedEntity(Ct, false, false, nil, "C")

    var g grammar = grammar{rules:[]rule{ruleString("E -> B C"),
                                         ruleString("D -> A E"),
                                         ruleString("B -> B"),
                                         ruleString("C -> C"),
                                         ruleString("A -> A")},
                            groups:make(map[string]group)}

    table := CYK_PARSE([]RecognizedEntity{Ar, Br, Cr}, g)

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
