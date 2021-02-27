package diff

// DiffType is the enum type for different diffing algs
type DiffType int

const (
	BASIC DiffType = iota
	MYERS
	LEVENSHTEIN
)

// ZenDiff is the interface that handles the diffing for the entire
// session
type ZenDiff interface {
	Calculate(out1, out2 []string) float64
	Print()
}

// NewDiff is a general purpose ZenDiff builder, just pass in the type!
func NewDiff(t DiffType, cname, tname string) ZenDiff {
	switch t {
	case BASIC:
		return &BasicDiff{CorrectName: cname, TestName: tname}
	case MYERS:
		return &MyersDiff{CorrectName: cname, TestName: tname}
	default:
		return nil
	}
}

// ZenNoDiffError is a general error struct
type ZenNoDiffError struct{}

func (z ZenNoDiffError) Error() string {
	return "No difference detected"
}
