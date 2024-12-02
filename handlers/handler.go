package handlers

import (
	"di/services"
)

type PostHandler struct {
	SharedClient *services.SharedClient
	service *services.Service
}

func (p *PostHandler) Construct(
	sharedClient *services.SharedClient,
	service *services.Service,
) {
	p.SharedClient = sharedClient
	p.service = service
}

func (p PostHandler) Handle() {
	println("post handler Handle fired")
	p.service.DoSmth()
	p.SharedClient.ClientLogic()
}

