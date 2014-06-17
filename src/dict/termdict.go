package dict

import (
    //"readline"
    "os"
    "io"
    "bytes"
    "fmt"
    "strconv"
    "common"
)

type termDict struct {
    freqs   int64
    id  map[int]*IdInfo     "iddict, key:id, val:IdInfo"
    hz  map[string][]int    "hzdict, key:term, val:ids"
    wordsnum    int
}

type IdInfo struct {
    id  int
    hz  string
    freq    int
    order   int
    rate    float64
    tag     []int
}


func (t *termDict)Load(file string)error {
    f, e:= os.Open(file)
    if e != nil {
        return e
    }
    defer f.Close()


//    bf := make([]byte, 1024)
//    var err error
    linenum := 0
    var allfreq int64 = 0
    lr := common.NewLineReader(f, 1024)
    for {
//        n , e:= readline.ReadLine(f, bf)
        bf , e:= lr.ReadLine()
        if e != nil && e!= io.EOF {
            return e
        }
//        if n == 0 {
//            if e == io.EOF {
//                break
//            }
//            continue
//        }

        if bf == nil && e == nil {
            continue
        }
        if bf == nil && e == io.EOF {
            break
        }

//        st := bytes.FieldsFunc(bf[:n], IsAscSpace)
        st := bytes.FieldsFunc(bf, IsAscSpace)
        if len(st) < 3 {
            return common.StrError(fmt.Sprintf("file %s format error with line %d", file, linenum+1))
        }

        hz := string(st[0])
        freq,_ := strconv.Atoi(string(st[1]))
        order,_ := strconv.Atoi(string(st[2]))
        t.id[linenum] = &IdInfo{id:linenum, hz:hz, freq:freq, order:order}
        allfreq += int64(freq)
        linenum +=1
        
        if e == io.EOF {
            break
        }
    }
    t.freqs = allfreq
    t.wordsnum = linenum
    return nil
}

func (t *termDict)GetWordsNum() int {
    return t.wordsnum
}

//found id, return val; else return nil
func (t *termDict)GetTermInfoById(id int) *IdInfo{
    in, ok := t.id[id]
    if ok {
        return in
    }
    return nil
}

//found id, return hz; else return ""
func (t *termDict)GetTermById(id int) string {
    in, ok:= t.id[id]
    if ok {
        return in.hz
    }
    return ""
}

//found term, return ids; else return nil
func (t *termDict)GetIdByTerm(term string) []int {
    h, ok := t.hz[term]
    if ok {
        return h
    }
    return nil
}

//found id, return freq; else return 0
func (t *termDict)GetTermFreqById(id int) int {
    in, ok := t.id[id]
    if ok {
        return in.freq
    }
    return 0
}

//found id, return rate; else return 0
func (t *termDict)GetTermRateById(id int) float64{
    in, ok := t.id[id]
    if ok {
        return in.rate
    }
    return 0
}

func (t *termDict)GetFreqs() int64 {
    return t.freqs
}

func (t *termDict)Save(file string) {
    f, err := os.Create(file)
    if err!= nil {
        fmt.Println(err)
    }
    ln := len(t.id)
    for i:=0;i <ln; i++ {
        v, ok := t.id[i]
        if ok {
            f.WriteString(v.hz)
            f.WriteString("\t")
            f.WriteString(strconv.Itoa(v.freq))
            f.WriteString("\n")
        }
    }
}

func NewTermDict() *termDict {
    td := new(termDict)
    td.id = make(map[int]*IdInfo)
    td.hz = make(map[string][]int)
    return td
}

func (t *termDict)Modify(old2new map[int]int)(nt *termDict) {
    nt = NewTermDict()
    nt.hz = t.hz
    nt.freqs= t.freqs

    for k, v := range t.id{  //k old id, v old idinfo
        nid, ok := old2new[k]   //id to be trans
        if !ok {    //not modify, copy
            nt.id[k] = v
        }else { //modify 
            v.id = nid
            nt.id[nid] = v
        }
    }

    for k, v:= range t.hz {
        for i, id := range v {
            nid, ok := old2new[id]
            if ok {
                t.hz[k][i] = nid
            }
        }
    }

    return nt
}
