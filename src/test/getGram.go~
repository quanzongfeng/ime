package main

import (
    "dict"
    "fmt"
    "os"
)

func main() {
    source := os.Args[1]
    dest := os.Args[2]
    e := dict.ChooseGramDict(source, dest)
    if e!= nil {
        fmt.Println(e)
    }
}

