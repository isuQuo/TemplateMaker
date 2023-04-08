package app

import (
	"github.com/google/wire"
	"github.com/isuquo/copper-test/pkg/logs"
	"github.com/isuquo/copper-test/pkg/sessions"
	"github.com/isuquo/copper-test/pkg/templates"
	"github.com/isuquo/copper-test/pkg/users"
)

var WireModule = wire.NewSet(
	sessions.WireModule,

	users.WireModule,

	logs.WireModule,

	templates.WireModule,
)
