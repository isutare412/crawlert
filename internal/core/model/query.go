package model

type QueryResult struct {
	Matched   bool
	Variables map[string]string
}
