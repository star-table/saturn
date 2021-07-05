package resp

import (
	"errors"
	"fmt"
)

type Resp struct {
	Code int    `json:"code"` // proxy resp code
	Suc  bool   `json:"suc"`
	Msg  string `json:"msg"`
}

func ErrResp(err error) Resp {
	return Resp{
		Code: -1,
		Suc:  false,
		Msg:  err.Error(),
	}
}

func SucResp() Resp {
	return Resp{
		Code: 0,
		Suc:  true,
		Msg:  "Success",
	}
}

func (r *Resp) Error() error {
	if r.Suc {
		return nil
	}
	return errors.New(fmt.Sprintf("err[%v]: %v", r.Code, r.Msg))
}
