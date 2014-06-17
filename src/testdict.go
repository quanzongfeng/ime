package main

import (
    "dict"
    "fmt"
    "sentence"
//    iconv "github.com/qiniu/iconv"
    "common"
    "os"
    "time"
    "runtime/pprof"
)

func main() {
    cpuprofile :="prorecord.txt"
    f , err := os.Create(cpuprofile)
    if err != nil {
        panic("create profile failed")
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    common.SetCpuNums(12)
//    tt := "你好"
//    cd, e := iconv.Open("gbk","utf-8")
//    if e!= nil {
//        fmt.Println(e)
//    }
//    defer cd.Close()
//    rt := cd.ConvString(tt)
//    fmt.Println(tt)
//    fmt.Println(rt)

//    fmt.Println(s.GetWordsByPy("ni'men"))
    s := new(dict.SysDict)
    s.InitDict("source/pinyin.txt","source/termfreq.txt","source/gram/", "gram.txt", "init.txt")
    fmt.Println("laod all dict")
//    time.Sleep(100*time.Second)
//    fmt.Println(s.GetIdByTerm(rt, "ni'hao"))

//    for i:=0;i<10000;i++ {
//        s.ModifyPenaltyParams(i)
//        fmt.Println(s.GetPenaltyParams())
//
//    }

//    py := os.Args[1]
//    spy := sentence.NewPinyin(py, s)
//    spy.BuildSylGraph()
//    spy.PrintSylGraph()
//
//    fmt.Println("try trans py to hz")
//    spy.BuildWordGraph()
//
//    spy.BuildLattice()
//    spy.PrintWordGraph()
//    p:=spy.GetResult(2)
//    for _, t := range p {
//        spy.PrintPath(t)
//    }
    file := os.Args[1]
    pg := sentence.NewPyGroup()
    pg.SetDict(s)
    pg.LoadTest(file)
    fmt.Println("laod all testfile")
    pg.StartChoose()
    time.Sleep(10*time.Second)
}
