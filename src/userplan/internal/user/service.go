package user

import userP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/port"

type service struct {
	repo userP.Repo
}

func New(r userP.Repo) userP.Service { return &service{repo: r} }
