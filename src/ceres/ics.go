package ceres

import (
    "sync"
    "CERES/src/utils"
)
type ICS struct {
    entityDictionary map[Word]*DictionaryEntry
    eDMutex sync.RWMutex

    initialized bool

    master *CERES
}

func (ics *ICS)Initialize(master *CERES) {
    ics.entityDictionary = make(map[Word]*DictionaryEntry)
    ics.master = master
    ics.initialized = true
}

type DictionaryEntry struct {
    entities []Entity
    DEMutex sync.RWMutex
}

/*
Returns all entities in the form of a generator.
May return an entity multiple time
*/
func (ics*ICS)allEntities() utils.Generator[Entity] {
    var mg *utils.MapGenerator[Word,*DictionaryEntry]
    // we look at all of the values
    key_gen := mg.Values(ics.entityDictionary)
    // for each value, we can access its entities attribute
    transform := func(de*DictionaryEntry)[]Entity {return de.entities}
    // we therefore have a generator of slice of entities
    entities_slice_gen := utils.Transform[*DictionaryEntry, []Entity](key_gen, transform)
    // and we can parse through a slice
    mesa_parser := func(s[]Entity)utils.Generator[Entity] {return utils.SliceGenerator[Entity](s)}
    // we combine them and return the last one
    return utils.Combine[[]Entity, Entity](entities_slice_gen, mesa_parser)
}

func (ics *ICS)createEntityType(w Word) *EntityType {
    et := new(EntityType)
    et.Initialize()

    ics.eDMutex.RLock()
    if DE, ok := ics.entityDictionary[w]; ok {
        ics.eDMutex.RUnlock()
        DE.DEMutex.Lock()
        DE.entities = append(DE.entities, Entity(et))
        DE.DEMutex.Unlock()
    } else {
        ics.eDMutex.RUnlock()
        DE := DictionaryEntry{entities:[]Entity{et}}
        ics.eDMutex.Lock()
        ics.entityDictionary[w] = &DE
        ics.eDMutex.Unlock()
    }

    return et
}


func (ics*ICS)createEntityInstance(w Word, et *EntityType) *EntityInstance{
    ei := new(EntityInstance)
    ei.Initialize()

    ics.eDMutex.RLock()
    if DE, ok := ics.entityDictionary[w]; ok {
        ics.eDMutex.RUnlock()
        DE.DEMutex.Lock()
        DE.entities = append(DE.entities, ei)
        DE.DEMutex.Unlock()
    } else {
        ics.eDMutex.RUnlock()
        DE := DictionaryEntry{entities:[]Entity{ei}}
        ics.eDMutex.Lock()
        ics.entityDictionary[w] = &DE
        ics.eDMutex.Unlock()
    }

    et.addChild(ei)

    return ei
}

func (ics*ICS)listOptionStrict(w Word, de *DictionaryEntry) []RecognizedEntity {

    var res []RecognizedEntity = make([]RecognizedEntity, 0, len(de.entities))

    for _, entity := range de.entities {
        re := MakeRecognizedEntity(entity, false, false, ics, string(w))
        res = append(res, re)
    }
    return res
}

func (ics*ICS)proposeOptions(w Word, ctx *CTX) []RecognizedEntity{
    var answer []RecognizedEntity
    if de, ok := ics.entityDictionary[w]; ok {
        answer = ics.listOptionStrict(w, de)
    }

    return answer
}

func (ics*ICS)computeP(re RecognizedEntity, ctx*CTX,
        previous_elements ...RecognizedEntity) float64{
    // TODO: code this
    return .5
}
