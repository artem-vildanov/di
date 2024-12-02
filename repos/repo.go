package repos

type Repo interface {
	GetAll() []string
}

type RepoImpl struct {
}

func (RepoImpl) GetAll() []string {
	println("GetAll from RepoImpl fired")
	return []string{
		"hello",
		"world",
	}
}