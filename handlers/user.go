package handlers

import (
	"di/services"
)

type UserHandler struct {
	Service *services.Service
}

func (u UserHandler) AfterInstantiated() {
}

func (u UserHandler) AfterEnriched() {
}

func (h *UserHandler) Handle() {
	h.Service.Repository.GetById("qwe")
	h.Service.DoSmth()
}

