package cache

import (
	"fmt"
)

type Cmder interface {
	Args() []interface{}
	String() string

	SetErr(error)
	Err() error
}

type baseCmd struct {
	args []interface{}
	err  error
}

func (cmd *baseCmd) SetErr(e error) {
	cmd.err = e
}

func (cmd *baseCmd) Args() []interface{} {
	return cmd.args
}

func (cmd *baseCmd) Err() error {
	return cmd.err
}

type StringCmd struct {
	baseCmd
	val string
}

var _ Cmder = (*StringCmd)(nil)

func NewStringCmd(value string, err error, args ...interface{}) *StringCmd {
	return &StringCmd{
		baseCmd: baseCmd{
			args: args,
			err:  err,
		},
		val: value,
	}
}

func (cmd *StringCmd) String() string {
	return fmt.Sprintf("commandType:%T, val: %s, args: %+v", cmd, cmd.val, cmd.args)
}

func (cmd *StringCmd) Result() (string, error) {
	return cmd.val, cmd.err
}

type SliceCmd struct {
	baseCmd

	val []interface{}
}

var _ Cmder = (*SliceCmd)(nil)

func NewSliceCmd(value []interface{}, err error, args ...interface{}) *SliceCmd {
	return &SliceCmd{
		baseCmd: baseCmd{
			args: args,
			err:  err,
		},
		val: value,
	}
}

func (cmd *SliceCmd) String() string {
	return fmt.Sprintf("commandType:%T, val: %s, args: %+v", cmd, cmd.val, cmd.args)
}

func (cmd *SliceCmd) Result() ([]interface{}, error) {
	return cmd.val, cmd.err
}

type StringSliceCmd struct {
	baseCmd

	val []string
}

var _ Cmder = (*StringSliceCmd)(nil)

func NewStringSliceCmd(value []string, err error, args ...interface{}) *StringSliceCmd {
	return &StringSliceCmd{
		baseCmd: baseCmd{
			args: args,
			err:  err,
		},
		val: value,
	}
}

func (cmd *StringSliceCmd) String() string {
	return fmt.Sprintf("commandType:%T, val: %s, args: %+v", cmd, cmd.val, cmd.args)
}

func (cmd *StringSliceCmd) Result() ([]string, error) {
	return cmd.val, cmd.err
}

type StatusCmd struct {
	baseCmd

	val string
}

var _ Cmder = (*StatusCmd)(nil)

func NewStatusCmd(value string, err error, args ...interface{}) *StatusCmd {
	return &StatusCmd{
		baseCmd: baseCmd{
			args: args,
			err:  err,
		},
		val: value,
	}
}

func (cmd *StatusCmd) String() string {
	return fmt.Sprintf("commandType:%T, val: %s, args: %+v", cmd, cmd.val, cmd.args)
}

func (cmd *StatusCmd) Result() (string, error) {
	return cmd.val, cmd.err
}

type IntCmd struct {
	baseCmd

	val int
}

var _ Cmder = (*IntCmd)(nil)

func NewIntCmd(value int, err error, args ...interface{}) *IntCmd {
	return &IntCmd{
		baseCmd: baseCmd{
			args: args,
			err:  err,
		},
		val: value,
	}
}

func (cmd *IntCmd) String() string {
	return fmt.Sprintf("commandType:%T, val: %d, args: %+v", cmd, cmd.val, cmd.args)
}

func (cmd *IntCmd) Result() (int, error) {
	return cmd.val, cmd.err
}
