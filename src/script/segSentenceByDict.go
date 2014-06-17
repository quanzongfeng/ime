package main

import (
    "os"
    "io"
    "common"
    "bytes"
    "fmt"
    "strconv"
)

func LoadDict(file string) (map[string]int, error) {
    f, e := os.Open(file)
    if e!= nil {
        return nil, e
    }
    defer f.Close()

    dt := make(map[string]int)
    lr := common.NewLineReader(f, 1024)
    linenum := 0
    for {
        linenum += 1
        bf, e := lr.ReadLine()
        if e!= nil && e!= io.EOF {
            return nil, e
        }
        if bf == nil && e== nil {
            continue
        }
        if bf == nil && e == io.EOF {
            break
        }

        st := bytes.FieldsFunc(bf, common.IsAscSpace)
        if len(st) < 3 {
            return  nil, common.StrError(fmt.Sprintf("file %s format error with line %d", file, linenum))
        }

        hz := string(st[0])
        freq, _ := strconv.Atoi(string(st[1]))
        
        t, ok:= dt[hz]
        if ok {
            if t< freq {
                dt[hz] = freq
            }
        }else {
            dt[hz] = freq
        }
        if e== io.EOF {
            break
        }
    }

    return dt, nil
}

func  segForthByDict(line string, dt map[string]int)([]string, error) {
//    defer func(){fmt.Println("segForthByDict for ",line, "succeed")}()
    indexs, e:= common.GetGbkHzIndexsList([]byte(line))
    if e!= nil {
        return nil, e
    }

    ln := len(indexs)

    re := make([]string, 0, 4)

    i := 0
    for i < ln-1 {
        start := indexs[i]
        j := i
        for j=ln-1; j>i; j-- {
            end := indexs[j]
            sr := line[start:end] 
            _, ok := dt[sr]
            if ok {
                re = append(re, sr)
                break
            }
        }
        if j == i{
            return nil ,common.StrError("not found words in dict")
        }
        i = j
    }

    return re, nil
}

func segBackByDict(line string, dt map[string]int) ([]string, error) {
//    defer func(){fmt.Println("segBackByDict for ",line, "succeed")}()
    indexs, e:= common.GetGbkHzIndexsList([]byte(line))
    if e!= nil {
        return nil, e
    }

    ln := len(indexs)

    re := make([]string, 0, 4)

    i := ln-1
    for i >= 0 {
        end := indexs[i]
        j := i
        for j= 0; j < i ;j++ {
            start := indexs[j]
            sr := line[start:end]
            _, ok := dt[sr]
            if ok {
                re = append(re, sr)
                break
            }
        }
        if j == i {
            return nil, common.StrError("not found words in dict")
        }

        i = j

    }

    lre := len(re)
    for i:=0; i< lre/2 ;i++ {
        re[i], re[lre-1-i] = re[lre-1-i], re[i]
    }
    return re, nil
}

func segSentence(line string, dt map[string]int) ([]string, error) {
    fst, e := segForthByDict(line, dt)
    bst, e2 := segBackByDict(line, dt)
    if e != nil && e2 != nil {  //都有错误发生
        return nil, e
    }else if e != nil {         //前向切分错误，则返回后向切分
        return bst, nil
    }else if e2 != nil {        //后向切分错误，则返回前向切分
        return fst, nil
    }


    if len(bst) == len(fst) { //前后向切分长度相同
        flag := 0
        for i, t:= range bst {
            if t != fst[i] {
                flag = 1
                break
            }
        }
        if flag == 0 {      //前后切分数据相同，返回
            return fst, nil
        }
    }

    //前后向切分均成功，但切分方法不同，则根据频率进行计算
    freqforth := 1
    for _, t:= range fst {
        ft, ok := dt[t]
        if !ok {
            panic( "segment forth error")
        }
        freqforth *= ft
    }

    freqback := 1
    for _, t := range bst {
        ft, ok := dt[t]
        if !ok {
            panic("segment back error")
        }
        freqback *= ft
    }

    if freqforth > freqback {
        return fst, nil
    }
    return bst, nil
}


func main() {
    if len(os.Args) < 2 {
        fmt.Println("seg $1 $2")
        return 
    }
    dtfile := os.Args[1]
    sourcefile := os.Args[2]

    dt, e := LoadDict(dtfile)
    if e!= nil {
        fmt.Println(e)
        return
    }
    f, e := os.Open(sourcefile)
    if e!= nil {
        fmt.Println(e)
        return
    }
    defer f.Close()

    lr := common.NewLineReader(f, 1024)
    linenum := 0
    for {
        linenum += 1
        bf, e := lr.ReadLine()
        if e!= nil && e!= io.EOF {
            fmt.Println(e)
            return 
        }
        if bf == nil && e== nil {
            continue
        }
        if bf == nil && e == io.EOF {
            break
        }

        line := bf
        a, e:= segSentence(string(line), dt)
//        fmt.Println(line, a)
        if e != nil {
            fmt.Println(e)
            continue
        }
        os.Stdout.Write(line)
        os.Stdout.Write([]byte("\t"))
        la := strconv.Itoa(len(a))
        os.Stdout.Write([]byte(la))
        os.Stdout.Write([]byte("\t"))
        for i, t := range a{
            os.Stdout.Write([]byte(t))
            if i != len(a) {
                os.Stdout.Write([]byte("\t"))
            }
        }
        os.Stdout.Write([]byte("\n"))
    }

    return
}

