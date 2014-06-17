package dict

import (
    "io"
    "os"
    "bytes"
    //"readline"
    "common"
    "fmt"
    "strconv"
)

type pyDict struct {
    p map[string][]int
}


func (p *pyDict)Load(file string) error {
    f, e := os.Open(file)
    if e!= nil {
        return e
    }
    defer f.Close()

//    var err error
    linenum := 0
//    bf := make([]byte, 10240)

    lr := common.NewLineReader(f, 10240)

    for {
//        n, e:= readline.ReadLine(f, bf)
        bf, e:= lr.ReadLine()
        if e!= nil && e!= io.EOF {
            return e
        }

        if bf == nil && e== nil{
            continue
        }
        if bf == nil && e == io.EOF {
            break
        }

//        if n == 0 {
//            if e == io.EOF {
//                break
//            }
//        }

        linenum +=1
//        st := bytes.FieldsFunc(bf[:n], IsAscSpace)
        st := bytes.FieldsFunc(bf, IsAscSpace)

        if len(st) < 4 {
            return common.StrError(fmt.Sprintf("file %s format error with line %d", file,linenum))
        }

        py := string(st[0])
//      need to ignore "'"? no. cause tian:ti'an not same
//        py = strings.Replace(py, "'", "", -1)

        idterms := st[2:]
        ids := make([]int, 0, 1)
        for i:=0; i<len(idterms); i+=2 {
            id,_ := strconv.Atoi(string(idterms[i]))
            ids = append(ids, id)
        }
        p.p[py] = ids

        if e == io.EOF {
            break
        }
    }
    return nil
}

func (p *pyDict)GetTermsIdByPy(py string)[]int {
    in, ok := p.p[py]
    if ok {
        return in
    }
    return nil
}

func (p *pyDict)Modify(old2new map[int]int){
    for _, v := range p.p { //k string, v ids
        for i, id := range v {
            if nid, ok := old2new[id]; ok { //id in old2new, to be trans
                v[i] = nid
            }
        }
    }
}

func (p *pyDict)Save(file string) {
}

func NewPyDict() *pyDict {
    tn := make(map[string][]int)
    return &pyDict{tn}
}


