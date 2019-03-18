package apps

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type AddNoticeReceiversAction actions.Action

// 添加通知接收人
func (this *AddNoticeReceiversAction) Run(params struct {
	AgentId string
	AppId   string
	Level   notices.NoticeLevel
}) {
	agentutils.InitAppData(this, params.AgentId, params.AppId, "noticeSetting")

	level := notices.FindNoticeLevel(params.Level)
	if level == nil {
		this.Fail("Level不存在")
	}

	this.Data["level"] = level

	// 媒介
	setting := notices.SharedNoticeSetting()
	mediaMaps := []maps.Map{}
	for _, media := range setting.Medias {
		if !media.On {
			continue
		}
		mediaType := notices.FindNoticeMediaType(media.Type)
		if mediaType == nil {
			continue
		}
		mediaMaps = append(mediaMaps, maps.Map{
			"id":               media.Id,
			"name":             media.Name,
			"typeName":         mediaType["name"],
			"type":             media.Type,
			"mediaDescription": mediaType["description"],
			"userDescription":  mediaType["user"],
		})
	}
	this.Data["medias"] = mediaMaps

	this.Show()
}

// 提交保存
func (this *AddNoticeReceiversAction) RunPost(params struct {
	AgentId string
	AppId   string
	Level   notices.NoticeLevel
	On      bool
	Name    string
	MediaId string
	User    string
	Must    actions.Must
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入接收人名称").
		Field("mediaId", params.MediaId).
		Require("请选择使用的媒介").
		Field("user", params.User).
		Require("请输入接收人标识")

	receiver := notices.NewNoticeReceiver()
	receiver.On = params.On
	receiver.Name = params.Name
	receiver.MediaId = params.MediaId
	receiver.User = params.User

	app.AddNoticeReceiver(params.Level, receiver)
	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
