package dict

import (
    "fmt"
    "common"
    "io"
    "os"
    "bytes"
//    "readline"
    "strconv"
//    "readline"
)

type CoocInfo struct {
    freq    int
    num     int
}

type gramDict struct {
    g       map[string]int
    info    map[int]*CoocInfo
}

func (g *gramDict)Load(file string)error {
    f, e := os.Open(file)
    if e != nil {
        return e
    }
    defer f.Close()

//    bf := make([]byte, 10240)
    lr := common.NewLineReader(f, 10240)
    linenum := 0
    for {
//        n, e:= readline.ReadLine(f, bf)
        bf, e:= lr.ReadLine()
        if e!= nil && e!= io.EOF {
            return e
        }
        if bf == nil && e == nil{
            continue
        }
        if bf == nil && e == io.EOF {
            break
        }
//        if n ==0 {
//            if e == io.EOF {
//                break
//            }
//            continue
//        }

//        st := bytes.FieldsFunc(bf[:n], IsAscSpace)
        st := bytes.FieldsFunc(bf, IsAscSpace)
        if len(st) < 3 {
            return common.StrError(fmt.Sprintf("file %s format error with line %d", file, linenum+1))
        }

        key := string(st[0])+" "+string(st[1])
        freq, _:= strconv.Atoi(string(st[2]))

        g.g[key] = freq

        idfrom, _ := strconv.Atoi(string(st[0]))
        _, ok := g.info[idfrom]
        if ok {
            g.info[idfrom].freq += freq
            g.info[idfrom].num += 1
        }else {
            g.info[idfrom] = &CoocInfo{freq:freq, num:1}
        }


        linenum +=1

        if e == io.EOF {
            break
        }
    }
    return nil
}

func (g *gramDict)GetTrans(id1, id2 int)int {
    t := strconv.Itoa(id1)+" "+strconv.Itoa(id2)
    gf, ok := g.g[t]
    if ok {
        return gf
    }
    return 0
}

func (g *gramDict)Save(file string) {
}

func NewGramDict() *gramDict {
    g := new(gramDict)
    g.g = make(map[string]int)
    g.info = make(map[int]*CoocInfo)
    return g
}

func ModifyGramDict(old2new map[int]int, source, dest string) error {
    f, e := os.Open(source)
    if e != nil {
        return e
    }
    defer f.Close()

    g, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer g.Close()

//    bf := make([]byte, 10240)
    linenum :=0
    lr := common.NewLineReader(f, 1024)


    for {
       // n, e:= readline.ReadLine(f, bf)
        bf, e:= lr.ReadLine()
        if e!= nil && e!= io.EOF {
            return e
        }
        if bf == nil && e != io.EOF {
            continue
        } 
        if bf == nil && e == io.EOF {
            break
        }
//        if n ==0 {
//            if e == io.EOF {
//                break
//            }
//            continue
//        }

        st := bytes.FieldsFunc(bf, IsAscSpace)
        if len(st) != 3 {
            return common.StrError(fmt.Sprintf("file %s format error with line %d", source, linenum+1))
        }

        idfrom,_ := strconv.Atoi(string(st[0]))
        idto,_ := strconv.Atoi(string(st[1]))

        idfWrite := string(st[0])
        idtWrite := string(st[1])

        idfromtemp, ok1 := old2new[idfrom]
        if ok1 {
            idfWrite = strconv.Itoa(idfromtemp)
        }

        idtotemp, ok2 := old2new[idto]
        if ok2 {
            idtWrite = strconv.Itoa(idtotemp)
        }

        idf, _ := strconv.Atoi(idfWrite)
        idt, _ := strconv.Atoi(idtWrite)

        if idf > LimitId || idt > LimitId {
            continue
        }

        towrite := idfWrite+" "+idtWrite+"\t"+string(st[2])+"\r\n"
        g.Write([]byte(towrite))
        linenum +=1
    }

    return nil
}

func ChooseGramDict(source, dest string)error {
    f, e := os.Open(source)
    if e != nil {
        return e
    }
    defer f.Close()

    g, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer g.Close()

//    bf := make([]byte, 10240)
    linenum :=0
    lr := common.NewLineReader(f, 1024)

    for {
        //n, e:= readline.ReadLine(f, bf)
        bf, e:= lr.ReadLine()
        if e!= nil && e!= io.EOF {
            return e
        }

        if bf == nil && e != io.EOF {
            continue
        } 
        if bf == nil && e == io.EOF {
            break
        }
//        if n ==0 {
//            if e == io.EOF {
//                break
//            }
//            continue
//        }

        st := bytes.FieldsFunc(bf, IsAscSpace)
        if len(st) != 3 {
            return common.StrError(fmt.Sprintf("file %s format error with line %d", source, linenum+1))
        }

        idfrom,_ := strconv.Atoi(string(st[0]))
        idto,_ := strconv.Atoi(string(st[1]))

        if idfrom >= LimitId || idto >= LimitId {
            continue
        }
        g.Write(bf)
        g.Write([]byte("\n"))
    }
    return nil
}

func (g *gramDict)GetTermCoocInfo(id int) *CoocInfo  {
    t, ok := g.info[id]
    if ok {
        return t
    }
    return nil
}
