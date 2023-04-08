package sessions

import "github.com/google/wire"

var WireModule = wire.NewSet(
	NewQueries,
)
