package common

import (
    "log"
    "io"
    "bytes"
    "fmt"
)

const (
    DefaultLineLength=1024
)

type LineReader struct {
    bf      []byte
    start   int         //start index of data to precessed
    end     int         //end index of data to processed

    r       io.Reader
    eflag   error //file end flag
    maxline int
}

func NewLineReader(rin io.Reader, maxline int) *LineReader{
    rout := new(LineReader)
    rout.r = rin
    if maxline == -1 {
        maxline = DefaultLineLength
    }
    if maxline < DefaultLineLength {
        maxline = DefaultLineLength
    }
    rout.maxline = maxline
    rout.bf = make([]byte, 5*maxline)   //可以保存至少5行数据
    return rout
}

func (lr *LineReader)clear() {
    n := lr.end - lr.start
    if n == 0 {     //no left data, return
        lr.start = 0
        lr.end = 0
        return
    }
    if lr.eflag != nil { //文件已到末尾或者出错，不需要再读入数据，则不需要整理
        return
    }
    if len(lr.bf)-lr.end > lr.maxline { //剩余空间足够保存一个新行
        return
    }
    if n > lr.maxline { //剩余数据足够一个新行
        return 
    }
    //left space not enough to store a new line 
    copy(lr.bf, lr.bf[lr.start:lr.end])
    lr.start = 0
    lr.end = n
    return
}

//return $2: nil 
func (lr *LineReader)getline() ([]byte, error) {
    if lr.start >= lr.end {
        if lr.eflag != nil {
            return nil, lr.eflag
        }
        return nil,nil 
    }

    linetag := bytes.IndexByte(lr.bf[lr.start:lr.end], '\n')
    if linetag == -1 {
        if lr.eflag == io.EOF{      //表明读到文件末尾
            re := make([]byte, lr.end - lr.start)
            copy(re, lr.bf[lr.start:lr.end])
            return re, lr.eflag
        }
        return nil,lr.eflag     //出错或未读入一行
    }

    //说明找到了换行符,这种情况下，不返回错误或者结束符
    newstart := lr.start + linetag + 1

    if linetag > 0 && lr.bf[lr.start + linetag -1] == '\r' {
        linetag -= 1
    }
    re := make([]byte, linetag)
    copy(re, lr.bf[lr.start:lr.start + linetag])
    lr.start = newstart
    return re, nil
}

func (lr *LineReader)ReadLine() ([]byte, error) {
    if lr.start < lr.end { //still have data not processed
        re, ef := lr.getline()
        if ef != nil {      //end or error, return
            return re, ef
        }else if re != nil {    //get line , return
            return re, ef
        }
        //no line ,no error
    }
    lr.clear()  //确保lr有足够的空间读入一个新行
    if lr.eflag == nil {
        n, e := lr.r.Read(lr.bf[lr.end:])
        lr.eflag = e                //读取文件时才会设置eflag
        lr.end += n
        if e!= nil && e!= io.EOF {
            return nil, e
        }
    
        //正常情况下，至此应该读入了一个新行
        re, ef := lr.getline() 
        if ef == nil && re == nil { //说明读入数据不够一个新行, 或者文件阻塞
            lr.eflag = StrError(fmt.Sprintf("read line blocked or line length more than %d", lr.maxline) )
            return nil, lr.eflag
        }
        return re, ef
    }

    return nil, lr.eflag
}


//read line less than 1024 bytes
func ReadLine(r io.ReadSeeker, bf []byte)(n int, e error) {
    newline := make([]byte, 10240)
    n, e = r.Read(newline)
    if  e != nil && e != io.EOF {
        return 0, e
    }

    lineTag := bytes.IndexByte(newline[:n], '\n')
    if lineTag == -1 {
        if e != io.EOF {
            return 0, StrError(fmt.Sprintf("line too long than %d bytes", n)) 
        }

        if n > len(bf) {
            return 0, StrError(fmt.Sprintf("too small bytes to hold line"))
        }
        copy(bf, newline[:n])
        return n, io.EOF
    }

    if lineTag > len(bf) {
        return 0, StrError(fmt.Sprintf("too small bytes to hold line"))
    }

    back := lineTag -n +1
    _, err:=r.Seek(int64(back), 1)
    if err != nil {
        log.Fatal("Read error")
    }

    copy(bf, newline[:lineTag])
    return lineTag, e
}




