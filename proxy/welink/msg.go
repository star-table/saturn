package welink

import (
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
)

func (w *welinkProxy) SendMsg(ctx *context.Context, req req.SendMsgReq) resp.SendMsgResp {
	return resp.SendMsgResp{Resp: resp.SucResp()}
}
