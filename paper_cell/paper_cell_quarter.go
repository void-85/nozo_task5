package paper_cell

type Quarter uint8

func From_text(char byte) Quarter {

	if char == '#' {
		return 0b00001111
	}

	return 0b00000000
}

func (q *Quarter) Invert_45() {

	*q = 0b00000000

}

func (q *Quarter) To_text() byte {

	switch *q {
	case 0b00001111:
		return '#'
	default:
		return '.'
	}

}
