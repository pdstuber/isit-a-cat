package idgenerator

import "github.com/rs/xid"

type Service struct{}

func (g *Service) GenerateID() string {
	return xid.New().String()
}
