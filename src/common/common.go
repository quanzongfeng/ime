package common
import (
    "runtime"
)


func IsAscSpace(t rune)bool {
    if t == ' ' || t== '\t' {
        return true
    }
    return false
}

func IsAscQuota(t rune)bool {
    if t == '\'' || t == '"' {
        return true
    }
    return false
}

func FilterSingleQuota(t []byte) []byte{
    ln := len(t)
    nbyte := make([]byte, ln)
    j:=0
    for i:=0; i< ln ;i++ {
        if t[i] != '\'' {
            nbyte[j] = t[i]
            j++
        }
    }
    return nbyte[:j]
}

func GetGbkDzLen(bf []byte) int {
    ln := len(bf)
    if ln == 0 {
        return 0
    }

    if bf[0] < 128 {
        return 1
    }else if bf[0] > 128 {
        if ln>1 {
            return 2
        }else {
            return -1
        }
    }
    return 1
}

//"我们爱祖国", return []int {2,4,6,8,10}
func GetGbkHzIndexsList(bf []byte)([]int, error) {
    ln := len(bf)
    if ln == 0 {
        return nil, nil
    }

    re := make([]int, 1, 1+ln/2)
    re[0] = 0

    for i:=0; i< ln; {
        li := GetGbkDzLen(bf[i:])
        if li < 0 {
            return nil, StrError("error code in gbk")
        }
        re = append(re, i+li)
        i += li
    }

    return re, nil
}

func SetProcessNums(n int) {
    t:= runtime.NumCPU()
    if n == -1 || n > t {
        n = t
    }
    runtime.GOMAXPROCS(n)
}

func SetCpuNums(n int) {
    t:= runtime.NumCPU()
    if n <=0  || n > t {
        n = t
    }
    runtime.GOMAXPROCS(n)
}
