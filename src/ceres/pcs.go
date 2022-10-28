package ceres

import (
    "strings"
)

const(
    // gender
    MALE    int8 =  1
    FEMALE  int8 = -1
    NEUTRAL int8 = -2
    UNKNOWN int8 =  2

    // number
    SINGULAR int8 = 0
    PLURAL   int8 = 1
    DUAL     int8 = 2

    // special codes
    SPEAKER     int8 = 127
    DESTINATOR  int8 = -127
)

func GN_to_Code(gender_and_number string) int8 {
    gender_and_number = strings.ToLower(gender_and_number)
    if gender_and_number == "speaker" || gender_and_number == "speak" {
        return SPEAKER
    } else if gender_and_number == "dest" || gender_and_number == "destinator" {
        return DESTINATOR
    }

    var gender_multiplier int8
    switch {
    case strings.Contains(gender_and_number, "male"):
        gender_multiplier = MALE
    case strings.Contains(gender_and_number, "female"):
        gender_multiplier = FEMALE
    case strings.Contains(gender_and_number, "neutral"):
        gender_multiplier = NEUTRAL
    default:
        gender_multiplier = UNKNOWN
    }
    var number_indicator int8=DUAL
    switch {
    case strings.Contains(gender_and_number, "singular") || strings.Contains(gender_and_number, "sing"):
        number_indicator = SINGULAR
    case strings.Contains(gender_and_number, "plural") || strings.Contains(gender_and_number, "pl"):
        number_indicator = PLURAL
    }

    return 8*gender_multiplier + number_indicator
}


type Pronoun struct {
    GenderAndNumber int8
    Posessive bool
    Adjective bool
}

func (p Pronoun) IsDestinator() bool {
    return p.GenderAndNumber == DESTINATOR
}

func (p Pronoun)IsSpeaker() bool {
    return p.GenderAndNumber == SPEAKER
}

func (p Pronoun)Gender() int8 {
    return p.GenderAndNumber / 8
}

func (p Pronoun)Number() int8 {
    return p.GenderAndNumber % 8
}

type PCS struct {
    pronounDictionary map[Word]Pronoun

    initialized bool
}

func (pcs *PCS)Initialize() {
    pcs.pronounDictionary = make(map[Word]Pronoun)

    pcs.initialized = true
}

func (pcs *PCS)IsPronoun(w Word) bool{
    _, ok := pcs.pronounDictionary[w]
    return ok
}

func (pcs *PCS) Match(w Word) Entity{
    // TODO: do this
    return nil
func (pcs*PCS)proposeOptions(w Word, ctx *CTX) []RecognizedEntity {
    if pronoun, ok := pcs.pronounDictionary[w]; ok {
        entities := make([]RecognizedEntity, 0, 64)
        g, n, p := pronoun.GNP_Sep()
        if p == PERSON1 && n == SINGULAR {
            return []RecognizedEntity{MakeRecognizedEntity(ctx.SPEAKER,
                    pronoun.Posessive, false)}
        } else if p == PERSON2 {
            return []RecognizedEntity{MakeRecognizedEntity(ctx.DESTINATOR,
                 pronoun.Posessive, false)}
        }
        for i := 0; i<ctx.expressed_buffer.Len(); i++ {
            buffered := ctx.expressed_buffer.Get(i)
            if buffered.GetGender() == g || g == UNKNOWN || buffered.GetGender() == UNKNOWN {
                if buffered.GetNumber() == n {
                    entities = append(entities, MakeRecognizedEntity(buffered,
                        pronoun.Posessive, false))
                }
            }
        }
        return entities
    } else {
        return nil
    }
}
