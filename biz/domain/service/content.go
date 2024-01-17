package service

import "github.com/google/wire"

type IContentDomainService interface {
}
type ContentDomainService struct {
}

var ContentDomainServiceSet = wire.NewSet(
	wire.Struct(new(ContentDomainService), "*"),
	wire.Bind(new(IContentDomainService), new(*ContentDomainService)),
)
