package main

import (
    "fmt"
    "os"
    "io"
    "common"
    "time"
)

func main() {
    f , err:= os.Open("temp")
    if err != nil {
        fmt.Println(err.Error())
        return 
    }
    defer f.Close()

    r := common.NewLineReader(f, 1024)

    bt := time.Now()
    for i:=0;;i++{
        lr, e := r.ReadLine()
        if e!= nil && e != io.EOF{
            fmt.Println(e.Error())
            break
        }
        if lr == nil && e == nil {
            continue
        }
        fmt.Println(lr,i, e==nil, e)
//        fmt.Println(string(lr))
        if e == io.EOF {
            break
        }
    }
    at := time.Now()
    dlr := at.Sub(bt)
    fmt.Println("new reader: ",dlr)

//    g, e:= os.Open("term.txt")
//    if e!= nil {
//        fmt.Println(e)
//        return 
//    }
//    defer g.Close()
//    bf := make([]byte, 1024)
//
//    nbt := time.Now()
//    for i:=0;;i++{
//        _, e:=common.ReadLine(g, bf)
//        if e!= nil && e != io.EOF{
//            fmt.Println(e.Error())
//            break
//        }
////        fmt.Printf("n=%p, e=%p\n", &n, e)
////        fmt.Println(bf[:n], e==nil)
////        fmt.Println(string(bf[:n]))
////        fmt.Println(i)
//        if e == io.EOF {
//            break
//        }
//    }
//
//    nat := time.Now()
//    ndlr := nat.Sub(nbt)
//    fmt.Println("old reader: ",ndlr)
//

}


