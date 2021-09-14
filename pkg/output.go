package pkg

type Output struct {
	Stdout string
	Stderr string
}

func (o Output) Equals(o2 Output) bool {
	return (o.Stderr == o2.Stderr) && (o.Stdout == o2.Stdout)
}
