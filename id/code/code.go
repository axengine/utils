package code

import "fmt"

var CHARS = []string{"F", "L", "G", "W", "5", "X", "C", "3",
	"9", "Z", "M", "6", "7", "Y", "R", "T", "2", "H", "S", "8", "D", "V", "E", "J", "4", "K",
	"Q", "P", "U", "A", "N", "B"}

const (
	CHARS_LENGTH = int(32)
	//CODE_LENGTH  = int(10)
	SALT   = int(123456789)
	PRIME1 = int(3)
	PRIME2 = int(11)
)

// Gen 根据id生成对应的字符串，通常用于 邀请码
func Gen(id int64, codeLen int) string {
	id = id*int64(PRIME1) + int64(SALT)

	b := make([]int64, codeLen)
	b[0] = id
	for i := 0; i < int(codeLen)-1; i++ {
		b[i+1] = b[i] / int64(CHARS_LENGTH)
		//按位扩散
		b[i] = (b[i] + int64(i)*b[0]) % int64(CHARS_LENGTH)
	}
	b[5] = (b[0] + b[1] + b[2] + b[3] + b[4]) * int64(PRIME1) % int64(CHARS_LENGTH)

	codeIndexArray := make([]int64, codeLen)
	for i := 0; i < int(codeLen); i++ {
		idx := int(i) * int(PRIME2) % int(codeLen)
		codeIndexArray[i] = b[int(idx)]
	}

	var codestr string
	for _, v := range codeIndexArray {
		codestr = codestr + fmt.Sprintf("%s", CHARS[v])
	}
	return codestr
}
