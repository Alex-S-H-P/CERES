package ceres

// "strings"

const (
	// gender
	MALE    int8 = 0
	FEMALE  int8 = 1
	NEUTRAL int8 = 2
	UNKNOWN int8 = 3

	// number
	SINGULAR int8 = 0
	PLURAL   int8 = 1
	DUAL     int8 = 2

	// special codes
	PERSON1 int8 = 0
	PERSON2 int8 = 1
	PERSON3 int8 = 2
)

type Pronoun struct {
	GNP       int8
	Posessive bool
	Adjective bool
}

func (p Pronoun) Gender() int8 {
	return (p.GNP / 4) % 4
}

func (p Pronoun) Number() int8 {
	return p.GNP % 4
}

func (p Pronoun) Person() int8 {
	return (p.GNP) / 16
}

func (p Pronoun) GNP_Sep() (int8, int8, int8) {
	return p.Gender(), p.Number(), p.Person()
}

func (p Pronoun) MakeGNP(gender int8, number int8, person int8) Pronoun {
	p.GNP = person*16 + number*4 + person
	return p
}

type PCS struct {
	pronounDictionary map[Word]Pronoun

	initialized bool
}

func (pcs *PCS) Initialize() {
	pcs.pronounDictionary = make(map[Word]Pronoun)

	pcs.initialized = true
}

func (pcs *PCS) IsPronoun(w Word) bool {
	_, ok := pcs.pronounDictionary[w]
	return ok
}

func (pcs *PCS) proposeOptions(w Word, ctx *CTX) []*RecognizedEntity {
	var re = new(RecognizedEntity)

	if pronoun, ok := pcs.pronounDictionary[w]; ok {
		entities := make([]*RecognizedEntity, 0, 64)
		g, n, p := pronoun.GNP_Sep()
		if (p == PERSON1 && n == SINGULAR) || (p == PERSON2) {

			*re = MakeRecognizedEntity(ctx.SPEAKER,
				pronoun.Posessive, false, pcs, string(w))
			return []*RecognizedEntity{re}
		}
		for i := 0; i < ctx.expressed_buffer.Len(); i++ {
			buffered := ctx.expressed_buffer.Get(i).Entity()
			if buffered.GetGender() == g || g == UNKNOWN || buffered.GetGender() == UNKNOWN {
				if buffered.GetNumber() == n {
					*re = MakeRecognizedEntity(buffered,
						pronoun.Posessive, false, pcs, string(w))
					entities = append(entities, re)
				}
			}
		}
		return entities
	} else {
		return nil
	}
}

func (pcs *PCS) computeP(re *RecognizedEntity, ctx *CTX,
	previous_elements ...*RecognizedEntity) float64 {
	// TODO: fix this
	return .25
}

func (pcs *PCS) name() string { return "PCS" }
