package easysyntax

import (
	"errors"
	"strconv"
)

// url 中可以嵌入的原生字符只有'-' '_' '~' '.' 四个，其他均会被转义。'_' 已经被用于 cquery 分隔多个 id，'-' '~' 不容易被肉眼分辨，因此这里使用 '.' 和 '-'
const digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.-"

var (
	base   int    = 64
	maxVal uint64 = uint64(1)<<uint(base) - 1
)

func Atoi_64(s string) (n uint64, err error) {
	for _, c := range []byte(s) {
		var d byte
		switch {
		case '0' <= c && c <= '9':
			d = c - '0'
		case 'a' <= c && c <= 'z':
			d = c - 'a' + 10
		case 'A' <= c && c <= 'Z':
			d = c - 'A' + 36
		case c == '.':
			d = '0' - '0' + 62
		case c == '-':
			d = '0' - '0' + 63
		}

		n *= uint64(base)
		n1 := n + uint64(d)
		if n1 < n || n1 > maxVal {
			return maxVal, rangeError("Atoi_64", s)
		}
		n = n1
	}
	return n, nil
}

func Itoa_64(i uint64) (s string) {
	if i == 0 {
		return "0"
	}

	b := []byte{}
	for j := 0; i != 0; j++ {
		mod := i % uint64(base)
		i /= uint64(base)

		b = append(b, digits[mod])
	}
	b = reverse(b)
	return string(b)
}

func reverse(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

var ErrRange = errors.New("value out of range")

// A NumError records a failed conversion.
type NumError struct {
	Func string // the failing function (ParseBool, ParseInt, ParseUint, ParseFloat, ParseComplex)
	Num  string // the input
	Err  error  // the reason the conversion failed (e.g. ErrRange, ErrSyntax, etc.)
}

func (e *NumError) Error() string {
	return "strconv." + e.Func + ": " + "parsing " + strconv.Quote(e.Num) + ": " + e.Err.Error()
}

func (e *NumError) Unwrap() error { return e.Err }

func rangeError(fn, str string) *NumError {
	return &NumError{fn, str, ErrRange}
}
