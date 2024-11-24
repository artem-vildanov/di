package services

import "di/repos"


type Service struct {
	InnerService *NestedService `di:"enrich"`
	Repo repos.RepoInterface `di:"enrich"`
	SharedClient *SharedClient `di:"enrich"`
}

func (Service) AfterEnriched() {

}

func (Service) AfterInstantiated() {

}

func (s *Service) DoSmth() {
	s.SharedClient.ClientLogic()
	s.InnerService.SomeFunc()
}


type NestedService struct {
	AnotherService *AnotherService `di:"instantiate"`
	SharedClient *SharedClient `di:"enrich"`
}

func (NestedService) AfterInstantiated() {

}

func (n NestedService) AfterEnriched() {
	n.SharedClient.ClientLogic()
	n.AnotherService.GreatLogic()
}

func (NestedService) SomeFunc() {
}


type AnotherService struct {
}

func (AnotherService) AfterInstantiated() {
}


func (AnotherService) AfterEnriched() {
}

func (a AnotherService) GreatLogic() {
	//fmt.Print("\n")
	//for i := 0; i < 10; i++ {
		//fmt.Print("_logic_")
	//}
	//fmt.Print("\n")
}

type SharedClient struct {
	SomeValue int
	AnotherValue string
	AnotherService *AnotherService `di:"enrich"`
}

func (s *SharedClient) AfterInstantiated() {
	s.SomeValue = 79
	s.AnotherValue = "akjwndkajwdbnajwhdkbawkjdhbawjdkhba"
}

func (SharedClient) AfterEnriched() {
}

func (s SharedClient) ClientLogic() {
	//fmt.Printf("\nSomeValue: %d, AnotherValue %s \n", s.SomeValue, s.AnotherValue)
	s.AnotherService.GreatLogic()
}