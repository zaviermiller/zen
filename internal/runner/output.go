package runner

// enum type for different outputs
type OutputType int

const (
	PROMPT OutputType = iota
	RESPONSE
)

// output type that can differentiate between prompt and response outputs
type ZenOutput struct {
	Type    OutputType
	Content string
}

func (out ZenOutput) Text() string {
	return out.Content + "\n"
}
