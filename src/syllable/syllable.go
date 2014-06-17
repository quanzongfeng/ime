package syllable
import (
    "sort"
//    "fmt"
)

const (
    ShengmuFlag = 0x1
    YunmuFlag = 0x2
    SyllableFlag = 0x4
)

//contain "" as zero_shengmu, zero_yunmu
var sm_list = []string {"b","p","m","f","d","t","n","l","g","k","h","j","q","x","zh","ch","sh","r","z","c","s","y","w",""}
var ym_list = []string {"a","o","e","i","u","v","ai","ei","ui","ao","ou","iu","ie","ve","er","an","en","in","un","ang","eng","ing","ong","ia","iao","ian","iang","iong","ua","uai","uan","uang","uo","van","vn","io","ion","on",""}
var syl6 = []string{"zhuang","shuang","chuang"}
var syl5 = []string{"shang","zhong","sheng","xiang","cheng","liang","chang","zheng","jiang","chong","guang","zhang","zhuan","qiang","chuan","huang","niang","kuang","xiong","qiong", "shuai","zhuai","chuai","jiong","shuan"}
var syl4 = []string{"shuo","zhan","shui","qing","tian","neng","xiao","nian","yong","yang","mian","shen","fang","dang","qian","xian","jing","hang","cong","dian","jian","tong","dong","ding","shou","ming","zhen","quan","guan","deng","gong","bian","bing","biao","ting","pian","zong","rang","ying","yuan","ping","kong","xing","shao","geng","jiao","wang","shan","kuai","feng","zhao","lian","suan","tiao","hong","zhou","shei","ling","rong","chan","huan","ceng","gang","zhun","chen","xuan","song","duan","chao","leng","meng","peng","diao","tuan","mang","huai","luan","guai","reng","liao","chai","bang","chun","tang","ruan","pang","long","kang","zhui","nong","shun","zhuo","zeng","ning","lang","zhua","cang","chou","niao","chui","piao","qiao","juan","miao","kuan","heng","nuan","teng","sang","zang","zuan","zhai","shua","seng","beng","weng","keng","shai","cuan","nang","chuo","zhei", "chua", "fiao"}
var syl3 = []string{"shi","zai","ren","you","zhe","lai","dao","jiu","men","wei","xia","yao","guo","hua","chu","dei","hui","kan","dui","mei","xin","hao","hou","dou","hai","ran","jia","wen","xue","duo","zhi","zuo","kai","suo","tou","hen","xie","ben","dan","jin","fen","lao","yin","san","zou","zui","gei","zhu","bei","wai","gao","yan","yue","tai","gan","bie","cai","jun","jue","kou","bai","min","lun","zen","nei","que","jie","huo","shu","dai","che","qin","wan","nan","chi","fan","gai","bao","fei","gen","qie","man","liu","zao","hei","lin","she","qiu","ban","luo","cun","lei","tan","yun","sui","cen","cuo","kao","lie","gou","pin","sha","nin","lou","mao","cha","fou","xiu","tui","pao","nao","mou","ruo","dun","gui","tao","mai","qun","lan","cao","pai","ken","lue","han","pan","hun","tie","rou","niu","xun","pei","sun","tuo","nuo","kuo","zha","gua","nai","zan","kun","sen","sai","mie","zun","can","lia","sao","gun","diu","qia","rui","die","rao","cui","sou","run","bin","kui","zei","kua","pen","tun","cou","ang","nie","nen","miu","pie","pou","nue","tei","dia","nou","kei","den","rua","eng", "num", "rua"}
var syl2 = []string{"yi","de","wo","le","bu","ta","ge","ni","da","zi","di","na","ye","he","li","ke","fa","me","qu","ba","er","ne","ma","qi","mu","ru","yu","ri","wu","ti","ji","nv","ci","xi","ai","an","bi","se","lu","si","te","xu","du","la","fu","hu","ze","ju","gu","pa","re","zu","tu","ku","su","ya","po","mi","pi","ha","mo","fo","bo","ka","lv","za","pu","ce","ao","en","ou","nu","wa","cu","sa","ca","ga","lo","yo","ei"}
var syl1 = []string{"a","e","o"}
var syl_list = [][]string {syl1,syl2,syl3,syl4,syl5,syl6}


type PyIdUnit struct {
    id      int
    start   int
    end     int
    flag    int
}

func (pi *PyIdUnit) GetId() int {
    return pi.id
}
func (pi *PyIdUnit)GetStart() int {
    return pi.start
}
func (pi *PyIdUnit)GetEnd() int {
    return pi.end
}
func (pi *PyIdUnit)GetFlag() int {
    return pi.flag
}

func (pi *PyIdUnit)IsSyllable() bool { 
    return pi.flag & SyllableFlag != 0
}

func (pi *PyIdUnit)IsShengmu() bool {
    return pi.flag & ShengmuFlag != 0
}

func (pi *PyIdUnit)IsYunmu() bool {
    return pi.flag & YunmuFlag != 0
}

func (pi *PyIdUnit)Len() int {
    return 1
}
//for sort
type pyIdSlice []*PyIdUnit 
func (p pyIdSlice) Len() int            {return len(p)}
func (p pyIdSlice) Less(i,j int) bool   {return p[i].flag > p[j].flag}
func (p pyIdSlice) Swap(i,j int)        {p[i], p[j] = p[j], p[i]}


type Syllable struct {
    Sm  []string    "shengmu list"
    Ym  []string    "yunmu  list"
    Syl [][]string  "yinjie lists"
    ordersyl []string   "ordered syls"
    index map[byte]int  "for prefix,first index"
    sylId map[string]int    "for small memory"
    idSyl map[int]string    "for full syl from id"
}

//for memory limit, to build syl_id_map
func (syl *Syllable)makeSylId() {
    if syl.sylId == nil {
        syl.sylId = make(map[string]int)
        syl.idSyl = make(map[int]string)
    }
    id := 0
    for _,s := range syl.Sm {
        if _, ok := syl.sylId[s]; !ok {
            syl.sylId[s] = id
            syl.idSyl[id] = s
            id += 1
        }
    }
    for _, y:= range syl.Ym {
        if _, ok := syl.sylId[y]; !ok {
            syl.sylId[y] = id
            syl.idSyl[id] = y
            id += 1
        }
    }

    for _, syls := range syl.Syl {
        for _, t:= range syls {
            if _, ok := syl.sylId[t]; !ok {
                syl.sylId[t] = id
                syl.idSyl[id] = t
                id +=1
            }
        }
    }
}

//get string from sylid
func (syl *Syllable)GetSylById(id int) string {
    if syl.sylId == nil {
        syl.makeSylId()
    }
    t, ok:= syl.idSyl[id]
    if !ok {
        return ""
    }
    return t
}

//get id from string
func (syl *Syllable)GetIdBySyl(s string) int {
    if syl.sylId == nil {
        syl.makeSylId()
    }
    t, ok := syl.sylId[s]
    if !ok {
        return -1
    }
    return t
}


//judge $1 is shengmu or not
func (syl *Syllable)IsShengMu(sm string) bool {
    for _, t := range syl.Sm {
        if sm == t {
            return true
        }
    }
    return false
}


func (syl *Syllable)IsYunMu(ym string) bool {
    for _, t := range syl.Ym {
        if ym == t {
            return true
        }
    }
    return false
}

//judge $1 is syllable or not
func (syl *Syllable)IsSyllable(py string)bool  {
    ln := len(py)
    if ln > len(syl.Syl) || ln==0 {
        return false
    }

    for _, s := range syl.Syl[ln-1] {
        if s == py {
            return true
        }
    }
    return false
}

//judge $1,$2 can be syllable or not
func (syl *Syllable) IsComposeAble(sm, ym string) bool {
    if syl.IsShengMu(sm) == false {
        return false
    }
    if syl.IsYunMu(ym) == false {
        return false
    }
    py := sm+ym
    return syl.IsSyllable(py)
}

//sort Syl for prefix get
func (syl *Syllable)sortSyllable() {
    if len(syl.Syl) == 0 {
        return
    }

    if syl.ordersyl == nil {
        syl.ordersyl = make([]string, 0, 1)
    }
    for _, sl := range syl.Syl{
        for _, ssl := range sl {
            syl.ordersyl = append(syl.ordersyl, ssl)
        }
    }

    sort.Strings(syl.ordersyl)

    if syl.index == nil {
        syl.index = make(map[byte]int)
    }

    lastch := syl.ordersyl[0][0]
    syl.index[lastch] = 0
    for i, sl := range syl.ordersyl {
        cl := sl[0]
        if cl == lastch {
            continue
        }else {
            syl.index[cl] = i
            lastch = cl
        }
    }
//    fmt.Println(syl.ordersyl)
//    fmt.Println(syl.index)
}

//get Syls by prefix py
func (syl *Syllable)GetSylByPrefixString(pres string) []string {
    ln := len(pres)
    if ln > len(syl.Syl) || ln == 0 {
        return nil
    }

    if syl.ordersyl == nil {
        syl.sortSyllable()
    }
    if syl.ordersyl == nil {
        return nil
    }
    
    ch := pres[0]
    ich, ok := syl.index[ch]
    if !ok {
        return nil
    }

//    fmt.Println(pres, ich, syl.ordersyl[ich])
    start := -1
    end := len(syl.ordersyl)
    for i:= ich; i < len(syl.ordersyl) ;i++ {
        cch := syl.ordersyl[i]

  //      fmt.Println(i, cch)
        if cch[0] != ch {
            break
        }
        if len(cch) >= ln && cch[:ln] == pres {
            if start == -1{
                start = i
            }
        }else {
            if start != -1 {
                end = i
                break
            }
        }
    }
    if (start == -1) {
        return nil
    }
    return syl.ordersyl[start:end]  //return no copy,means all goroutine use the same []
}
//segment py[start:],return all pyUnidt 
func (syl *Syllable)SegPy(py string, start int) []*PyIdUnit{
    ln := len(py)
    re := make([]*PyIdUnit, 0,0)
    flag := 0
    n := 0
    for i:= 1; i<= ln-start && i <= 6 ; i++ {
        flag = 0
        spy := py[start:start+i]
        id := syl.GetIdBySyl(spy)
        if id == -1 {
            continue
        }

        if IsShengMu(spy) {
            flag |= ShengmuFlag
        }
        if IsYunMu(spy) {
            flag |= YunmuFlag
        }
        if IsSyllable(spy) {
            flag |= SyllableFlag
        }

        if flag == 0 {
            panic("syllable id error")
        }

        idU := new(PyIdUnit)
        idU.id = id
        idU.start = start
        idU.end = start + i
        idU.flag = flag
        re = append(re, idU)
        n += 1
    }
    if n == 0 {
        return nil
    }
    
    sort.Sort(pyIdSlice(re))
    return re
}

var DefaultPySyllable *Syllable = &Syllable{ Sm : sm_list,Ym : ym_list,  Syl : syl_list}

func IsShengMu(sm string) bool {
    return DefaultPySyllable.IsShengMu(sm)
}

func IsYunMu(ym string) bool {
    return DefaultPySyllable.IsYunMu(ym)
}
func IsSyllable(py string) bool {
    return DefaultPySyllable.IsSyllable(py)
}
func IsComposeAble(sm, ym string) bool {
    return DefaultPySyllable.IsComposeAble(sm,ym)
}

func GetSylByPrefixString(py string) []string {
    return DefaultPySyllable.GetSylByPrefixString(py)
}
func SegPy(py string, start int) []*PyIdUnit {
    return DefaultPySyllable.SegPy(py, start) 
}
func GetSylById(id int) string {
    return DefaultPySyllable.GetSylById(id)
}
