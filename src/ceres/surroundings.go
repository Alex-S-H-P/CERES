package ceres

import (
    "math"
    "sync"
    "strings"
    "strconv"
    "encoding/binary"
)

type surroundingList struct {
    surr []*surrounding
    rwm sync.RWMutex
}

func (sl*surroundingList)save(m map[Entity]int) []byte {
    b := make([]byte, 0, 20*len(sl.surr))
    first := true
    for _, sur := range sl.surr {
        if first {
            b = append(b, []byte(UnderSEP)...)
            first = false
        }
        b = append(b, sur.save(m)...)
    }
    return b
}

func (sl*surroundingList)load(b string, m map[int]Entity) {
    if len(b) == 0 {return}
    C := strings.Split(b, UnderSEP)
    for _, c := range C {
        var s = new(surrounding)
        s.load(c, m)
        sl.surr = append(sl.surr, s)
    }
}

func (sl*surroundingList)Len()int {return len(sl.surr)}

func (sl*surroundingList)Add(s *surrounding){
    sl.rwm.Lock()
    sl.surr = append(sl.surr, s)
    sl.rwm.Unlock()
}

func (sl*surroundingList)Get(i int)*surrounding{
    sl.rwm.RLock()
    defer sl.rwm.RUnlock()
    return sl.surr[i]
}

type surrounding struct {
    prox []surroundingToken
    coherence float64
}

func (s*surrounding)save(m map[Entity]int)[]byte{
    b := make([]byte, 0, 8+len(s.prox)*8)
    bits := math.Float64bits(s.coherence)
    bytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(bytes, bits)
    first := true
    b = append(b, bytes...)
    for _, p := range s.prox {
        if first {
            b = append(b, []byte(":")...)
            first = false
        }
        b = append(b,[]byte(strconv.Itoa(m[p.stype]))...)
        b = append(b, []byte(":")...)
        b = append(b,[]byte(strconv.Itoa(p.pos))...)
    }
    return b
}

func (s*surrounding)load(b string, m map[int]Entity){
    bits := binary.LittleEndian.Uint64([]byte(b[:8]))
    s.coherence = math.Float64frombits(bits)
    parts := strings.Split(b[8:], ":")
    for i:=0; i <len(parts); i+=2 {
        token := new(surroundingToken)
        // stype
        idx, err := strconv.Atoi(parts[i])
        if err != nil {panic(err)}
        token.stype = m[idx].(*EntityType)
        // pos
        pos, err := strconv.Atoi(parts[i+1])
        if err != nil {panic(err)}
        token.pos = pos
        s.prox = append(s.prox,*token)
    }
}

type surroundingToken struct {
    stype *EntityType
    pos  int
}
