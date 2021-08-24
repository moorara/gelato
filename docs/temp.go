package lookup

type Request struct {
	ID string
}

type Response struct {
	Name string
}

type Service interface {
	Lookup(*Request) (*Response, error)
}

type Provider struct {
	service Service
}

func (p *Provider) Provide(id string) (string, error) {
	req := &Request{
		ID: id,
	}

	resp, err := p.service.Lookup(req)
	if err != nil {
		return "", err
	}

	return resp.Name, nil
}
