package services

import (
	"di/repos"
	"fmt"
)

type Service struct {
	InnerService *NestedService
	Repository repos.RepoInterface
}

func (Service) AfterInstantiated() {
}
func (s *Service) AfterEnriched() {
}

func (s *Service) DoSmth() {
	s.InnerService.SomeFunc()
}


type NestedService struct {
	AnotherService *AnotherService `di:"instantiate"`
}

func (NestedService) AfterInstantiated() {
}
func (n NestedService) AfterEnriched() {
	n.AnotherService.GreatLogic()
}

func (NestedService) SomeFunc() {
	fmt.Println("hello from nested service")
}


type AnotherService struct {
}

func (a AnotherService) GreatLogic() {
	fmt.Print("\n")
	for i := 0; i < 10; i++ {
		fmt.Print("_logic_")
	}
	fmt.Print("\n")
}