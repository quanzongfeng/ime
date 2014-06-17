package sentence

import (
    "syllable"
    "fmt"
)

type SylGraph struct {
    sg [][]*syllable.PyIdUnit
}

func MakeSylGraph(py string) (re *SylGraph) {
    ln := len(py)
    re = new(SylGraph)
    re.sg = make([][]*syllable.PyIdUnit, 0, ln)

    num := 0
    for i:=0; i <ln ; i++ {
        ts := syllable.SegPy(py, i)
        if ts == nil {
            ts = make([]*syllable.PyIdUnit,0)
        }else {
            num += 1
        }
        re.sg = append(re.sg, ts)
    }

    if num == 0 {
        return nil
    }
    re = re.eraseRedundance()
    return re
}

func (sg *SylGraph)Build() {
}

func (sg *SylGraph)Print() {
    for i, line := range sg.sg {
        for j, st := range line {
            fmt.Println("[",i,",",j,"]",syllable.GetSylById(st.GetId()),st.GetStart(), st.GetEnd(), st.GetFlag())
        }
    }
}

//去除冗余，保留音节，以最短音节为步长，进行下一个筛选
//只支持全拼音节
func (sg *SylGraph)eraseRedundance() (ne *SylGraph) {
    ln := len(sg.sg)

    ne = new(SylGraph)
    ne.sg = make([][]*syllable.PyIdUnit, ln, ln)

    for i:=0; i< ln; i++{
        newline := make([]*syllable.PyIdUnit, 0)
        flag := 0
        for _, st := range sg.sg[i] {
            if st.IsSyllable() {
                newline = append(newline, st) //only save syllable
                flag +=1
            }
        }

        ne.sg[i] = newline
    }

    endgraph := make([]int, ln)
    endgraph[0] = 1
    for i:=0; i< ln ; i++ {
        newline := make([]*syllable.PyIdUnit,0, len(ne.sg[i]))
        for _, st := range ne.sg[i] {
            end := st.GetEnd()
            if end == ln {
                newline = append(newline, st)
                continue
            }
            if len(ne.sg[end]) != 0 { //no syllable start at end, erase
                newline = append(newline, st)
                endgraph[end] = 1
            }
        }
        ne.sg[i] = newline
    }

    for i:=0; i< ln ;i++ {
        if endgraph[i] == 0{
            ne.sg[i] = make([]*syllable.PyIdUnit,0)
        }
    }


    return 
}



