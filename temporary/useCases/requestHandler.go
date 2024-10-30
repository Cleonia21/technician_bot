package useCases

import (
	entities "technician_bot/temporary/entities"
)

type DataBase interface {
	GetResponse(request entities.Request) (entities.Response, error)
}

type RequestHandler struct {
	db DataBase
}

func NewRequestHandler(db DataBase) RequestHandler {
	return RequestHandler{db: db}
}

func (rh *RequestHandler) Handler(request entities.Request) (entities.Response, error) {
	return rh.db.GetResponse(request)
}
