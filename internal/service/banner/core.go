package banner

import (
	"avito-tech/internal/storage"
)

type Service interface {
}

type ServiceImpl struct {
	Storage storage.Storage
}
