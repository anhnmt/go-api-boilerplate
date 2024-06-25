package authbusiness

import (
	userquery "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/query"
)

type Business struct {
	userQuery *userquery.Query
}

func New(
	userQuery *userquery.Query,
) *Business {
	return &Business{
		userQuery: userQuery,
	}
}
