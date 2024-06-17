package userbusiness

type Business interface {
}

type business struct {
}

func New() Business {
	return &business{}
}
