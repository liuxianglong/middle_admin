package code

import (
	"context"
	"fmt"
	"middle/internal/model"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"

	"strconv"
	"strings"
)

const (
	SuccessCode = 0
	CommonCode  = iota + 9
	AuthCode
)

type ICfg interface {
	GetCode(errMsg string) int
	BuildCode(code int) int
}

var CodeMap = map[string]ICfg{
	"system": &BaseSystemCfg{ //注system暂时不定preCode，其他的需定义
		Cfg: systemMap,
		//PreCode: SystemCode,
	},
	"common": &BaseSystemCfg{
		Cfg:     commonMap,
		PreCode: CommonCode,
	},

	"auth": &BaseSystemCfg{
		Cfg:     authMap,
		PreCode: AuthCode,
	},
}
var CodeError = new(codeError)

type codeError struct {
	error
}

func (this *codeError) ErrorNew(ctx context.Context, i18n *model.I18n, errMsg string, values ...interface{}) error {
	err := gerror.New(i18n.Tf(ctx, errMsg, values...))
	code := this.GetErrCode(errMsg)
	return gerror.WrapCode(gcode.New(code, "", nil), err, "")
}

func (this *codeError) New(ctx context.Context, errMsg string, values ...interface{}) error {
	err := gerror.New(gi18n.Tf(ctx, errMsg, values...))
	code := this.GetErrCode(errMsg)
	return gerror.WrapCode(gcode.New(code, "", nil), err, "")
}

func (this *codeError) GetErrCode(errMsg string) int {
	code := -1
	split := strings.Split(errMsg, ".")
	firstMsg := split[0]

	if _, ok := CodeMap[firstMsg]; ok {
		obj := CodeMap[firstMsg]
		code = obj.GetCode(errMsg)
	}
	return code
}

type BaseSystemCfg struct {
	Cfg     map[string]int
	PreCode int
}

func (cfg *BaseSystemCfg) GetCode(errMsg string) int {
	if _, ok := cfg.Cfg[errMsg]; !ok {
		return -1
	}
	if cfg.PreCode > 0 {
		return cfg.BuildCode(cfg.Cfg[errMsg])
	}
	return cfg.Cfg[errMsg]
}

// BuildCode 由6位数字组成的code
func (cfg *BaseSystemCfg) BuildCode(code int) int {
	codeStr := fmt.Sprintf("%02d%04d", cfg.PreCode, code)
	atoi, _ := strconv.Atoi(codeStr)
	return atoi
}
