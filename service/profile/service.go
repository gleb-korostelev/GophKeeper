package profile

import "github.com/gleb-korostelev/GophKeeper/tools/db"

type service struct {
	db db.IAdapter
}

func NewService(db db.IAdapter) *service {
	return &service{db: db}
}
