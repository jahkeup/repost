package config

import (
	"text/template"

	"github.com/jahkeup/repost/pkg/emailfuncs"
)

var funcMap = template.FuncMap{
	"emailName":   emailfuncs.Name,
	"emailDomain": emailfuncs.Domain,
	"emailUser":   emailfuncs.User,
	"emailTag":    emailfuncs.Tag,
}
