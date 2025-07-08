package aprint

import (
	"errors"
	"fmt"
	"strings"
)

type PrintAlign int

const (
	LEFT_ALIGN PrintAlign = iota
	RIGHT_ALIGN
)

func New(expectedColumns int, defaultDelimeter string, defaultAlign PrintAlign) *AlignedPrint {
	pb := &AlignedPrint{
		expectedColumns: expectedColumns,
		maximumLines:    make([]int, expectedColumns),
		contents:        make([][]string, 0),
		delimeters:      make([]string, expectedColumns-1),
		align:           make([]PrintAlign, expectedColumns),
	}

	for idx := range expectedColumns {
		pb.align[idx] = LEFT_ALIGN
	}
	for idx := range pb.delimeters {
		pb.delimeters[idx] = defaultDelimeter
	}

	return pb
}

type AlignedPrint struct {
	expectedColumns int
	maximumLines    []int
	contents        [][]string
	delimeters      []string
	align           []PrintAlign
}

func (p *AlignedPrint) SetDelimeter(column int, delimeter string) {
	p.delimeters[column] = delimeter
}
func (p *AlignedPrint) SetAlignment(column int, align PrintAlign) {
	p.align[column] = align
}

func (p *AlignedPrint) Add(content ...string) error {
	if len(content) != p.expectedColumns {
		return errors.New(
			fmt.Sprintf("Expected %d columns, got %d", p.expectedColumns, len(content)))
	}

	for idx, value := range content {
		length := len(value)
		if length >= p.maximumLines[idx] {
			p.maximumLines[idx] = length
		}
	}

	p.contents = append(p.contents, content)

	return nil
}

func (p *AlignedPrint) String() string {
	sb := strings.Builder{}

	for idx, line := range p.contents {
		for col, content := range line {
			if p.align[col] == LEFT_ALIGN {
				sb.WriteString(
					fmt.Sprintf("%-*s", p.maximumLines[col], content),
				)
			} else {
				sb.WriteString(
					fmt.Sprintf("%+*s", p.maximumLines[col], content),
				)
			}

			if col+1 != len(line) {
				sb.WriteString(p.delimeters[col])
			}
		}

		if idx+1 != len(p.contents) {
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}
