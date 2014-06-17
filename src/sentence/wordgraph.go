package sentence

import (
    "fmt"
    "dict"
    "syllable"
)

type hzIdInfo struct {
    id      int
    start   int
    end     int
    flag    int
    lpath   []*path
    rpath   []*path
    bestcost    float64
}

type WordGraph struct {
    wg      [][]*hzIdInfo
    dt      *dict.SysDict
}

type backPair   struct {
    start   int
    syl     string
}

const (
    MixHz = 0x100
)

func MakeWordGraph(sg *SylGraph,dt *dict.SysDict, py string)(wg *WordGraph) {
    ln := len(sg.sg)
    wg = new(WordGraph)
    wg.wg = make([][]*hzIdInfo, ln)
    wg.dt = dt
    //need a sylgraph to store syls end
    endsylgraph := make([][]*syllable.PyIdUnit, ln+1)
  
    for i:=0; i < ln ; i++ {
        wg.wg[i] = make([]*hzIdInfo,0)
        for _, pyidunit := range sg.sg[i] {
            id := pyidunit.GetId()
            syl := syllable.GetSylById(id)
            if syl == "" {
                panic("error syl")
            }

            hzids := dt.GetTermsIdByPy(syl)
            start := pyidunit.GetStart()
            if start != i {
                panic("error syl start")
            }

            end := pyidunit.GetEnd()
            flags := pyidunit.GetFlag()

            for _, id:= range hzids {       //first, push simple hz
                hz := &hzIdInfo{id:id, start:start, end:end, flag:flags}
                wg.wg[i] = append(wg.wg[i], hz)
            }

            if endsylgraph[end] == nil {
                endsylgraph[end] = make([]*syllable.PyIdUnit, 0)
            }
            endsylgraph[end] = append(endsylgraph[end],pyidunit)

            //广度遍历回溯
            index := 0
            newstack := make([]backPair,0, 10)
            newstack = append(newstack, backPair{start, syl} )
            for index != len(newstack) {
                sstart := newstack[index].start
                csyl := newstack[index].syl

                for _, ssyl := range endsylgraph[sstart] {
                    comsyl := syllable.GetSylById(ssyl.GetId()) + "'"+csyl
                    shzids := dt.GetTermsIdByPy(comsyl)
                    if shzids == nil {
                        continue
                    }
                    newstart := ssyl.GetStart()
                    newflags :=  flags | MixHz
                    for _, id := range shzids {
                        hz := &hzIdInfo{id:id, start:newstart, end:end,flag:newflags}
                        wg.wg[newstart] = append(wg.wg[newstart], hz)
                    }
                    newstack = append(newstack, backPair{newstart, comsyl})
                }
                index +=1
            }

        }
    }
    return wg
}

func (wg *WordGraph)Print() {
    for _, wds := range wg.wg {
        for _, hz := range wds {
//            fmt.Println(hz)
            term := wg.dt.GetTermById(hz.id)
            fmt.Println("[",hz.start, ",", hz.end, "]", hz.id, term, hz.flag)
        }
    }
}
