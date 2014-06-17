package dict

import (
    "fmt"
    "log"
    "math"
    "strconv"
)

const (
    LimitId = 65535
    InvalidWeight = 65535
    FixedStrategyStartID = 50000
    NumLimitStrategyFix = 50
    NumMaxStrategyFreq = 100
    NumMaxStrategyMutal = 100
    MaxWeightChoice = 4
)

const (
    CommonFlag = 0
    ApproximateFlag = 1
    NotExistFlag = 2
    DanziPenaltyFlag = 4
)

const (
    PathPenaltyMin = 0
    PathPenaltyMax = 32
    DanziPenaltyMin = 0
    DanziPenaltyMax = 32
)

type Parameters struct {
    time            int         "0~10000"
    PathPenalty     float64     "0~32, 100次"
    DanziPenalty    float64     "0~32, 100次"
}

type SysDict struct {
    pd  *pyDict
    td  *termDict 
    gd  *gramDict

    sourcegdpath    string
    sourcegdname    string
    currgram        string
    all int64
    modifynums  int

    p           Parameters
}

//init dict from pyfile, termfile, gramfile
func (s *SysDict)InitDict(pyfile, termfile, dest, gram, curr string) {
    if s.pd == nil {
        s.pd = NewPyDict()
    }
    e :=s.pd.Load(pyfile)
    fmt.Println("load pyfile:", pyfile)
    if e!= nil {
        fmt.Println(pyfile, e)
    }

    if s.td == nil {
        s.td = NewTermDict()
    }

    fmt.Println("load termfile:", termfile)
    e = s.td.Load(termfile)
    if e!= nil {
        fmt.Println(termfile,e)
    }
    s.all = s.td.GetFreqs()
    fmt.Println("all freqs is: ", s.all)

    if s.gd == nil {
        s.gd = NewGramDict()
    }

    fullgdfile := dest +"/"+curr
    fmt.Println("load gram file:", fullgdfile)
    e = s.gd.Load(fullgdfile)
    if e!= nil {
        fmt.Println(curr, e)
    }

    if dest[len(dest)-1] != '/' {
        dest = dest + "/"
    }

    s.sourcegdpath = dest
    s.sourcegdname = gram
    s.currgram = curr
}

//found words by pinyin
func (s *SysDict)GetWordsByPy(py string)[]string {
    ids :=s.pd.GetTermsIdByPy(py)
    re := make([]string, 0)
    for _, id:= range ids {
        term := s.td.GetTermById(id)
        if term == "" {
            continue
        }
        re = append(re, term)
    }
    if len(re) == 0 {
        return nil
    }
    return re
}

//found word by id
func (s *SysDict)GetTermById(id int) string {
    return s.td.GetTermById(id)
}

//found word ids by pinyin
func (s *SysDict)GetTermsIdByPy(py string)[]int {
    return s.pd.GetTermsIdByPy(py)
}

//get ids from term
func (s *SysDict)GetIdByTerm(term, py string) []int {
    if py == "" {
        return s.td.GetIdByTerm(term)
    }

    ids := s.pd.GetTermsIdByPy(py)
    for _, id := range ids {
        if term == s.td.GetTermById(id) {
            return []int{id}
        }
    }
    return nil
}

//get id rate
func (s *SysDict)GetRateById(id int) float64 {
    return s.td.GetTermRateById(id)
}

func (s *SysDict)GetWeight(id1, id2, choose int) (float64, int) {
    switch(choose) {
        case 0:
            return s.GetTransWeight(id1,id2)
        case 1:
            return s.GetTransWeight2(id1,id2)
        case 2:
            a, c :=s.GetPenaltyWeight(id1,id2)
            b, d :=s.GetTransWeight(id1,id2)
            return a+b, c|d
        case 3:
            a, c :=s.GetPenaltyWeight(id1,id2)
            b, d :=s.GetTransWeight2(id1,id2)
            return a+b, c|d
    }
    return 0, 0
}

//get transfer weight
func (s *SysDict)GetTransWeight(id1, id2 int) (float64,int) {
    num_12 := int64(s.gd.GetTrans(id1,id2))
    num1 := int64(s.td.GetTermFreqById(id1))
    flag :=0 
    if num_12 == 0 {
        num_12 = s.GetCoocNumsBySmooth(id1,id2)
        flag = 1
    }

    if num_12 > num1 {
        fmt.Println("error gramdict and sysdict")
        panic(str(id1)+"\t"+str(id2)+"\t"+str(num_12)+"\t"+str(num1))
        num1 += num_12
        exit(-1)
    }

    if num1 == 0 {
        return 32.0,2
    }
    if num_12 == 0 {
        return 32.0,2
    }

    rate := -1.0* math.Log(float64(num_12)/float64(num1))
    return rate, flag

//    return float64(s.gd.GetTrans(id1,id2)), 0
}

func (s *SysDict)GetPenaltyParams() (float64,float64) {
    return s.p.PathPenalty, s.p.DanziPenalty
}

func (s *SysDict)GetPenaltyWeight(id1,id2 int) (float64, int) {
    ppath, pdanzi := s.GetPenaltyParams()
    idif1 := s.td.GetTermInfoById(id1)
    idif2 := s.td.GetTermInfoById(id2)
    if len(idif1.hz) < 4 && len(idif2.hz) < 4 {
        return ppath + pdanzi, DanziPenaltyFlag
    }
    return ppath,0
}

func (s *SysDict)GetTransWeight2(id1, id2 int)(float64,int) {
    num_12 := int64(s.gd.GetTrans(id1,id2))
    if num_12 == 0 {
        return s.GetCoocRateByInfo(id1), 1
    }
    num1 := int64(s.td.GetTermFreqById(id1))
    if num1 ==0 {
        return 32.0, 2
    }

    rate := -1.0* math.Log(float64(num_12)/float64(num1))
    return rate, 0
}


func (s *SysDict)GetCoocNumsBySmooth(id1, id2 int) int64 {
    n1 := s.td.GetTermFreqById(id1)
    n2 := s.td.GetTermFreqById(id2)

    return int64(n1)*int64(n2)/s.all
}


func (s *SysDict)GetCoocRateByInfo(id1 int) float64{
    freq := s.td.GetTermFreqById(id1)
    ci := s.gd.GetTermCoocInfo(id1)

    wordnum := s.td.GetWordsNum()
    if ci == nil {
        return -1.0*math.Log(float64(1.0)/float64(wordnum))
    }
    fr := float64(freq-ci.freq)/float64(wordnum-ci.num) /float64(freq)
    return  -1.0*math.Log(fr)
}


func IsAscSpace(t rune)bool {
    if t == ' ' || t== '\t' {
        return true
    }
    return false
}

//通过传入一个mword_dict: key word_id; val:not used
//return $1:map:key old_word_id; val  new_word_id
//调整策略为新词条从某个ID_end 向前交换
func (s *SysDict)GetModifyDict(mword map[int]int, modifynums int) (map[int]int, int)  {

    if modifynums < 100 {
        return s.ModifyByFreq(mword, modifynums)
    }else {
        return s.ModifyByMutal(mword, modifynums)
    }
}

func (s *SysDict)ModifyByFreq(mword map[int]int, modifynums int) (map[int]int, int) {
    ln := len(mword)
    if ln < NumLimitStrategyFix  {
        return s.modifyByFix(mword, modifynums)
    }

    return s.modifyByFreq(mword, modifynums)
}

//Now state to be defined
func (s *SysDict)ModifyByMutal(mword map[int]int, modifynums int) (map[int]int, int) {
    return nil, modifynums+1
}

//In this method, exchange ids in $1 and [startid,startid+ln]
func (s *SysDict)modifyByFix(mword map[int]int, modifynums int) (map[int]int, int) {
    ln := len(mword)
//    if ln > NumLimitStrategyFix {
//        panic("error mofify method")
//    }
    startid := FixedStrategyStartID + modifynums * NumLimitStrategyFix
    old2new := make(map[int]int)

    i := 0
    for k, _ := range mword {
        old2new[k] = startid + i
        old2new[startid + i] = k
        i += 1
    }

    if i != ln {
        panic("error nums")
    }

    return old2new, modifynums+1
}

func (s *SysDict)modifyThird(mword map[int]int, modifynums int) (map[int]int, int) {
    ln := len(mword)
//    if ln > NumLimitStrategyFix {
//        panic("error mofify method")
//    }
    startid := FixedStrategyStartID + modifynums * NumLimitStrategyFix
    old2new := make(map[int]int)

    i := 0
    for k, v := range mword {
        old2new[k] = startid + i
        if v == -1 {
            old2new[startid + i] = k
        }else{
            old2new[startid + i ] = v
        }

        i += 1
    }

    if i != ln {
        panic("error nums")
    }

    return old2new, modifynums+1
}
//从前FixedStrategyStartID中，依次比较，替换freq小的ID，对被替换的ID和未替换的ID，调用
//modifyByFix()
func (s *SysDict)modifyByFreq(mword map[int]int, modifynums int) (map[int]int, int) {
//    ln := len(mword)

    old2new := make(map[int]int )   //store old to new ids
    
    tempdict := make(map[int]int)    //保存前FixedStrategyStartID
    for k, v:= range mword {
        freqk := s.td.GetTermFreqById(k)
        flag := 0

        id := 1
        for id = 1;id < FixedStrategyStartID; id ++ {
            freqid := s.td.GetTermFreqById(id)
            if freqk > freqid {
                flag = 1
                break
            }
        }

        if flag == 1 {
            old2new[k] = id
            tempdict[id] = k
        }else {
            tempdict[k] = -1
        }
    }

    n2, nnums := s.modifyThird(tempdict, modifynums)
    for k, v := range n2 {
        old2new[k] = v
    }

    return old2new, nnums
}



func (s *SysDict)ModifyTermDict(old2new map[int]int, dest string, cc chan int) {
    s.td = s.td.Modify(old2new)
    if cc != nil {
        cc <- 1
    }
    s.td.Save(dest)
}

func (s *SysDict)ModifyPinyinDict(old2new map[int]int, dest string, cc chan int) {
    s.pd.Modify(old2new)
    if cc != nil {
        cc <- 1
    }
//    s.py.Save(dest)
}

func (s *SysDict)ModifyGramDict(old2new map[int]int, source,dest string, cc chan int) {
    ModifyGramDict(old2new, source, dest)   //因为source_gram要占用极大的内存，所以这里要对文件操作，并重新载入gram
    s.gd = NewGramDict()
    s.gd.Load(dest)
    if cc != nil {
        cc <- 1
    }
}
func un(s string) {
    log.Println("exit: ", s)
}
func trace(s string)string {
    log.SetFlags(log.Ltime)
    log.Println("start: ", s)
    return s
}


func (s *SysDict)ModifyDict(md map[int]int, modifynums int ) {
    defer un(trace(fmt.Sprintf("modify dict %dth", modifynums)))
    nd, _ := s.GetModifyDict(md, modifynums)
    s.modifynums = modifynums+1
    cc := make(chan int, 3)
    go s.ModifyTermDict(nd, s.sourcegdpath+"term."+strconv.Itoa(modifynums)+".txt", cc)
    go s.ModifyPinyinDict(nd, "", cc)
    go s.ModifyGramDict(nd, s.sourcegdpath+s.sourcegdname, s.sourcegdpath+"new."+strconv.Itoa(modifynums)+".txt", cc)
    for i:=0;i<3;i++ {
        <- cc
    }
}

func (s *SysDict)ModifyPenaltyParams(modstart int) {
    s.p.time = modstart
    s.p.PathPenalty = float64(s.p.time/50)* float64(PathPenaltyMax-PathPenaltyMin)/50.0 + float64(PathPenaltyMin)
    s.p.DanziPenalty = float64(s.p.time%50) * float64(DanziPenaltyMax-DanziPenaltyMin)/50.0 + float64(DanziPenaltyMin)
}

func (s *SysDict)EraseNoUsedGram(id int) {
//    cp, e := os.Getwd() 
//    if e {
//        return
//    }
//
//    os.Chdir(s.sourcegdpath)
//
//    fileName := s.sourcegdpath+ strconv.Itoa(id) + ".txt"
}
