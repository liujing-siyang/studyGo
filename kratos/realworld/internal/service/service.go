package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSetService = wire.NewSet(NewRealWorldService)
