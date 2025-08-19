package plan

import planP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/port"

type service struct {
	repo planP.Repo
}

func New(r planP.Repo) planP.Service { return &service{repo: r} }
