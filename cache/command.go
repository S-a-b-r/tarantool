package cache

type StringCmd struct {
	val string
	err error
}

func NewStringCmd(val string, err error) *StringCmd {
	return &StringCmd{val: val, err: err}
}

func (cmd *StringCmd) Result() (string, error) {
	return cmd.val, cmd.err
}
