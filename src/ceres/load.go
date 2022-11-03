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

	c.ics.load(b1, b2)
	fmt.Println(c, "\n", c.root)
	c.pcs.load(b3)
	fmt.Println(c)
	c.ucs.load(b4)
	fmt.Println(c)
	return err
}

/*
Loads ics.

Should not run while any other goroutine has access to the ICS.
*/
func (ics*ICS)load(b1,b2 *[]byte){
	// setting b2
	m := make(map[int]Entity)
	C := strings.Split(string(*b2), UnitSEP)
	for _, descr := range C {
		c := strings.Split(descr, UnderSEP)
		index, err := strconv.Atoi(c[1])
		if err != nil {
			panic(err)
		}
		parent_index, err := strconv.Atoi(c[2])
		if err != nil {
			panic(err)
		}
		var e Entity
		switch c[0] {
		case "type":
			et := new(EntityType)
			if parent_index >= 0 {
				m[parent_index].(*EntityType).addChild(et)
			} else {
				fmt.Println("Setting root", et)
				ics.master.root = et
				ics.master.ucs.ceres_main = et
			}
			e = Entity(et)
			et.attributes = new(AttributeTypeList)
			for i := 3; i<len(c); i++ {
				Cidx, err := strconv.Atoi(c[i])
				if err != nil {
					panic(err)
				}
				et.attributes.attrs = append(et.attributes.attrs, m[Cidx].(*EntityType))
			}
		case "inst":
			ei := new(EntityInstance)
			m[parent_index].(*EntityType).addChild(ei)
			e = Entity(ei)
			ei.values = new(AttributeInstanceList)
			for i := 3; i<len(c); i++ {
				d := strings.SplitN(c[i], ":", 2)
				Cidx, err := strconv.Atoi(d[0])
				if err != nil {
					panic(err)
				}
				attr := m[Cidx].(*EntityType)
				ei.values.attrs = append(ei.values.attrs, attr)
				ei.values.values[attr] = Word(d[1])
			}
		}
		m[index] = e
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
		fmt.Println("LOADING PCS", c)
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

func (ucs*UCS)load(b*[]byte){
	C := strings.Split(string(*b), UnderSEP)
	ucs.unrecognized_words = make([]Word, 0, len(C))
	for _, c := range C {
		ucs.unrecognized_words = append(ucs.unrecognized_words, Word(c))
	}
}

func (c *CERES) save(fn string) error {
	b1, b2, e := c.ics.save()
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

/*
Returns an index for every entity, while adding the description of the entity to a slice.

The indexes of two equal entities will be equal.
Otherwise, every index is unique.

Parent entities are always stored into the dictionary before the children.
*/
func safeIndexEntity(e Entity, m map[Entity]int, entityDict *[]byte) int {

	i, ok := indexEntity(e, m)
	if ok { // we have a new entity. Let's add it to the list,
		var s string
		switch e.(type) {
		case *EntityType:
			et := e.(*EntityType)
			fmt.Println("Indexing", *et, "...")
			if et.parent != nil {
				s = "␟type␣" + strconv.Itoa(i) + "␣" +
					strconv.Itoa(safeIndexEntity(et.parent, m, entityDict)) + "␣"
			} else {
				s = "␟type␣" + strconv.Itoa(i) + "␣-1␣"
			}
			for _, attr := range et.attributes.attrs {
				fmt.Println("Attr ::", attr)
				s += strconv.Itoa(safeIndexEntity(attr, m, entityDict)) + "␣"
			}
		case *EntityInstance:
			ei := e.(*EntityInstance)
			s = "␟inst␣" + strconv.Itoa(safeIndexEntity(ei, m, entityDict)) + "␣" +
				strconv.Itoa(safeIndexEntity(ei.typeOf, m, entityDict)) + "␣"
			for _, attr := range ei.values.attrs {
				val := ei.values.values[attr]
				s += strconv.Itoa(safeIndexEntity(attr, m, entityDict)) + ":" +
					string(val) + "␣"
			}
		}
		s = strings.TrimSuffix(s, "␣")
		if len(*entityDict) == 0 {
			s = strings.TrimPrefix(s, "␟")
		}
		fmt.Println("Adding", s+"("+strconv.Itoa(len(s))+") to", string(*entityDict))
		*entityDict = append(*entityDict, []byte(s[:len(s)])...)
	}
	return i
}

func (ics *ICS) save() ([]byte, []byte, error) {
	b1, b2 := make([]byte, 0, 2048*2048), make([]byte, 0, 2048)
	var m map[Entity]int = make(map[Entity]int)
	var initial bool = true
	for w, entries := range ics.entityDictionary {
		fmt.Println("entry [\""+string(w)+"\"]->", entries)
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
			fmt.Println("EDITING B1 :", understood, first)
			first = false
			b1 = append(b1, []byte(understood)...)
		}

	}
	fmt.Println(string(b1), string(b2))
	return b1[:len(b1)], b2[:len(b2)], nil
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
