package referral

import (
	"errors"
	"strings"
)

var (
	base    = "E8S2DZX9WYLTN6BGQF7P5IK3MJUAR4HV"
	decimal = 32
	pad     = "C"
	length  = 6
)

func SetBase(b string) {
	b = strings.ToUpper(strings.TrimSpace(b))
	if len(b) <= 0 {
		return
	}
	decimal = len(base)
	base = b
}

func SetPad(p string) error {
	p = strings.ToUpper(p)
	if strings.Contains(base, p) {
		return errors.New("pad should not exists in base")
	}
	pad = p
	return nil
}

func SetLength(n int) {
	length = n
}

// Encode uid to referral code
func Encode(uid uint64) (referral string) {
	id := uid
	mod := uint64(0)
	for id != 0 {
		mod = id % uint64(decimal)
		id = id / uint64(decimal)
		referral += string(base[mod])
	}
	resLen := len(referral)
	if resLen < length {
		referral += pad
		for i := 0; i < length-resLen-1; i++ {
			referral += string(base[(int(uid)+i)%decimal])
		}
	}
	return
}

// Decode code to uid
func Decode(referral string) (id uint64) {
	lenCode := len(referral)
	baseArr := []byte(base)       // string decimal to byte array
	baseRev := make(map[byte]int) // decimal data key to map
	for k, v := range baseArr {
		baseRev[v] = k
	}
	// find cover char addr
	isPad := strings.Index(referral, pad)
	if isPad != -1 {
		lenCode = isPad
	}
	r := 0
	for i := 0; i < lenCode; i++ {
		// if cover char , continue
		if string(referral[i]) == pad {
			continue
		}
		index, ok := baseRev[referral[i]]
		if !ok {
			return 0
		}
		b := uint64(1)
		for j := 0; j < r; j++ {
			b *= uint64(decimal)
		}
		id += uint64(index) * b
		r++
	}
	return id
}
