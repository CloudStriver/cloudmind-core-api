package service

import "github.com/google/wire"

type ITradeDomainService interface {
}
type TradeDomainService struct {
}

var TradeDomainServiceSet = wire.NewSet(
	wire.Struct(new(TradeDomainService), "*"),
	wire.Bind(new(ITradeDomainService), new(*ITradeDomainService)),
)
