package api

import (
	"github.com/google/wire"

	basev1 "github.com/jr-dragon/dynamic_link/api/base/v1"
	linkv1 "github.com/jr-dragon/dynamic_link/api/link/v1"
)

var ProviderSet = wire.NewSet(
	basev1.NewRoute,
	linkv1.NewRoute,
)
