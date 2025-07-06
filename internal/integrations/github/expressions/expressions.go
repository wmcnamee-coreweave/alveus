package expressions

type Expression string

func (e Expression) String() string {
	return To(string(e))
}

func To(val string) string {
	return "${{" + val + "}}"
}

func Inputs(val string) Expression {
	return Expression("inputs." + val)
}

func Secrets(val string) Expression {
	return Expression("secrets." + val)
}
