package gordle

import "strings"

type hint byte

const (
	ABSENT_CHARACTER hint = iota
	WRONG_POSITION
	CORRECT_POSITION
)

func (h hint) String() string {
	switch h {
	case ABSENT_CHARACTER:
		return "⬜️"
	case WRONG_POSITION:
		return "🟨"
	case CORRECT_POSITION:
		return "🟩"
	default:
		return "🟥"
	}
}

type feedback []hint

func (fb feedback) StringConcat() string {
	out := strings.Builder{}

	for _, h := range fb {
		out.WriteString(h.String())
	}

	return out.String()
}

func (fb feedback) Equal(other feedback) bool {
	if len(fb) != len(other) {
		return false
	}
	for i, v := range fb {
		if v != other[i] {
			return false
		}
	}
	return true
}
