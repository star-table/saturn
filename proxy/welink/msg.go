package welink

import (
	"github.com/galaxy-book/saturn/model/context"
	"github.com/galaxy-book/saturn/model/req"
	"github.com/galaxy-book/saturn/model/resp"
)

func (w *welinkProxy) SendMsg(ctx *context.Context, req req.SendMsgReq) resp.SendMsgResp {
	return resp.SendMsgResp{Resp: resp.SucResp()}
}
