package dict

import (
    "os"
    "bytes"
    "strconv"
    "readline"
    "stringError"
    "math"
)

const (
    minrate = 1e-10
)

type idInfo struct {
    id      int         "hz id in system"
    hz      string      "hz"
    freq    int         "freq with static"
    rate    float64     "freq/All"
    order   int         "order value"
    tag     []int       "for postag, not used now"
}

type pyDict struct {
    map[string][]int
}
type termDict struct {
    map[int]*idInfo
}

type Dict struct {
    pyID    map[string][]int
    idHz    map[int]*idInfo
    hzId    map[string]int
}

func (t *Dict)Load(py, ih string) (e error) {
    t.pyID, e = t.LoadPyID(py)
    if e != nil {
        return e
    }
    t.hz, e = t.LoadTerms(ih)
    return e
}

func IsAscSpace(t rune)bool {
    if t == ' ' || t== '\t' {
        return true
    }
    return false
}

func (t *Dict)LoadPyID(pyfile string) (map[string][]int, error) {
    f, e:= os.Open(pyfile)
    if e != nil {
        return nil, e
    }

    dt := make(map[string][]string)
    bf := make([]byte, 1024)
    var err error
    linenum := 0
    for {
        n, e := ReadLine(r, bf)
        if e != nil && e!= io.EOF {
            return nil,e
        }
        if n == 0 {
            err = e
            break
        }

        linenum +=1
        st := bytes.FieldsFunc(bf[:n], IsAscSpace)

        if len(st) < 4 {
            return nil, stringError.StringError(fmt.Sprintf("file format error with line %d", linenum))
        }

        py := string(st[0])
        py = strings.Replace(py, "'","",-1)
        idterms := st[2:]
        ids := make([]int,0,1)
        for i:=0; i<len(idterms); i+=2 {
            id := strconv(idterms[i])
            ids = append(ids, id)
        }

        dt[py] = ids
    }

    return dt, err
}

func (t *Dict) LoadTerms(idfile string) (map[string]*idInfo, error) {
    f, e:= os.Open(pyfile)
    if e != nil {
        return nil, e
    }

    hzdt := make(map[string]*idInfo)
    bf := make([]byte, 1024)
    var err error
    linenum := 0
    allfreq := 0
    
    for {
        n, e := ReadLine(r, bf)
        if e != nil && e!= io.EOF {
            return nil,e
        }
        if n == 0 {
            err = e
            break
        }

        st := bytes.FieldsFunc(bf[:n], IsAscSpace)
        if len(st) < 3 {
            return nil, stringError.StringError(fmt.Sprintf("file format error with line %d", linenum+1))
        }

        hz = string(st[0])
        freq = strconv.Atoi(string(st[1]))
        order = strconv.Atoi(string(st[2]))
        hzdt[linenum] = &idInfo{hz, freq, order}

        linenum +=1
        allfreq += freq
    }

    if err != nil {
        return nil, err
    }
    if allfreq == 0 {
        return hzdt, stringError(StringError("0 freqs"))
    }

    for _, v := range hzdt {
        v.rate = float64(v.freq)/float64(allfreq)
        if v.rate != 0 {
            v.rate = -1.0*math.Log(v.rate)
        }else {
            v.rate = minrate
        }
    }

    return hzdt, nil
}


func (t *Dict)GetHzFromId(int id) (string, error) {
    if t, ok:= t.idHz(id); ok {
        return t.hz, nil
    }
    return "", stringError(StringError("not found"))
}













