package api

import (
	"github.com/google/wire"

	linkv1 "github.com/jr-dragon/dynamic_link/api/link/v1"
)

var ProviderSet = wire.NewSet(linkv1.NewRoute)
