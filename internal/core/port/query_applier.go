package port

import "github.com/isutare412/crawlert/internal/core/domain"

type QueryApplier interface {
	ApplyQuery(jsonBytes []byte) (domain.QueryResult, error)
}
