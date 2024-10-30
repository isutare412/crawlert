package port

import "github.com/isutare412/crawlert/internal/core/model"

type QueryApplier interface {
	ApplyQuery(jsonBytes []byte) (model.QueryResult, error)
}
