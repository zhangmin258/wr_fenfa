package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"math/rand"
	"math"
	"time"
	"regexp"
	"encoding/base64"
	"encoding/json"
	"strings"
	"crypto/des"
	"crypto/cipher"
	"bytes"
)

// 随机数种子
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// md5加密
func MD5(data string) string {
	m := md5.Sum([]byte(data))
	return hex.EncodeToString(m[:])
}

func PageCount(count, pagesize int) int {
	if count%pagesize > 0 {
		return count/pagesize + 1
	} else {
		return count / pagesize
	}
}

func StartIndex(page, pagesize int) int {
	if page > 1 {
		return (page - 1) * pagesize
	}
	return 0
}

func ErrNoRow() string {
	return "<QuerySeter> no row found"
}

// 获取数字随机字符
func GetRandDigit(n int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(n)+"d", rnd.Intn(int(math.Pow10(n))))
}

// 验证是否是手机号
func Validate(mobileNum string) bool {
	reg := regexp.MustCompile(Regular)
	return reg.MatchString(mobileNum)
}

// 处理验证码获取5次处理Or 登录错误5次处理
func CheckPwd5Time(key string) int {
	if RunMode == "debug" {
		return 1
	}
	count := 5
	if Rc.IsExist(key) {
		count, _ = Rc.RedisInt(key)
	}
	count--
	if count < 0 {
		return count
	}
	Rc.Put(key, count, GetTodayLastSecond())
	return count
}

// 获取今天剩余秒数
func GetTodayLastSecond() time.Duration {
	today := GetToday(FormatDate) + " 23:59:59"
	end, _ := time.ParseInLocation(FormatDateTime, today, time.Local)
	return time.Duration(end.Unix()-time.Now().Local().Unix()) * time.Second
}

func GetToday(format string) string {
	today := time.Now().Format(format)
	return today
}

// base64 encrypt
func Base64Encrypt(origData []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(origData))
}

func StringsToJSON(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}
	return jsons
}

// 序列化
func ToString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

// 隐藏手机号中间四位
func AccountDispose(account string) string {
	var err error
	var reg *regexp.Regexp
	var placeStr1 string
	if len(account) == 11 {
		reg, err = regexp.Compile("(\\d{3})\\d{4}(\\d{4})")
		placeStr1 = "****"
	} else {
		return account
	}
	if err != nil {
		return account
	}
	if reg.MatchString(account) == true {
		submatch := reg.FindStringSubmatch(account)
		return submatch[1] + placeStr1 + submatch[2]
	}
	return account
}

// 数组反转
func ExchangeList(idListDesc []int) (idListAsc []int) {
	l := len(idListDesc)
	idListAsc = make([]int, l)
	for i, id := range idListDesc {
		idListAsc[l-i-1] = id
	}
	return
}

// 截取小数点后几位
func SubFloatToString(f float64, m int) string {
	n := strconv.FormatFloat(f, 'f', -1, 64)
	if n == "" {
		return ""
	}
	if m >= len(n) {
		return n
	}
	newn := strings.Split(n, ".")
	if m == 0 {
		return newn[0]
	}
	if len(newn) < 2 || m >= len(newn[1]) {
		return n
	}
	return newn[0] + "." + newn[1][:m]
}

// 拼接get请求
func JoinGetUrl(url string, params map[string]string) (getUrl string) {
	url += "?"
	i := 1
	for k, v := range params {
		url += k + "=" + v
		if i < len(params) {
			url += "&"
		}
		i++
	}
	return url
}

// 格式化时间
func FormatTimeToString(t time.Time) string {
	t1 := t.Format(FormatDateTime)
	switch t1 {
	case "0001-01-01 00:00:00":
		return "-"
	default:
		return t1
	}
}

// float64转字符串 保留2位小数
func Float64ToString(m float64) string {
	return fmt.Sprintf("%.2f", m)
}

// 银行卡后四位
func BankCardFormat(card string) string {
	if len(card) < 10 {
		return card
	}
	rs := []rune(card)
	length := len(rs) - 1
	return string(rs[length-4:length])
}

const key = `h*.d;cy7x_12dkx?#j39fdl!` //api数据加密、解密key

func DesBase64Decrypt(crypted []byte) []byte {
	result, _ := base64.StdEncoding.DecodeString(string(crypted))
	origData, err := TripleDesDecrypt(result, []byte(key))
	if err != nil {
		panic(err)
	}
	return origData
}

// 3DES解密
func TripleDesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//des3 + base64 encrypt
func DesBase64Encrypt(origData []byte) []byte {
	result, err := TripleDesEncrypt(origData, []byte(key))
	if err != nil {
		panic(err)
	}
	return []byte(base64.StdEncoding.EncodeToString(result))
}

// 3DES加密
func TripleDesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

