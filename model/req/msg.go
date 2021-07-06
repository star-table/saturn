package req

type SendMsgReq struct {
	OpenIds []string `json:"openIds"`
	UserIds []string `json:"userIds"`
	DeptIds []string `json:"deptIds"`
	ChatIds []string `json:"chatIds"`
	MsgType string   `json:"msgType"` // Tag: text,image,share_chat,interactive
	Msg     Msg      `json:"content"`
}

type Msg struct {
	Title    string    `json:"title"`
	Color    string    `json:"color"`
	Url      Url       `json:"url"`
	Elements []Element `json:"elements"`

	Text        string `json:"text"`        // 文本消息
	ImageKey    string `json:"imageKey"`    // 图片消息，image_key 可以通过图片上传接口获得
	ShareChatID string `json:"shareChatId"` // 发送群名片， 群聊id
}

type Element struct {
	Type               string                 `json:"tag"`                // Tag: div,hr,img,action,note,button,select_static,select_person,picker_datetime
	ContentText        string                 `json:"contentText"`        // 内容文字
	ContentFields      []MsgField             `json:"contentFields"`      // 内容字段
	ImgKey             string                 `json:"imgKey"`             //图片key
	ImgTitle           string                 `json:"imgTitle"`           // 图片标题
	ImgAlt             string                 `json:"imgAlt"`             // 图片hover时展示说明，为空则不展示
	ImgWidth           int                    `json:"imgWidth"`           // 自定义图片的最大展示宽度。278~580
	ImgMode            string                 `json:"imgMode"`            // 图片显示模式。crop_center：居中裁剪模式，对长图会限高，并居中裁剪后展示，fit_horizontal：平铺模式，完整展示上传的图片
	ImgPreview         bool                   `json:"imgPreview"`         // 点击后是否放大图片
	ButtonText         string                 `json:"buttonText"`         // 按钮文本
	ButtonUrl          Url                    `json:"buttonUrl"`          // 跳转链接
	ButtonType         string                 `json:"buttonType"`         // 配置按钮样式 default/primary/danger
	ButtonValue        map[string]interface{} `json:"buttonValue"`        // 点击后返回业务方
	ButtonConfirm      MsgConfirm             `json:"buttonConfirm"`      // 二次确认的弹框
	SelectPlaceholder  string                 `json:"selectPlaceholder"`  // 占位符，无默认选项时必须有
	SelectDefaultValue string                 `json:"selectDefaultValue"` // 默认选项的value字段值
	SelectOptions      []MsgOption            `json:"selectOptions"`      // 待选选项
	SelectConfirm      MsgConfirm             `json:"selectConfirm"`      // 二次确认的弹框
	SelectValue        map[string]interface{} `json:"selectValue"`        // 点击后返回业务方
	DateDefaultValue   string                 `json:"dateDefaultValue"`   // 日期默认数据 yyyy-MM-dd HH:mm
	DatePlaceholder    string                 `json:"datePlaceholder"`    // 占位符
	DateValue          map[string]interface{} `json:"dateValue"`          // 点击后返回业务方
	DateConfirm        MsgConfirm             `json:"dateConfirm"`        // 二次确认的弹框
	Notes              []Element              `json:"notes"`              // 备注 text对象或image元素
	Actions            []Element              `json:"actions"`            // 交互
	ActionLayout       string                 `json:"actionLayout"`       // 交互元素布局, "bisected", "trisection", "flow"
}

type MsgField struct {
	IsShort bool   `json:"isShort"`
	Text    string `json:"text"`
}

type MsgConfirm struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type MsgOption struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

type Url struct {
	Url     string `json:"url"`
	Android string `json:"android"`
	IOS     string `json:"ios"`
	PC      string `json:"pc"`
}
