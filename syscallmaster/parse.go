package syscallmaster

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type ParserResult struct {
	Number   int
	Audit    string
	Files    string
	Decl     string
	Comments string
}

type semantics int

const (
	NUMBER semantics = iota
	AUDIT
	FILES
	DECL
	COMMENTS
)

type state struct {
	text     []byte
	pos      int
	len      int
	sem      semantics
	ln       int
	col      int
	filename string
}

func Parse(dat []byte, filename string) []ParserResult {
	s := &state{
		text:     dat,
		pos:      0,
		ln:       1,
		col:      1,
		len:      len(dat),
		sem:      NUMBER,
		filename: filename,
	}

	return doParse(s)
}

func doParse(s *state) []ParserResult {
	ret := []ParserResult{}
	row := &ParserResult{}
	numrange := ""
	for {
		err := s.skipSpace()
		if err == io.EOF && s.sem == NUMBER {
			return ret
		} else if err != nil {
			fail(*s, err)
		}

		if s.sem == NUMBER {
			for !isNumber(s.peek()) {
				s.until('\n')
				s.pop()
				if s.pos == s.len {
					return ret
				}
			}
		}

		switch s.sem {
		case NUMBER:
			row = &ParserResult{}
			numrange = s.word()
			s.sem = AUDIT
			continue
		case AUDIT:
			fallthrough
		case FILES:
			w := s.word()
			if s.sem == AUDIT {
				row.Audit = w
				s.sem = FILES
			} else {
				row.Files = w
				s.sem = DECL
			}
			continue
		case DECL:
			fallthrough
		case COMMENTS:
			if s.peek() != '{' {
				break
			}
			start := s.pos
			s.until('}')
			s.pop()
			end := s.pos
			w := string(s.text[start:end])
			if s.sem == DECL {
				row.Decl = w
				s.sem = COMMENTS
			} else {
				row.Comments = w
				s.sem = NUMBER
			}
			continue
		default:
			fail(*s, fmt.Errorf("??"))
		}

		nr := [2]int{}
		if strings.Contains(numrange, "-") {
			nums := strings.Split(numrange, "-")
			nr[0], err = strconv.Atoi(nums[0])
			if err != nil {
				fail(*s, err)
			}
			nr[1], err = strconv.Atoi(nums[1])
			if err != nil {
				fail(*s, err)
			}
			nr[1] += 1
		} else {
			nr[0], err = strconv.Atoi(numrange)
			if err != nil {
				fail(*s, err)
			}
			nr[1] = nr[0] + 1
		}

		for n := nr[0]; n < nr[1]; n++ {
			row := &ParserResult{
				Number:   n,
				Audit:    row.Audit,
				Files:    row.Files,
				Decl:     row.Decl,
				Comments: row.Decl,
			}
			ret = append(ret, *row)
		}
		s.sem = NUMBER
	}

	return ret
}

func isPrintable(b byte) bool {
	return 33 <= b && b <= 126
}

func (s *state) untilSpace() error {
	for isPrintable(s.peek()) {
		_, err := s.pop()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *state) untilNumber() error {
	for s.peek() < '0' || '9' < s.peek() {
		_, err := s.pop()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *state) until(c byte) (err error) {
	for s.peek() != c {
		_, err = s.pop()
		if err != nil {
			return err
		}
	}
	return nil
}

func isNumber(b byte) bool {
	return '0' <= b && b <= '9'
}

func (s *state) skipSpace() error {
	for !isPrintable(s.peek()) {
		_, err := s.pop()
		if err != nil {
			return err
		}
	}
	return nil
}

func fail(s state, err error) {
	fmt.Printf("%s:%v:%v, Pos %v, expect %v\n", s.filename, s.ln, s.col, s.pos, s.sem)
	errpos := s.col - 1
	start := s.pos - (s.col - 1)
	end := s.pos
	for end < s.len && s.text[end] != '\n' {
		end += 1
	}
	fmt.Printf("%s", s.text[start:end])
	for range errpos {
		fmt.Print("~")
	}
	fmt.Print("^\n")
	panic(err)
}

func (s *state) word() string {
	start := s.pos
	s.untilSpace()
	// s.pop()
	end := s.pos
	ret := s.text[start:end]
	return string(ret)
}

func (s *state) peek() byte {
	if s.pos == s.len {
		fail(*s, fmt.Errorf("EOF"))
	}
	return s.text[s.pos]
}

func (s *state) pop() (byte, error) {
	if s.pos == s.len {
		return 0, io.EOF
	}
	ret := s.text[s.pos]
	s.pos += 1
	if ret == '\n' {
		s.ln += 1
		s.col = 1
	} else {
		s.col += 1
	}
	return ret, nil
}
