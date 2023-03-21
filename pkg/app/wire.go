package app

import (
	"github.com/google/wire"
	"github.com/isuquo/copper-test/pkg/logs"
	"github.com/isuquo/copper-test/pkg/templates"
)

var WireModule = wire.NewSet(
	logs.WireModule,

	templates.WireModule,
)
