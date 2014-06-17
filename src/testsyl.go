package main

import (
    "syllable"
    "fmt"
)

func main() {
    a := "ha"
    fmt.Println(syllable.IsShengMu(string(a[0])))
    fmt.Println(syllable.IsYunMu(string(a[1])))
    fmt.Println(syllable.IsSyllable(a))
    fmt.Println(syllable.GetSylByPrefixString(a))

    a = "nimenhaoa"
    for i:=0; i< len(a); i++ {
        re := syllable.SegPy(a, i) 
        for _, su := range re {
            fmt.Println(syllable.GetSylById(su.GetId()), su.GetStart(), su.GetEnd(), su.GetFlag())
        }
    }


}
