package errorx_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/krizsx/gopkg/errorx"
)

func TestE(t *testing.T) {
	testCases := map[string]struct {
		args   []interface{}
		wanted *errorx.Errorx
	}{
		"空参数构造": {
			[]interface{}{}, &errorx.Errorx{},
		},
		"只指定错误码": {
			[]interface{}{errorx.ErrCode(1001)}, &errorx.Errorx{ErrCode: 1001},
		},
		"只指定operation,通过string": {
			[]interface{}{"错误信息"}, &errorx.Errorx{ErrOperation: "错误信息"},
		},
		"只设定Operation": {
			[]interface{}{errorx.ErrOperation("a.b")}, &errorx.Errorx{ErrOperation: "a.b"},
		},
		"全部设定": {
			[]interface{}{
				errorx.ErrCode(2021),
				errorx.ErrOperation("a.b"),
				errorx.ErrLogLevel("debug"),
				errors.New("新的错误"),
			},
			&errorx.Errorx{
				ErrCode:      2021,
				ErrOperation: "a.b",
				ErrLogLevel:  "debug",
				Err:          errors.New("新的错误"),
			},
		},
	}

	for name, tc := range testCases {
		assert.Equal(t, tc.wanted, errorx.E(tc.args...), name)
	}
}

func TestWrapErr(t *testing.T) {
	err := errors.New("原始错误")
	err = errorx.E(err, "包装第一层", errorx.ErrCode(1001))
	err = errorx.E(err, "包装第二层")
	t.Logf("错误码:%v, 操作列表: %v: 原始错误: %v", errorx.Code(err), errorx.Operations(err), err)
	errorx.Log(err)
}
