package sentence

import (
    "heap"
    "fmt"
)

type WeightGraph struct {
    tw [][]*path
    h   *heap.HeapFactory
}

type path struct {
    idfrom  int
    idto    int
    left    *hzIdInfo
    right   *hzIdInfo
    weight  float64
    flag    int
    fx  float64
    gx  float64
    next    *path

}


func(p *path)Value() float64 {
    return p.fx
}

func (p *path)Compare(n heap.Value) (bool, error) {
    return p.Value() < n.Value(), nil
}


func BuildLattice(wg *WordGraph)(*WordGraph) {
    ln := len(wg.wg)
    newwg := new(WordGraph)     //get new wg

    newwg.wg = make([][]*hzIdInfo,1)    //for init start '$'
    newwg.wg[0] = make([]*hzIdInfo,1)
    startNode := &hzIdInfo{id:0,start:-1, end:0}    //here, end point to index in newwg
    newwg.wg[0][0] = startNode

    newwg.wg = append(newwg.wg, wg.wg...)   //copy hzinfo

    end := make([]*hzIdInfo,1)          //add end '$'
    end[0] = &hzIdInfo{id:0, start:ln, end:-1}
    newwg.wg = append(newwg.wg, end)

    nln := len(newwg.wg)
    for i:=0; i<nln; i++ {      //一次遍历，构建所有的右图和左图
        for _, hzsn := range newwg.wg[i] {
            end := hzsn.end                     //end 表示原wg中的数据end，即其在拼音串中的end，而不是在newwg中的end
            if end == -1 {          //reach right end
                continue
            }
            for _, rphz := range newwg.wg[end+1] {      //newwg 比wg多了一个起点.但wg中的start，end保持不变，所以这里要+1
                if hzsn.rpath == nil {
                    hzsn.rpath = make([]*path, 0,1)
                }
                nrp:=&path{idfrom:hzsn.id, idto:rphz.id, left:hzsn, right:rphz}
                hzsn.rpath = append(hzsn.rpath, nrp)

                if rphz.lpath == nil {
                    rphz.lpath = make([]*path, 0,1)
                }
                rphz.lpath = append(rphz.lpath, nrp)
            }
        }
    }
    newwg.dt = wg.dt
    return newwg
}

func (wg *WordGraph)Vitebi(choose int) (*WeightGraph, float64) {
    ln := len(wg.wg)
    if ln==0 {
        return nil, 0.0
    }

    tw := new(WeightGraph)
    tw.tw = make([][]*path, ln)
    //path 从右往左推，故0无用
    //    tw.tw[0] = make([]*path,1)
    //    tw.tw[0][0] = &path{idfrom:-1, idto:0,weight:0, right:wg.wg[0][0]}

    for i:=0; i< ln-1; i++ {
        if tw.tw[i] == nil {
            tw.tw[i] = make([]*path, 0)
        }
        for _, lnode := range wg.wg[i] {
            for _, rp := range lnode.rpath {
                rp.weight, rp.flag = wg.dt.GetWeight(rp.idfrom,rp.idto, choose)

                tw.tw[i] = append(tw.tw[i], rp)
                rnode := rp.right
                if rnode.bestcost == 0.0 {
                    rnode.bestcost = lnode.bestcost + rp.weight
                }else {
                    n := lnode.bestcost + rp.weight
                    if  n < rnode.bestcost {
                        rnode.bestcost = n
                    }
                }
            }
        }
    }

    if tw.tw[ln-1]== nil {
        tw.tw[ln-1] = make([]*path,1)
        tw.tw[ln-1][0] = &path{idfrom:0, idto:-1, left:wg.wg[ln-1][0], fx:wg.wg[ln-1][0].bestcost }
    }

    tw.h = heap.NewHeap()

    tw.h.Push(tw.tw[ln-1][0])

    return tw, wg.wg[ln-1][0].bestcost
}

func (tw *WeightGraph)Print() {
    for _, tws := range tw.tw {
        for _, p:= range tws {
            ln := p.left
            rn := p.right
            if rn != nil {
                fmt.Println("[",ln.start, ",", rn.end, "]", p.idfrom , p.idto, p.weight, p.flag)
            }else {
                fmt.Println("best result is :",p.fx)
            }
        }
    }
}


func (tw *WeightGraph)Next() ([]*path, float64, bool) {
    if tw.h == nil {
        return nil, 0,false
    }

//    fmt.Println("start get one path")
    index := 0
    for !tw.h.Empty() {
        mtop := tw.h.Pop()
        top, ok := mtop.(*path)
        if !ok {
            panic("error heap element")
        }
//        fmt.Println("Now is :", top)
//        fmt.Println(index)
        index +=1

        if top.idfrom == 0 && top.idto != -1 {
            re := make([]*path, 0, 10)
            for t:= top; t != nil && t.idto!= -1; t = t.next {
                re = append(re, t)
            }
            return re, top.fx, true
        }

        for _, p := range top.left.lpath {
            p.gx = top.gx + p.weight
            p.fx = p.left.bestcost + p.gx
            p.next = top
//            fmt.Println("Push into heap:", p)
            tw.h.Push(p)
        }
    }
    return nil,0.0, false
}

func (tw *WeightGraph)GetPath(num int) (re [][]*path, cl []float64){
    re = make([][]*path, 0, 1)
    cl = make([]float64,0, 1)
    flag := 0
    for i:=0; i< num ;i++  {
        r,z, ok := tw.Next() 
        if !ok {
            break
        }
        re = append(re, r)
        cl = append(cl, z)
        flag = 1
    }
    if flag == 0 {
        return nil, nil
    }
    return
}

