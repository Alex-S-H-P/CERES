package ceres

import (
	"fmt"
	"strconv"
	"strings"

	api "github.com/Alex-S-H-P/NOT_API"
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

func (c *CERES) Initialize(workers int) {
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
Tries to load the data stored in the mainSave file.
*/
func (c *CERES) Load(mainSave, lastSave string) (error, error) {
	err := c.load(mainSave)
	if err == nil {
		return nil, nil
	} else if lastSave != "" {
		return err, c.load(lastSave)
	} else {
		return err, fmt.Errorf("last save was also not defined. CERES is now blank")
	}
}

/*
Creates a valid entityType.

Adds it to the CERES dictionary.

*/
func (c *CERES) createEntityType(w Word) {
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

func (c *CERES) allOptions(w Word) []*RecognizedEntity {
	var total_capacity int

	icss := c.ics.proposeOptions(w, c.ctx)
	pcss := c.pcs.proposeOptions(w, c.ctx)
	ucss := c.ucs.proposeOptions(w, c.ctx)

	total_capacity = len(icss) + len(pcss) + len(ucss)

	var options = make([]*RecognizedEntity, 0, total_capacity)

	options = append(options, icss...)
	options = append(options, pcss...)
	options = append(options, ucss...)

	return options
}

/*
This function exists for the NOT API.

You need to pass, in "|" separated values, these fields :

 - **word**, the word used to add this onto the database
    - Formulated as `{word}`
    - required
 - **isType**, whether the word refers to a type or a specific request
    - Formulated as `{istype=(y|n)}`
    - required
 - **parent**, what that word is a hyponym of (for example )
    - formulated as `--parent {parent_word} {parent_index}`
    - **not** required
 - **grammar group**
    - formulated as `--ggroup {ggroup_name}`
    - **not** required
*/
func (c *CERES) AddEntryMethod(p *api.Process,
	args string,
	cID uint32) {
	var w string
	var isType bool

	m, err := api.ProcessOptions(args,
		[]string{"word", "isType"},
		map[string]string{"--parent": "{parent_word} {parent_index}",
			"--ggroup": "{ggroup_name}"})
	if err != nil {
		return
	}
	w = m["word"][0]
	isType = strings.ToLower(m["isType"][0]) == "y"

	delete(m, "word")
	delete(m, "isType")

	err = c.AddEntry(Word(w), isType, m)
	fmt.Println(err)
}

func (c *CERES) AddEntry(w Word, isType bool, options map[string][]string) error {

	const errorPrefix = "could not process addEntry request"

	var parent *EntityType = nil
	var ggroup group
	if parentArg, ok := options["--parent"]; ok {
		index, err := strconv.Atoi(parentArg[1])
		if err != nil {
			return fmt.Errorf("%s : %v", errorPrefix, err)
		}
		de := c.ics.entityDictionary[Word(parentArg[0])]
		if len(de.entities) <= index {
			return fmt.Errorf("%s : index given was out of range (%v >= %v)",
				errorPrefix, index, len(de.entities))
		}

		parent, ok = de.entities[index].(*EntityType)
		if !ok {
			return fmt.Errorf("%s : entity at item %v is instance, not type.",
				errorPrefix, index)
		}
	}
	if ggroupArg, ok := options["--ggroup"]; ok {
		if found, ok := c.grammar.groups[ggroupArg[0]]; ok {
			ggroup = found
		} else if isType {
			ggroup = group{name: ggroupArg[0], instanceSolver: nil}
		} else {
			return fmt.Errorf("%s : Group %s not found",
				errorPrefix, ggroupArg[0])
		}
	}

	if isType {
		return c.addEI(w, parent, ggroup)
	} else {
		return c.addET(w, parent, ggroup)
	}
}

func (c *CERES) addEI(w Word, parent *EntityType, ggroup group) error {
	return nil
}

func (c *CERES) addET(w Word, parent *EntityType, ggroup group) error {
	var et = new(EntityType)

	if ggroup.instanceSolver == nil {
		ggroup.instanceSolver = et
	}

	return nil
}
