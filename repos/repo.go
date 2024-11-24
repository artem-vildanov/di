package repos


type RepoInterface interface {
	GetById(id string)
}

type RepoProxy struct {
	RepoImpl *RepoImpl `di:"enrich"`
}

func (RepoProxy) AfterInstantiated() {}
func (RepoProxy) AfterEnriched() {}

func (r RepoProxy) GetById(id string) {
	r.RepoImpl.GetById(id)
}

type RepoImpl struct {

}

func (RepoImpl) AfterInstantiated() {}
func (RepoImpl) AfterEnriched() {}

func (RepoImpl) GetById(id string) {
}