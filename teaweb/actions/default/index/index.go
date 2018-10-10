package index

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaweb/helpers"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct {
	Auth *helpers.UserMustAuth
}) {
	this.RedirectURL("/dashboard")
}
