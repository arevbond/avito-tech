package banner

import (
	"banners/internal/storage"
)

type Service interface {
}

type ServiceImpl struct {
	Storage storage.Storage
}
