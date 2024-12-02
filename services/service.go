package services

import (
	"di/repos"
	"fmt"
)


type Service struct {
	InnerService *NestedService
	SharedClient *SharedClient
	Repo repos.Repo
}

func (s *Service) Construct(
	innerService *NestedService, 
	sharedClient *SharedClient, 
	repo repos.Repo,
) {
	s.InnerService = innerService
	s.SharedClient = sharedClient
	s.Repo = repo
}

func (s *Service) DoSmth() {
	println("DoSmth from Service fired")
	fmt.Printf("data from repo: %v", s.Repo.GetAll())
	s.SharedClient.ClientLogic()
	s.InnerService.SomeFunc()
}

type NestedService struct {
	AnotherService *AnotherService
	SharedClient *SharedClient
}

func (n *NestedService) Construct(
	anotherService *AnotherService,
	sharedClient *SharedClient,
) {
	n.AnotherService = anotherService
	n.SharedClient = sharedClient
}

func (NestedService) SomeFunc() {
	println("SomeFunc from NestedService fired")
}


type AnotherService struct {
}

func (a AnotherService) GreatLogic() {
	println("another service fired clientlogic")
}

type SharedClient struct {
	AnotherService *AnotherService
}

func (s *SharedClient) Construct(anotherService *AnotherService) {
	s.AnotherService = anotherService
}

func (s SharedClient) ClientLogic() {
	println("shared service fired clientlogic")
	s.AnotherService.GreatLogic()
}