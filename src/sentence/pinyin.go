package sentence

import (
    "dict"
    "strings"
    "fmt"
    "time"
//    "log"
)


type Pinyin struct {
    py      string      "input pinyin"
    pyquotalist []string
    exp     string
    hz      []string    "expected hz segment"
    pylist  []string    "expected hz's py, with hz, can get expected from sd"
    sg      *SylGraph       //syl graph
    wg      *WordGraph      //word graph
    sd      *dict.SysDict   //dict userd
    tw      *WeightGraph    //transfer weight graph

    result_hz          [][]string      //result path
    result          [][]int         //continue all ids
    result_trans    []float64   //curr result trans
    bestweight      float64         //best weight value
//    expected_word_ids   []pairInt //expected hz in word graph
    expected            []int        //contain all expected ids and start and end '$'
    expected_trans      []float64   //expected hz transfer weight

    mwords  map[int]int
    sxflag  int
    foundflag  int

    id  int
    choice int
}

func NewPinyin(py string, st *dict.SysDict) *Pinyin{
    return &Pinyin{py:py, sd:st}
}

func (py *Pinyin)SetID(id int) {
    py.id = id
}

func (py *Pinyin)SetHzAndPy(hz, pylist []string) {
    py.hz = make([]string, len(hz))
    py.pylist = make([]string, len(pylist))
    copy(py.hz, hz)
    copy(py.pylist, pylist)

    py.expected = make([]int, len(hz)+2)
    py.expected[0] = 0
    for i:=0; i<len(hz);i++ {
        nindex := py.sd.GetIdByTerm(hz[i], pylist[i])
//        fmt.Println(nindex)
        if nindex != nil {
            py.expected[i+1] = nindex[0]
        }
    }
    py.expected[len(hz)+1] = 0
//    fmt.Println(py.expected)
}

func (py *Pinyin)ResetHzAndPy() {
    py.expected = make([]int, len(py.hz)+2)
    py.expected[0] = 0

    for i:=0; i<len(py.hz);i++ {
        nindex := py.sd.GetIdByTerm(py.hz[i], py.pylist[i])
//        fmt.Println(nindex)
        if nindex != nil {
            py.expected[i+1] = nindex[0]
        }
    }
    py.expected[len(py.hz)+1] = 0
}

//make pinyin syllable graph
func (py *Pinyin)BuildSylGraph() {
    py.sg = MakeSylGraph(py.py)
}

func (py *Pinyin)PrintSylGraph(){
    py.sg.Print()
}

//make pinyin word graph, not contain start '$' and end '$'
func (py *Pinyin)BuildWordGraph() {
    py.wg = MakeWordGraph(py.sg, py.sd, py.py)
}

func (py *Pinyin)PrintWordGraph() {
    py.wg.Print()
}

//make lattice for vitebi, add start '$' and end '$' in wg
func (py *Pinyin)BuildLattice() {
    py.wg = BuildLattice(py.wg)
}

//set sysdict used in py
func (py *Pinyin)SetDict(t *dict.SysDict) {
    py.sd = t
}

//find expected_trans in py.tw
func (py *Pinyin)FoundExpectedInWeightGraph() {
    if py.tw == nil {
        return
    }
    if py.expected == nil {
        return
    }

    ln := len(py.expected)
    py.expected_trans = make([]float64,0, ln)

    for i:=1; i<ln; i++ {
        toid := py.expected[i]
        fromid := py.expected[i-1]
        tw,_ := py.sd.GetWeight(fromid, toid, py.choice)
        py.expected_trans = append(py.expected_trans, tw)
    }

    return
}

//use vitebi algorithmic to get paths, parameters $1 defined path nums 
func (py *Pinyin)GetResult(n int, choose int ) [][]*path {
    py.choice = choose
    py.tw,_ = py.wg.Vitebi(choose)
//    py.tw.Print()

    if py.tw == nil {
        panic("error no tw")
    }
    re, ce :=py.tw.GetPath(n)
    if re == nil {
        panic("no result")
    }
    py.bestweight = ce[0]

//    for _, r := range re {
//        py.PrintPath(r)
//    }

    py.result = make([][]int, n)
    py.result_hz = make([][]string, n)
    for i, r := range re {
        py.result[i] =  make([]int, 0, len(r))
        py.result_hz[i] = make([]string, 0, len(r))
        for j, p := range r {
            term := py.sd.GetTermById(p.idfrom)
            py.result[i] = append(py.result[i], p.idfrom)
            py.result_hz[i] = append(py.result_hz[i], term)

            if j + 1 == len(r) {
                term = py.sd.GetTermById(p.idto)
                py.result[i] = append(py.result[i], p.idto)
                py.result_hz[i] = append(py.result_hz[i], term)
            }
        }
    }

    py.result_trans = make([]float64, len(re[0]))
    for i, p := range re[0] {
        py.result_trans[i] = p.weight
    }

    return re
}

func (py *Pinyin)PrintPath(rp []*path) {
    for _, t := range rp {
        ln := t.left
        rn := t.right
//        if rn.end == -1 {
//            continue
//        }
        fromterm := py.sd.GetTermById(ln.id)
        toterm := py.sd.GetTermById(rn.id)
        fmt.Println("[",ln.start,",", rn.start, "]", ln.id,":",fromterm, "->", rn.id, ":",toterm, t.weight)
    }
}

//diff results and expected, store diff in py.mwords. 
//If $1 is true, if expected not shouxuan, static differs; else
//static if expected is not found
func (py *Pinyin)DiffResult(sx_flag bool) {
    py.foundflag = 0
    py.sxflag = 0
    first_same := -1
    for i, re := range py.result {
        if len(re) < 2 || len(py.result_hz[i]) < 2 {
            continue
        }

        hz := strings.Join(py.result_hz[i][1:len(py.result_hz[i])-1], "")
        fmt.Println(py.exp, py.expected, re, hz,py.result_hz[i] )

        if hz != py.exp {
            continue
        }
        if first_same == -1 {
            first_same = i 
        }else {         //already found
            continue
        }

        change_expected_flag := 0
        if len(py.expected) != len(re) {
            change_expected_flag = 1
        }else {
            for  j:=0; j<len(re);j++ {
                if re[j] != py.expected[j] {
                    change_expected_flag = 1
                    break
                }
            }
        }

        if change_expected_flag == 1 {
            py.expected = make([]int, len(re))
            copy(py.expected,re) 
            //重置hz和pylist
            py.hz = make([]string, len(py.result_hz[i])-2)
            py.pylist = make([]string, 0, len(py.result_hz[i])-2)
            copy(py.hz, py.result_hz[i][1:len(py.result_hz[i])-1])
            j := 0
            for _, t :=range py.hz {
                lt := len(t)
                npy := strings.Join(py.pyquotalist[j:j+(lt+1)/2], "'")
                py.pylist = append(py.pylist, npy)
                j += (lt+1)/2
            }
//            fmt.Println(py.hz, len(py.hz))
//            fmt.Println(py.pyquotalist,len(py.pyquotalist), py.pylist)
        }

        py.foundflag = 1

        if i == 0 {
            py.sxflag= 1
        }
    }

    if sx_flag && py.sxflag == 0 {
        py.mwords = make(map[int]int)

        for _, id := range py.expected {
            if id > dict.LimitId{
                py.mwords[id] += 1
            }
        }
        return
    }

    if py.foundflag== 0 {
        py.mwords = make(map[int]int)

        for _, id := range py.expected {
            if id > dict.LimitId{
                py.mwords[id] += 1
            }
        }
    }

    return
}

//Get words which to modify ids
func(py *Pinyin)GetModifiedWords() map[int]int {
    return py.mwords
}

func (py *Pinyin)FreeMemory() {
    py.sg = nil
    py.wg = nil
    py.tw = nil
    py.result = nil
    py.result_trans = nil
}

func (py *Pinyin)Process(choice int, cc chan <- int, dc chan <- map[int]int, buildflag bool) {


    start := time.Now()
    py.BuildSylGraph()
    py.BuildWordGraph()
    py.BuildLattice()

    if (buildflag) {
        py.ResetHzAndPy()
    }

    py.GetResult(3, choice )
//    pa :=py.GetResult(3, choice )
//    for _, p := range pa {
////        fmt.Println(p)
////        log.Println()
//        py.PrintPath(p)
//    }
    py.DiffResult(true)
    py.FoundExpectedInWeightGraph()
    end := time.Now()
    delta := end.Sub(start)
    fmt.Println("expected trans are:", py.expected_trans, "took time of:%s\n", delta)
    py.FreeMemory()
    dc <-py.mwords
    cc <- py.id
}

