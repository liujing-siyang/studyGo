package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSetBiz = wire.NewSet(NewSocialUsecase, NewUserUsecase)
