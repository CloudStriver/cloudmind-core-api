package etcd

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/cloudwego/kitex/pkg/discovery"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
)

func NewEtcd(c *config.Config) discovery.Resolver {
	r, err := etcd.NewEtcdResolver(c.EtcdConf.Hosts)
	if err != nil {
		log.Fatal(err)
	}
	return r
}
