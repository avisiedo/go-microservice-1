package echo

import (
	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	"github.com/avisiedo/go-microservice-1/internal/interface/interactor"
	presenter "github.com/avisiedo/go-microservice-1/internal/interface/presenter/sync/echo"
)

type private struct {
	i interactor.Private
}

// NewPrivate
func NewPrivate(i interactor.Private) presenter.Private {
	if i == nil {
		panic(common_err.ErrNil("i"))
	}
	return newPrivate(i)
}

func newPrivate(i interactor.Private) *private {
	return &private{
		i: i,
	}
}

// TODO Add below the implementation for the Private API for privatePresenter
