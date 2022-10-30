package ceres

import (
	"fmt"
	re "regexp"
	"sync"
	"time"
	"CERES/src/utils"
)

var RegexpToken *re.Regexp = re.MustCompile(utils.TokenPattern)

type dijkstraPossibility struct {
	curP float64
	tokens []string
	analysedUntil int// this index excluded
	analyseResult []RecognizedEntity
}

func (idp *dijkstraPossibility) makeChild(re RecognizedEntity, offset int) *dijkstraPossibility{
	child := new(dijkstraPossibility)
	child.curP = 0
	child.tokens = idp.tokens
	child.analysedUntil = idp.analysedUntil + offset
	child.analyseResult = append(make([]RecognizedEntity, 0, cap(idp.analyseResult)),
				idp.analyseResult...)
	child.analyseResult = append(child.analyseResult, re)
	fmt.Println("CHILD :", child.analyseResult)
	return child
}

func (idp *dijkstraPossibility) completed() bool { return len(idp.tokens) == idp.analysedUntil }

type organizedDijkstraList struct {
	list []*dijkstraPossibility
	rwm sync.Mutex
}

func (odl *organizedDijkstraList)pop() *dijkstraPossibility{
	odl.rwm.Lock()
	defer odl.rwm.Unlock()

	if len(odl.list) > 0 {
		a := odl.list[0]
		odl.list = odl.list[1:]
		return a
	} else {
		return nil
	}
}

func (odl *organizedDijkstraList)put(possibility *dijkstraPossibility) {
	fmt.Println("Putting", possibility, "into", *odl)
	fmt.Println("lockable")
	odl.rwm.Lock()
	fmt.Println("-- locked")
	defer odl.rwm.Unlock()

	i := odl.getRank(possibility.curP)
	tmp := odl.list[i:]
	odl.list = append(odl.list[:i], possibility)
	odl.list = append(odl.list, tmp...)
}

func (odl *organizedDijkstraList)getRank(P float64) int {
	a, b := 0, len(odl.list)
	for a - b > 1 {
		c := (b-a)/2 + a
		if P > odl.list[c].curP {
			// we are in descending order, therefore i_P should be before c
			b = c
		} else {
			a = c
		}
	}
	return b
}

func (odl *organizedDijkstraList)finished()bool {
	if len(odl.list) > 0 {
		return odl.list[0].completed()
	}
	return true
}

func (c *CERES)ParseSentence(sentence string)[]RecognizedEntity{
	tokens := RegexpToken.FindAllString(sentence, len(sentence)/2)
	var n int = len(tokens)
	fmt.Println(tokens[:n])

	var odl = new(organizedDijkstraList)
	odl.list = make([]*dijkstraPossibility, 0, 2048)
	odl.list = append(odl.list, &dijkstraPossibility{curP:0,
		 tokens:tokens, analysedUntil:0})
	 wg := new(sync.WaitGroup)

	if a := c.dijkstra_main(odl, wg); a != nil {
		return a.analyseResult
	} else {
		return nil
	}
}

func (c *CERES) dijkstra_main(odl *organizedDijkstraList,
	wg *sync.WaitGroup) *dijkstraPossibility {
	// setting up all workers
	handle_chan := make(chan *dijkstraPossibility, 4*c.sentence_analyser_workers)
	counter_chans := make([]chan *sync.WaitGroup, 0, c.sentence_analyser_workers)
	now_empty_chans := make([]chan bool, 0, c.sentence_analyser_workers)

	wg.Add(c.sentence_analyser_workers)
	for w_id := 0; w_id < c.sentence_analyser_workers; w_id++ {
		counter_chan := make(chan *sync.WaitGroup)
		counter_chans = append(counter_chans, counter_chan)
		now_empty_chan := make(chan bool)
		now_empty_chans = append(now_empty_chans, now_empty_chan)
		go c.worker_main(odl, wg, handle_chan, counter_chan, now_empty_chan)
	}

	for {
		odl.rwm.Lock()
		candidate := odl.pop()
		odl.rwm.Unlock()
		if candidate.completed() {
			// this candidate is the best one
			return candidate
		} else  {
			// we inform all workers
			counter := new(sync.WaitGroup)
			counter.Add(c.sentence_analyser_workers)
			for _, channel := range counter_chans {
				channel <- counter
			}
			// we evolve the possibility
			c.evolve(candidate, handle_chan, now_empty_chans)
			counter.Wait()
		}
	}
}

func (c *CERES)worker_main(odl *organizedDijkstraList, wg *sync.WaitGroup,
	handle_chan chan *dijkstraPossibility, counter_chan chan *sync.WaitGroup,
	now_empty_chan chan bool) {
	if wg != nil {
		defer wg.Done()
	}

MAINLOOP:
	for {
		select {
		case <- time.After(3*time.Second):
			break MAINLOOP
		case counter := <- counter_chan:
			for {
				select {
				case poss := <- handle_chan:
					poss.curP = c.computeP(poss)
					odl.put(poss)
				case <-now_empty_chan:
					counter.Done()
					continue MAINLOOP
				}
				fmt.Println("put element.Now waiting for more")
			}
		}
	}
}

func (c *CERES)parseOptions(w Word, handler chan *dijkstraPossibility,
	cur *dijkstraPossibility, offset int) {
	var i int
	var mutex, pm, im sync.Mutex
	pm.Lock()
	im.Lock()
	go func () {
		for _, proposition := range c.pcs.proposeOptions(w, c.ctx) {
			fmt.Println("PCS : ", proposition)
			handler <- cur.makeChild(proposition, offset)
			mutex.Lock()
			i ++
			mutex.Unlock()
		}
		pm.Unlock()
	}()
	go func () {
		for _, proposition := range c.ics.proposeOption(w, c.ctx) {
			handler <- cur.makeChild(proposition, offset)
			mutex.Lock()
			i ++
			mutex.Unlock()
		}
		im.Unlock()
	}()
	pm.Lock()
	defer pm.Unlock()
	im.Lock()
	defer im.Unlock()

	if i == 0 {
		p := c.ucs.proposeOptions(w, c.ctx)[0]
		handler <- cur.makeChild(p, offset)
	}
}

func (c *CERES)evolve(cur *dijkstraPossibility,
	handler chan *dijkstraPossibility, DoneWarner []chan bool){
	// safety
	if cur.analysedUntil == len(cur.tokens) {
		panic(fmt.Errorf("Cannot evolve when there is no token to evolve onto"))
	}
	// proper function
	curTokenT := recognizeType(cur.tokens[cur.analysedUntil])
	if curTokenT == TOKEN_TYPE_WORD {
		w := Word(cur.tokens[cur.analysedUntil])
		c.parseOptions(w, handler, cur, 1)
		for j := cur.analysedUntil + 1; j < len(cur.tokens); j++ {
			tokenT := recognizeType(cur.tokens[j])
			if tokenT == TOKEN_TYPE_WORD {
				c.parseOptions(w, handler, cur, j - cur.analysedUntil + 1)
			} else {
				break
			}
		}
	} else {
		handler <- cur.makeChild(c.makeNonWordEntity(cur.tokens[cur.analysedUntil]), 1)
	}

	for _, warner := range DoneWarner {
		go func(){warner <- true}()
	}
}


func (c *CERES)makeNonWordEntity(token string)RecognizedEntity {
	// TODO handle all cases.
	switch recognizeType(token) {
	case TOKEN_TYPE_PRIC:
		// make a temporary price entity
	case TOKEN_TYPE_CURR:
		// replace with the currency's entity.
	case TOKEN_TYPE_INTN, TOKEN_TYPE_NUMB:
		// make temporary number entity
	case TOKEN_TYPE_UNKN:
		// fetch unknown. If failed make one.
	}

	return RecognizedEntity{}
}
