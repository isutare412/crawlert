package query

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"

	"github.com/isutare412/crawlert/internal/core/model"
)

type Applier struct {
	checkQuery      *gojq.Code
	variableQueries map[string]*gojq.Code
}

func NewApplier(checkQuery string, variableQueries map[string]string) (*Applier, error) {
	check, err := compileJQQuery(checkQuery)
	if err != nil {
		return nil, fmt.Errorf("compiling check query: %w", err)
	}

	variables := make(map[string]*gojq.Code, len(variableQueries))
	for key, query := range variableQueries {
		q, err := compileJQQuery(query)
		if err != nil {
			return nil, fmt.Errorf("compiling query of variable %s: %w", key, err)
		}

		variables[key] = q
	}

	return &Applier{
		checkQuery:      check,
		variableQueries: variables,
	}, nil
}

func (e *Applier) ApplyQuery(jsonBytes []byte) (model.QueryResult, error) {
	var target any
	if err := json.Unmarshal(jsonBytes, &target); err != nil {
		return model.QueryResult{}, fmt.Errorf("unmarshaling into json: %w", err)
	}

	checkResult, err := queryFirstItem(e.checkQuery, target)
	if err != nil {
		return model.QueryResult{}, fmt.Errorf("applying check query: %w", err)
	}

	variables := make(map[string]string, len(e.variableQueries))
	for key, query := range e.variableQueries {
		result, err := queryFirstItem(query, target)
		if err != nil {
			return model.QueryResult{}, fmt.Errorf("applying variable '%s' query: %w", key, err)
		}
		variables[key] = result
	}

	return model.QueryResult{
		Matched:   isTruthyValue(checkResult),
		Variables: variables,
	}, nil
}

func queryFirstItem(query *gojq.Code, target any) (string, error) {
	var result string
	iter := query.Run(target)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok := v.(error); ok {
			return "", fmt.Errorf("iterating result: %w", err)
		}

		encoded, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("json marshaling query result: %w", err)
		}

		result = string(encoded)
		break
	}

	return result, nil
}

func compileJQQuery(s string) (*gojq.Code, error) {
	query, err := gojq.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("parsing jq query: %w", err)
	}

	code, err := gojq.Compile(query)
	if err != nil {
		return nil, fmt.Errorf("compiling jq query: %w", err)
	}

	return code, nil
}

func isTruthyValue(s string) bool {
	s = strings.TrimSpace(s)

	if strings.EqualFold(s, "true") {
		return true
	}

	num, err := strconv.Atoi(s)
	switch {
	case err != nil:
		return false
	case num > 0:
		return true
	}

	return false
}
