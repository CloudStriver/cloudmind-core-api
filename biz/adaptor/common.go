package adaptor

import (
	"context"
	"errors"
	bizerrors "github.com/CloudStriver/go-pkg/utils/errors"
	"github.com/CloudStriver/go-pkg/utils/util/log"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	hertz "github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
	"strings"

	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
)

var _ propagation.TextMapCarrier = &headerProvider{}

type headerProvider struct {
	headers *protocol.ResponseHeader
}

// Get a value from metadata by key
func (m *headerProvider) Get(key string) string {
	return m.headers.Get(key)
}

// Set a value to metadata by k/v
func (m *headerProvider) Set(key, value string) {
	m.headers.Set(key, value)
}

// Keys Iteratively get all keys of metadata
func (m *headerProvider) Keys() []string {
	out := make([]string, 0)

	m.headers.VisitAll(func(key, value []byte) {
		out = append(out, string(key))
	})

	return out
}

func PostProcess(ctx context.Context, c *app.RequestContext, req, resp any, err error) {
	//log.CtxInfo(ctx, "[%s] req=%s, resp=%s, err=%v", c.Path(), util.JSONF(req), util.JSONF(resp), err)
	b3.New().Inject(ctx, &headerProvider{headers: &c.Response.Header})

	switch {
	case err == nil:
		c.JSON(hertz.StatusOK, resp)
	case errors.Is(err, consts.ErrNotAuthentication):
		c.JSON(hertz.StatusUnauthorized, err.Error())
	case errors.Is(err, consts.ErrForbidden):
		c.JSON(hertz.StatusForbidden, err.Error())
	default:
		startIndex := strings.Index(err.Error(), "desc =")
		if startIndex == -1 {
			log.CtxError(ctx, "internal error, err=%s", err)
			code := hertz.StatusInternalServerError
			c.String(code, hertz.StatusMessage(code))
		} else {
			startIndex += len("desc = ")
			c.JSON(http.StatusBadRequest, &bizerrors.BizError{
				Msg: err.Error()[startIndex:],
			})
		}
	}
}
