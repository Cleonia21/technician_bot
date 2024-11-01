package useCases

import (
	"technician_bot/temporary/entities"
)

type Repository interface {
	Get(request entities.Request) (entities.Response, error)
}

type Presenter interface {
	Send(entities.Response) error
}

type Handler struct {
	repo      Repository
	presenter Presenter
}

func NewHandler(repo Repository, presenter Presenter) *Handler {
	return &Handler{repo: repo, presenter: presenter}
}

func (rh *Handler) Handler(request entities.Request) error {
	response, err := rh.repo.Get(request)
	if err != nil {
		return err
	}
	_ = rh.presenter.Send(response)
	return nil
}
