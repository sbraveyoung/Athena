package pool

import "unsafe"

//copy from std lib
//Use simple []byte instead of bytes.Buffer to avoid large dependency.
type Buffer []byte

func (b *Buffer) Write(p []byte) {
	*b = append(*b, p...)
}

func (b *Buffer) WriteString(s string) {
	*b = append(*b, s...)
}

func (b *Buffer) WriteByte(c byte) {
	*b = append(*b, c)
}

func (b Buffer) Len() int {
	return len(b)
}

func (b Buffer) String() string {
	//copy from strings.Builder
	return *(*string)(unsafe.Pointer(&b))
}

// https://segmentfault.com/a/1190000005006351
// func (b *Buffer) NewfromString(s string) []byte {
// x := (*[2]uintptr)(unsafe.Pointer(&s))
// h := [3]uintptr{x[0], x[1], x[1]}
// return *(*[]byte)(unsafe.Pointer(&h))
// }

// Join concatenates the elements of s to create a new byte slice. The separator
// sep is placed between elements in the resulting slice.
func Join(s []Buffer, sep []byte) Buffer {
	if len(s) == 0 {
		return Buffer{}
	}
	if len(s) == 1 {
		// Just return a copy.
		return append(Buffer{}, s[0]...)
	}
	n := len(sep) * (len(s) - 1)
	for _, v := range s {
		n += len(v)
	}

	b := make(Buffer, n)
	bp := copy(b, s[0])
	for _, v := range s[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], v)
	}
	return b
}
