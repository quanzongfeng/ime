package sentence

import (
    "io"
    "bytes"
    "fmt"
    "readline"
    "os"
    "dict"
    "common"
    "strconv"
//    "log"
    "time"
    "syllable"
)

const (
    ModifyNumLimits = 20
    ModifyNumRate = 0.1
    ModifyFreqLimits = 0    //for test =0
    MaxModifyDictNums = 200
)

type PyGroup struct {
    pys     []*Pinyin
    mwords  map[int]int
    sd      *dict.SysDict

    moddictnums         int     //标识当前mod词典次数
    bestmoddictnum      int     //标识当前最优词典的id
    modparamnums        int
    parameters      []float64   //当前使用的参数
    bestparams       []float64   //最优首选率下对应的参数
    bestchoose          int
    bestmoddict         int

    shouxuanrate    float64 //当前词典和参数下的首选率
    shangpingrate   float64 //当前词典和参数下的上屏率
    shangpingnum    int     //当前词典和参数下上屏次数
    shouxuannum     int     //当前词典和参数下首选次数

    bestsxrate      float64 //到当前测试为止的最优首选率
    bestsprate      float64 //到当前测试为止最优首选率下的上屏率

    dataChan    chan map[int]int    //该channel用于Pinyin中传送数据
    controlChan chan int            //该channel用于传送控制命令
    endChan     chan int
}


func (pg *PyGroup)LoadTest(file string)( error){
    f, e := os.Open(file)
    if e!= nil {
        return  e
    }
    defer f.Close()

    bf := make([]byte, 10240)
    linenum :=0

    for {
        n, e := readline.ReadLine(f, bf)
        if e!= nil && e!= io.EOF {
            return  e
        }
        if n ==0 {
            if e== io.EOF {
                break
            }
            continue
        }

        st := bytes.FieldsFunc(bf[:n], common.IsAscSpace)
        if len(st) < 3 {
            return  common.StrError(fmt.Sprintf("file %s format error with line %d", file, linenum+1))
        }

        //        py := string(st[0])
        pynoquota := common.FilterSingleQuota(st[0])
        hz := string(st[1])
        num, _ := strconv.Atoi(string(st[2]))
        var hzl []string = nil
        var pyl []string = nil
        var pyq []string = nil

        pyt := bytes.FieldsFunc(st[0], common.IsAscQuota)
        pyq = make([]string, 0,len(pyt))
        for _, t := range pyt {
            pyq = append(pyq, string(t))
        }

        if (num != 0) {
            hzlist := make([]string, 0)
            pylist := make([]string, 0)
            if len(st) < 4 {
                return common.StrError(fmt.Sprintf("file %s format error with line %d", file, linenum+1))
            }
            fn, _ := strconv.Atoi(string(st[3]))
            if len(st) < 4+fn {
                return common.StrError(fmt.Sprintf("file %s format error with line %d", file, linenum+1))
            }

            for i:=4; i<4+fn;i++ {
                hzlist = append(hzlist, string(st[i]))
            }
            //            pytstring := make([]string, len(pyt))
            //            for i:=0; i< len(pyt);i++ {
            //                pytstring[i] = string(pyt[i])
            //            }

            j := 0

            //            fmt.Println(string(st[0]))
            //            fmt.Println(len(pytstring), pytstring)
            for _, t := range hzlist {
                lt := len(t)
                npy := bytes.Join(pyt[j: j+(lt+1)/2], []byte{'\''})

                pylist = append(pylist, string(npy))
                j+=(1+lt)/2
            }

            if j != len(pyt) {
                return common.StrError(fmt.Sprintf("file %s format error with line %d", file, linenum+1))
            }

            hzl = hzlist
            pyl = pylist
            //            fmt.Println(hzl)
            //            fmt.Println(pyl)
        }

        npy := NewPinyin(string(pynoquota), pg.sd)
        if hzl != nil {
            npy.SetHzAndPy(hzl, pyl)
        }
        npy.SetID(linenum+1)
        npy.exp = hz
        npy.pyquotalist = pyq
        pg.pys = append(pg.pys, npy)
        linenum +=1
    }
    return nil
}

func NewPyGroup() *PyGroup {
    pg := new(PyGroup)
    pg.pys = make([]*Pinyin, 0, 1)
    pg.controlChan = make(chan int)
    pg.dataChan = make(chan map[int]int)
    pg.endChan = make(chan int)
    return pg
}

func (pg *PyGroup)SetDict(s *dict.SysDict) {
    pg.sd = s
}

func (pg *PyGroup)ProcessControlChanGo() {
    endflag := 0
    controlDict := make(map[int]int)
    for t := range pg.controlChan {
//        fmt.Println(t)
        if t == -1 {
            endflag = 1
        }else {
            _, ok := controlDict[t]
            if !ok {
                controlDict[t] = t
            }else {
                controlDict[t] ^= t
            }
            if controlDict[t] == 0 {
                delete(controlDict, t)
            }
        }

        if len(controlDict)==0 && endflag == 1{
            break
        }
    }

    end := make(map[int]int)
    end[-1] = -1
    pg.dataChan <- end    //send endflag
}

func (pg *PyGroup)ProcessDataChanGo() {
    pg.mwords = make(map[int]int)

    for t := range pg.dataChan {
        if t == nil {
            continue
        }

        k, ok := t[-1] 
        if ok && k==-1 {
            break
        }

        for k, v:= range t {
            _, ok := pg.mwords[k]
            if !ok {
                pg.mwords[k] = v
            }else {
                pg.mwords[k] += v
            }
        }
    }

    pg.endChan <- 1

}

//this func should used after processonce
func (pg *PyGroup)StaticRate() {
    ln := len(pg.pys)
    all_shouxuan := 0
    all_shangping := 0
    for _, t:= range pg.pys {
        all_shangping += t.foundflag
        all_shouxuan  += t.sxflag
    }

    pg.shangpingrate = float64(all_shangping)/float64(ln)
    pg.shouxuanrate = float64(all_shouxuan)/float64(ln)
}

//判断是否需要调整2元关系
func (pg *PyGroup)NeedModify(md map[int]int) bool {
    ln := len(md)
    wordnum := len(pg.mwords)

    if float64(ln) > (float64(ModifyNumLimits)/ModifyNumRate){
        if float64(wordnum)/float64(ln) > ModifyNumRate {
            return true
        }
    }else {
        if wordnum > ModifyNumLimits {
            return true
        }
    }
    return false
}

func (pg *PyGroup)GetModifyWords() map[int]int {
    nm := make(map[int]int)
    for k, v := range pg.mwords {
        if v > ModifyFreqLimits{
            nm[k] = v
        }
    }
    return nm
}

//计算当前2元和参数下的结果
func (pg *PyGroup)ProcessOnce(choice int, rebuild bool) {   
    ln := len(pg.pys)

    go pg.ProcessControlChanGo()
    go pg.ProcessDataChanGo()

    for i:=0; i< ln ;i++ {
        py:= pg.pys[i]
        pg.controlChan <- py.id
        go py.Process(choice, pg.controlChan, pg.dataChan, rebuild)   //次序很重要，先发data，再发control
    }

    pg.controlChan <- -1    //-1 is end flag
    <- pg.endChan        //here means Proceed all
    //getallmodifyedwords
    //calc shouxuan,shangping rate and record
}

func (pg *PyGroup)GetParameters() {
    if pg.parameters==nil {
        pg.parameters = make([]float64,2)
    }
    pp, dp := pg.sd.GetPenaltyParams()
    pg.parameters[0] = pp
    pg.parameters[1] = dp
}

func (pg *PyGroup)StartChoose() {
    choose := 0
    n:= 0
    modifydict := true
    syllable.SegPy("nihao",0)
    for n=0;n <10000; n++ {
        if pg.modparamnums >= 2500 {
            break
        }
        if n == 1{
            break
        }

        fmt.Println("start :", n)
        s_start := time.Now()

        choice := choose % dict.MaxWeightChoice
        pg.ProcessOnce(choice, modifydict)

        pg.StaticRate()
        pg.GetParameters()
        s_end := time.Now()
        delta:=s_end.Sub(s_start)
        fmt.Printf("this routinue use time:%s\n", delta)
        modifydict = false

        if pg.bestsxrate == 0 {
            pg.bestchoose = choice
            pg.bestsxrate = pg.shouxuanrate
            pg.bestsprate = pg.shangpingrate
            pg.bestmoddictnum  = pg.moddictnums
            pg.bestparams = pg.parameters
            pg.bestmoddict = pg.moddictnums
        }else if pg.bestsxrate < pg.shouxuanrate {
            pg.bestchoose = choice
            pg.bestsxrate = pg.shouxuanrate
            pg.bestsprate = pg.shangpingrate
            pg.bestmoddictnum = pg.moddictnums
            pg.bestparams = pg.parameters
            pg.bestmoddict = pg.moddictnums
        }

        td := pg.GetModifyWords()
        fmt.Println("to mod are:", len(td), "\t", td)
        fmt.Println("result: ",pg.bestsxrate, pg.bestsprate,pg.bestchoose, pg.bestparams, pg.bestmoddict)
        fmt.Println(n , choose, choice, pg.moddictnums, pg.modparamnums)
//        if pg.NeedModify(td) && pg.moddictnums < MaxModifyDictNums{
//            pg.sd.ModifyDict(td, pg.moddictnums)
//            pg.moddictnums++ 
//            modifydict = true
//            pg.modparamnums = 0
//            pg.sd.ModifyPenaltyParams(pg.modparamnums)
//            choose = 0
//            continue
//            //            break
//        }
        choose += 1
        if choose % dict.MaxWeightChoice ==0 {
            pg.sd.ModifyPenaltyParams(pg.modparamnums)
            pg.modparamnums +=1
            //
            //        go pg.sd.EraseNoUsedGram(pg.bestmoddictnum)
        }
    }
    fmt.Println("End with bestrate :", pg.bestsxrate, pg.bestsprate,pg.bestchoose, pg.bestparams)
}

