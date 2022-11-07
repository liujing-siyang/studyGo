package server

import (
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSetServer = wire.NewSet(NewGRPCServer, NewHTTPServer)
