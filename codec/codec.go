package codec

const (
	ENCODE_DEFAULT = iota
	ENCODE_BIT_NOT
	ENCODE_BYTE_RVS
	ENCODE_LOOP_XOR
)

func Encode(tp int, buff []byte) {
	switch tp {
	case ENCODE_BIT_NOT:
		{
			for i := 0; i < len(buff); i++ {
				buff[i] = ^buff[i]
			}
		}
	case ENCODE_BYTE_RVS:
		{
			length := len(buff)
			if len(buff)%2 != 0 {
				length -= 1
			}
			for i := 0; i < length; i += 2 {
				buff[i], buff[i+1] = buff[i+1], buff[i]
			}
		}
	case ENCODE_LOOP_XOR:
		{
			length := len(buff)
			if length <= 1 {
				break
			}
			for i := 0; i < length-1; i++ {
				buff[i+1] ^= buff[i]
			}
			buff[0] ^= buff[length-1]
		}
	case ENCODE_DEFAULT:
	default:

	}
}

func Decode(tp int, buff []byte) {
	switch tp {
	case ENCODE_BIT_NOT:
		{
			for i := 0; i < len(buff); i++ {
				buff[i] = ^buff[i]
			}
		}
	case ENCODE_BYTE_RVS:
		{
			length := len(buff)
			if len(buff)%2 != 0 {
				length -= 1
			}
			for i := 0; i < length; i += 2 {
				buff[i], buff[i+1] = buff[i+1], buff[i]
			}
		}
	case ENCODE_LOOP_XOR:
		{
			length := len(buff)
			if length <= 1 {
				break
			}
			buff[0] ^= buff[length-1]
			for i := length - 1; i > 0; i-- {
				buff[i] ^= buff[i-1]
			}
		}
	case ENCODE_DEFAULT:
	default:

	}
}
