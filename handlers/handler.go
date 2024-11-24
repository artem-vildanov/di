package handlers

import (
	"di/services"
)

type UserHandler struct {
	Service *services.Service `di:"enrich"`
}

func (h *UserHandler) Handle() {
	h.Service.DoSmth()
	h.Service.InnerService.AnotherService.GreatLogic()
}

type PostHandler struct {
	SharedClient *services.SharedClient `di:"enrich"`
}

func (p PostHandler) AfterInstantiated() {
}

func (PostHandler) AfterEnriched() {
}

func (p PostHandler) Handle() {
	p.SharedClient.SomeValue = 4545
	p.SharedClient.AnotherValue = "768766768768768766878"
	p.SharedClient.ClientLogic()
}

