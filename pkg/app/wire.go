package app

import (
	"github.com/google/wire"
	"github.com/isuquo/copper-test/pkg/templates"
)

var WireModule = wire.NewSet(
	templates.WireModule,
)
