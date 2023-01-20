package ceres

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	sto_fmt_split = "␞"
	UnderSEP = "␣"
	UnitSEP  = "␟"
)

const storage_format string = "%s" + sto_fmt_split + "%s" + sto_fmt_split +
	"%s" + sto_fmt_split + "%s" + sto_fmt_split + "%d\n"

func (c *CERES) load(fn string) error {
	var b1, b2, b3, b4 *[]byte = new([]byte), new([]byte), new([]byte), new([]byte)
	var workers int
	var contents *string = new(string)
	f, e := os.Open(fn)
	if e != nil {
		return e
	}
	defer f.Close()

	_, err := fmt.Fscanf(f, "%s\n", contents)
	if err != nil {
		return err
	}

	C := strings.Split(*contents, sto_fmt_split)
	*b1, *b2, *b3, *b4 = []byte(C[0]), []byte(C[1]), []byte(C[2]), []byte(C[3])

	workers, err = strconv.Atoi(C[4])
	c.sentence_analyser_workers = workers
	var future_action_groups map[string]int

	c.grammar, future_action_groups, err = grammar_load(fn + ".grammar")
	if err != nil {return err}

	solver := c.ics.load(b1, b2, c.grammar.groups)
	//fmt.Println(c, "\n", c.root)
	c.grammar.resolve_future_actions_on_loading(solver, future_action_groups)

	c.pcs.load(b3)
	//fmt.Println(c)
	c.ucs.load(b4)
	//fmt.Println(c)
	return err
}


func (et*EntityType) load(c[]string,
			grammar_groups map[string]group,
			m map[int] Entity) {
	et.attributes = new(AttributeTypeList)
	et.links = make([]Link, 0, len(c))
	var i int = 2
	fmt.Println(c)
	for i < len(c)-2 {

		fmt.Println(i, ">", c[i])
		switch {
		case len(c[i]) > 2 && c[i][0] == '@':
			attribute, err := strconv.Atoi(c[i])
			if err != nil {panic(err)}

			et.attributes.attrs = append(et.attributes.attrs,
				m[attribute].(*EntityType))
		case len(c[i]) > 1 && c[i][0] == '@':
			i++
			continue
		case len(c[i])>1 &&c[i][0]!= '@':
			d := strings.Split(c[i], "-")
			typeOfLink,linkTo_string := d[0], d[1]
			linkTo, err := strconv.Atoi(linkTo_string)
			destinationFound := m[linkTo]
			if destinationFound == nil {
				panic(fmt.Sprintf("Could not find item @%v in %v",
					linkTo, m))
			}
			if err != nil {panic(err)}

			/*fmt.Println("CURRENT", c[0], c[1], "#", typeOfLink, linkTo, "#",
				FindListType(typeOfLink).typeOfLink(),
				"\""+typeOfLink+"\"",
				destinationFound)*/
			link := FindListType(typeOfLink).set(et, destinationFound)
			et.links = append(et.links, link)
		}
		i++
	}
	et.word = Word(c[len(c)-2])
	et.grammar_group = grammar_groups[c[len(c)-1]]
}

func (ei*EntityInstance) load(c[]string,
			grammar_groups map[string]group,
			m map[int] Entity) {

	ei.values = new(AttributeInstanceList)
	ei.values.values = make(map[*EntityType]Word)
	typeOfIndex, err := strconv.Atoi(c[1])
	if err != nil { panic(err)	}
	ei.typeOf = m[typeOfIndex].(*EntityType)
	for _, attrSegment := range c[2:] {
		d := strings.SplitN(attrSegment, ":", 2)

		attrIndex, err := strconv.Atoi(d[0])
		if err != nil {panic(err)}

		var attr *EntityType = m[attrIndex].(*EntityType)
		ei.values.values[attr] = Word(d[1])
	}
}

func newEntityToLoad(fields []string) Entity {
	switch fields[0] {
	case "type":
		et := new(EntityType)
		return et
	case "inst":
		ei := new(EntityInstance)
		return ei
	default:
		panic(fmt.Sprintf("Bad entity category \"%s\" (either \"inst\" or \"type\")",  fields[0]))
	}
}


/*
Loads ics.

Should not run while any other goroutine has access to the ICS.
*/
func (ics*ICS)load(b1,b2 *[]byte, grammar_groups map[string]group) map[int]Entity{
	// setting b2
	m := make(map[int]Entity)
	C := strings.Split(string(*b2), UnitSEP)
	for _, descr := range C {
		c := strings.Split(descr, UnderSEP)
		index, err := strconv.Atoi(c[1])
		if err != nil {
			panic(err)
		}
		m[index] = newEntityToLoad(c)
	}
	for _, descr := range C {
		c := strings.Split(descr, UnderSEP)
		index, err := strconv.Atoi(c[1])
		if err != nil {
			panic(err)
		}
		m[index].load(c, grammar_groups, m)

		if len(m[index].directTypeOf()) == 0 {
			ics.master.root = m[index].(*EntityType)
		}
	}

	// b2 set
	// setting b1
	B := strings.Split(string(*b1), UnitSEP)
	for _, b := range B {
		b_ := strings.Split(b, "⚯")
		index, err := strconv.Atoi(b_[1])
		if err != nil {
			panic(err)
		}
		var w = Word(b_[0])
		if DE, ok := ics.entityDictionary[w];ok {
			DE.entities = append(DE.entities, m[index])
		} else {
			DE := DictionaryEntry{entities:[]Entity{m[index]}}
			ics.entityDictionary[w] = &DE
		}
	}

	// b1 set
	return m
}

func (pcs*PCS)load(b*[]byte){

	if pcs.pronounDictionary == nil {
		pcs.pronounDictionary = make(map[Word]Pronoun)
	}
	if len(*b) == 0 {
		return
	}

	C := strings.Split(string(*b), UnderSEP)
	fmt.Println(C, b, len(C), len(*b))
	for _, c := range C {
		//fmt.Println("LOADING PCS", c)
		w := Word(c[:len(c)-2])
		gnp := rune(c[len(c)-2])
		t, err := strconv.Atoi(string(c[len(c)-1]))
		if err != nil {
			panic(err)
		}
		p := Pronoun{GNP:int8(gnp)}
		p.Posessive = (t / 2) == 1
		p.Adjective = (t % 2) == 1
		pcs.pronounDictionary[w] = p
	}
}

func grammar_load(path string) (*grammar, map[string]int, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, nil, e
	}
	defer f.Close()

	var g *grammar = new(grammar)
	var m = make(map[string]int)


	for {
		var pline = new(string)
		fmt.Fscanf(f, "%s\n", pline)

		line := (*pline)
		if (line) == "␃" {
			break
		}

		if len(line) == 0 {continue}

		sline := strings.Split(line, UnitSEP)
		//fmt.Printf("loader : \"%s\"=>\"%s\", \"%s\"\n", line, sline[0], sline[1])

		if len(sline) != 2 {
			return nil, nil, fmt.Errorf("Cannot process line \"%s\" (%v elements found instead of 2)", line, len(sline))
		}
		if len(sline[0]) != 0{
			r := ruleString(sline[0])
			g.rules = append(g.rules, r)
		}
		if len(sline[1]) != 0 { // we give an entityInstance -> group link
			b := strings.Split(sline[1], ">")
			//fmt.Println(b[0], "|>", b[1])
			id, err := strconv.Atoi(b[1])
			if err != nil {
				return nil, nil, fmt.Errorf("Could not extract entityID :%s", err.Error())
			}
			m[b[0]] = id
		}
	}
	g.groups = make(map[string]group)

	return g, m, nil
}

func (g*grammar) resolve_future_actions_on_loading(solver map [int]Entity, future_actions map[string]int){
	//fmt.Println("resolver : ", solver, future_actions)
	for name, entityID := range future_actions {
		et := solver[entityID].(*EntityType)
		et.grammar_group = g.find(name)
		//fmt.Printf("resolver \"%s\" : %v for %v\n", name, entityID, et)
	}
}

func (ucs*UCS)load(b*[]byte){
	C := strings.Split(string(*b), UnderSEP)
	ucs.unrecognized_words = make([]Word, 0, len(C))
	for _, c := range C {
		ucs.unrecognized_words = append(ucs.unrecognized_words, Word(c))
	}
}

func (c *CERES) save(fn string) error {


	b1, b2, m, e := c.ics.save()
	if e != nil {
		return e
	}

	e = c.grammar.save(fn + ".grammar", m)
	if e != nil {
		return e
	}

	b3, e := c.pcs.save()
	if e != nil {
		return e
	}

	b4, e := c.ucs.save()
	if e != nil {
		return e
	}

	f, e := os.Create(fn)
	if e != nil {
		return e
	}
	defer f.Close()
	_, err := fmt.Fprintf(f, storage_format, b1, b2, b3, b4, c.sentence_analyser_workers)
	return err
}

/*
Grants every entity an index

If the map already contains this entity (or a deep copy), sends the same index.

The boolean value indicates whether we learned this entity during this function call.
*/
func indexEntity(e Entity, m map[Entity]int) (int, bool) {
	for k, v := range m {
		if e.Equal(k) {
			return v, false
		}
	}
	n := len(m)
	m[e] = n
	return n, true
}

func (et*EntityType) store(i int, m map[Entity]int, entityDict *[]byte) string {
	var s string = "␟type␣" + strconv.Itoa(i) + UnderSEP
	//fmt.Println("Indexing", et, "...")
	for _, link := range et.links {
		s += fmt.Sprintf("%s-%d␣", link.typeOfLink(),
					safeIndexEntity(link.GetB(), m, entityDict))
	}
	for _, attr := range et.attributes.attrs {
		fmt.Println("Attr ::", attr)
		s += "@" + strconv.Itoa(
					safeIndexEntity(attr, m, entityDict)) + UnderSEP
	}
	s += string(et.word) + UnderSEP
	s += et.grammar_group.String() + UnderSEP

	return s
}

func (ei*EntityInstance) store(i int, m map[Entity]int, entityDict *[]byte) string {
	var s string
	s = "␟inst␣" + strconv.Itoa(safeIndexEntity(ei, m, entityDict)) + "␣" +
	strconv.Itoa(safeIndexEntity(ei.typeOf, m, entityDict)) + "␣"
	for _, attr := range ei.values.attrs {
		val := ei.values.values[attr]
		s += strconv.Itoa(safeIndexEntity(attr, m, entityDict)) + ":" +
		string(val) + "␣"
	}

	return s
}

/*
Returns an index for every entity, while adding the description of the entity to a slice.

The indexes of two equal entities will be equal.
Otherwise, every index is unique.

Parent entities are always stored into the dictionary before the children.
*/
func safeIndexEntity(e Entity, m map[Entity]int, entityDict *[]byte) int {

	i, ok := indexEntity(e, m)
	if ok { // we have a new entity. Let's add it to the list,
		var s string = e.store(i ,m, entityDict)

		s = strings.TrimSuffix(s, "␣")
		if len(*entityDict) == 0 {
			s = strings.TrimPrefix(s, "␟")
		}
		//fmt.Println("Adding", s+"("+strconv.Itoa(len(s))+") to", string(*entityDict))
		*entityDict = append(*entityDict, []byte(s[:len(s)])...)
	}
	return i
}

func (ics *ICS) save() ([]byte, []byte, map[Entity]int, error) {
	b1, b2 := make([]byte, 0, 2048*2048), make([]byte, 0, 2048)
	var m map[Entity]int = make(map[Entity]int)
	var initial bool = true
	for w, entries := range ics.entityDictionary {
		//fmt.Println("entry [\""+string(w)+"\"]->", entries)
		if !initial {
			b1 = append(b1, []byte(UnitSEP)...)
		}
		initial = false
		var first bool = true
		for _, e := range entries.entities {
			i := safeIndexEntity(e, m, &b2)
			var understood string
			if first {
				understood = string(w) + "⚯" + strconv.Itoa(i)
			} else {
				understood = "␣" + string(w) + "⚯" + strconv.Itoa(i)
			}
			//fmt.Println("EDITING B1 :", understood, first)
			first = false
			b1 = append(b1, []byte(understood)...)
		}

	}
	fmt.Println(string(b1), string(b2))
	return b1[:len(b1)], b2[:len(b2)], m, nil
}

func (pcs *PCS) save() ([]byte, error) {
	b := make([]byte, 0, 32*len(pcs.pronounDictionary))
	var first bool = true
	for w, pronoun := range pcs.pronounDictionary {

		var s, t string
		if !first {
			s = UnderSEP
		}
		switch {
		case pronoun.Posessive && pronoun.Adjective:
			t = "3"
		case pronoun.Posessive:
			t = "2"
		case pronoun.Adjective:
			t = "1"
		default:
			t = "0"
		}
		s += string(w) + string(rune(pronoun.GNP)) + t
		first = false
		b = append(b, []byte(s)...)
	}
	return b, nil
}

func (ucs *UCS) save() ([]byte, error) {
	b := make([]byte, 0, 12*len(ucs.unrecognized_words))
	var first bool = true
	for _, w := range ucs.unrecognized_words {
		if len(string(w)) > 0 {
			if !first {
				b = append(b, []byte(UnderSEP)...)
			}
			b = append(b, []byte(w)...)
			fmt.Println(string(b))
		}
	}
	fmt.Println(string(b))
	return b[:len(b)], nil
}


func (g*grammar)save(path string, m map[Entity]int)error {
	if g == nil {return fmt.Errorf("Cannot save non-existant grammar")}
	var contents string

	f, e := os.Create(path)
	if e != nil {
		return e
	}
	defer f.Close()

	key_arr := make([]*EntityType, len(m))
	var i int = 0
	for e := range m {
		if et, ok := e.(*EntityType); ok {
			key_arr[i] = et
			i++
		}
	}

	for line := 0; line < len(m) || line < len(g.rules); line ++ {
		var ruleSub, entitySub string
		if line < len(g.rules) {
			ruleSub = g.rules[line].String()
		}
		if line < len(key_arr){
			if key_arr[line] != nil {
				et := key_arr[line]

				entitySub = fmt.Sprintf("%s>%v", et.grammar_group.String(),m[Entity(et)])

			} else if line >= len(g.rules){
				break
			}
		}
		contents += ruleSub + UnitSEP + entitySub + "\n"
	}

	_, err := fmt.Fprintf(f, "%s␃\n", contents)
	return err
}
