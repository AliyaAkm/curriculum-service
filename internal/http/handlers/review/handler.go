package review

import "github.com/go-playground/validator/v10"

type Handler struct {
	client   client
	validate *validator.Validate
}

func NewHandler(client client, validate *validator.Validate) *Handler {
	return &Handler{
		client:   client,
		validate: validate,
	}
}
