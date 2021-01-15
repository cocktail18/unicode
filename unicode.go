package unicode

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"regexp"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

const (
	unicodeReplacementChar = '\uFFFD'     // Unicode replacement character
	maxRune                = '\U0010FFFF' // Maximum valid Unicode code point.

	// 0xd800-0xdc00 encodes the high 10 bits of a pair.
	// 0xdc00-0xe000 encodes the low 10 bits of a pair.
	// the value is those 20 bits plus 0x10000.
	unicodeSurr1 = 0xd800
	unicodeSurr2 = 0xdc00
	unicodeSurr3 = 0xe000

	unicodeSurrSelf = 0x10000
)

func UnicodeToString(from string) (string, error) {
	reg1 := regexp.
		MustCompile("(\\\\u[0-9a-fA-F]{4})+")

	var loc []int
	toBuff := strings.Builder{}
	for {
		loc = reg1.FindStringIndex(from)
		if loc == nil {
			toBuff.WriteString(from)
			break
		}
		toBuff.WriteString(from[:loc[0]])
		tmp, err := unicodeToString(from[loc[0]:loc[1]])
		if err != nil {
			return "", err
		}
		toBuff.WriteString(tmp)
		if len(from) <= loc[1] {
			break
		}
		from = from[loc[1]:]
	}
	return toBuff.String(), nil
}

func unicodeToString(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	s := make([]uint16, 0, 0)
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		err = binary.Read(br, binary.BigEndian, &r)
		if err != nil {
			return
		}
		s = append(s, r)
	}
	to = DecodeUTF16ToString(s)
	return
}

// 参考https://github.com/lianggx6/goutf16
func DecodeUTF16ToString(s []uint16) string {
	n := 0
	for i := 0; i < len(s); i++ {
		switch r := s[i]; {
		case r < unicodeSurr1, unicodeSurr3 <= r:
			// normal rune
			n += utf8.RuneLen(rune(r))
		case unicodeSurr1 <= r && r < unicodeSurr2 && i+1 < len(s) &&
			unicodeSurr2 <= s[i+1] && s[i+1] < unicodeSurr3:
			// valid surrogate sequence
			n += utf8.RuneLen(utf16.DecodeRune(rune(r), rune(s[i+1])))
			i++
		default:
			// invalid surrogate sequence
			n += utf8.RuneLen(unicodeReplacementChar)
		}
	}
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < len(s); i++ {
		switch r := s[i]; {
		case r < unicodeSurr1, unicodeSurr3 <= r:
			// normal rune
			b.WriteRune(rune(r))
		case unicodeSurr1 <= r && r < unicodeSurr2 && i+1 < len(s) &&
			unicodeSurr2 <= s[i+1] && s[i+1] < unicodeSurr3:
			// valid surrogate sequence
			b.WriteRune(utf16.DecodeRune(rune(r), rune(s[i+1])))
			i++
		default:
			// invalid surrogate sequence
			b.WriteRune(unicodeReplacementChar)
		}
	}
	return b.String()
}

func EncodeStringToUTF16(s string) []uint16 {
	n := 0
	for _, v := range s {
		n++
		if v >= 0x10000 {
			n++
		}
	}

	a := make([]uint16, n)
	n = 0
	for _, v := range s {
		switch {
		case 0 <= v && v < unicodeSurr1, unicodeSurr3 <= v && v < unicodeSurrSelf:
			// normal rune
			a[n] = uint16(v)
			n++
		case unicodeSurrSelf <= v && v <= maxRune:
			// needs surrogate sequence
			r1, r2 := utf16.EncodeRune(v)
			a[n] = uint16(r1)
			a[n+1] = uint16(r2)
			n += 2
		}
	}
	return a[:n]
}
