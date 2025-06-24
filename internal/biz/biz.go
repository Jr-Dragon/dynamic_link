package biz

import (
	"github.com/google/wire"
	"github.com/jr-dragon/dynamic_link/internal/biz/link"
)

var ProviderSet = wire.NewSet(
	link.NewLink,
)
