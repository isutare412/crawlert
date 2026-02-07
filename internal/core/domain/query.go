package domain

type QueryResult struct {
	Matched   bool
	Variables map[string]string
}
