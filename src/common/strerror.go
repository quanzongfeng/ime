package common

type StrError string

func (s StrError)Error() string{
    return string(s)
}

var StrError1 = StrError("1")
var StrError0 = StrError("0")
