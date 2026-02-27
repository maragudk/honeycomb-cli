package cmd_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"maragu.dev/is"

	"github.com/maragudk/honeycomb-cli/cmd"
	"github.com/maragudk/honeycomb-cli/honeycomb"
)

func newQueryServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/1/queries/requests":
			_ = json.NewEncoder(w).Encode(honeycomb.QueryResponse{ID: "q1"})

		case r.Method == http.MethodPost && r.URL.Path == "/1/query_results/requests":
			// Return immediately complete result for fast tests
			result := honeycomb.QueryResult{
				ID:       "r1",
				Complete: true,
			}
			result.Data.Results = []map[string]any{
				{"COUNT": float64(42), "status_code": "200"},
				{"COUNT": float64(7), "status_code": "500"},
			}
			_ = json.NewEncoder(w).Encode(result)
		}
	}))
}

func TestQueryCommand(t *testing.T) {
	t.Run("runs a query and displays table results", func(t *testing.T) {
		server := newQueryServer(t)
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"query", "--dataset", "requests", "--calculation", "COUNT",
			"--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		output := buf.String()
		is.True(t, contains(output, "42"))
		is.True(t, contains(output, "200"))
	})

	t.Run("runs a query and displays JSON results", func(t *testing.T) {
		server := newQueryServer(t)
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"query", "--dataset", "requests", "--calculation", "COUNT",
			"--json", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		var results []map[string]any
		is.NotError(t, json.Unmarshal(buf.Bytes(), &results))
		is.Equal(t, 2, len(results))
	})

	t.Run("handles no results", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodPost && r.URL.Path == "/1/queries/empty":
				_ = json.NewEncoder(w).Encode(honeycomb.QueryResponse{ID: "q1"})
			case r.Method == http.MethodPost && r.URL.Path == "/1/query_results/empty":
				result := honeycomb.QueryResult{ID: "r1", Complete: true}
				_ = json.NewEncoder(w).Encode(result)
			}
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"query", "--dataset", "empty", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		is.True(t, contains(buf.String(), "No results"))
	})
}

func TestParseCalculation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantOp   string
		wantCol  string
		wantErr  bool
	}{
		{name: "COUNT without column", input: "COUNT", wantOp: "COUNT"},
		{name: "AVG with column", input: "AVG:duration_ms", wantOp: "AVG", wantCol: "duration_ms"},
		{name: "P99 with column", input: "P99:duration_ms", wantOp: "P99", wantCol: "duration_ms"},
		{name: "case insensitive", input: "count", wantOp: "COUNT"},
		{name: "AVG without column returns error", input: "AVG", wantErr: true},
		{name: "unknown operator returns error", input: "YOLO", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calc, err := cmd.ParseCalculation(test.input)
			if test.wantErr {
				is.True(t, err != nil)
				return
			}
			is.NotError(t, err)
			is.Equal(t, test.wantOp, calc.Op)
			is.Equal(t, test.wantCol, calc.Column)
		})
	}
}

func TestParseFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantCol  string
		wantOp   string
		wantVal  string
		wantErr  bool
	}{
		{name: "equals", input: "status_code = 200", wantCol: "status_code", wantOp: "=", wantVal: "200"},
		{name: "not equals", input: "status_code != 500", wantCol: "status_code", wantOp: "!=", wantVal: "500"},
		{name: "greater than", input: "duration_ms > 100", wantCol: "duration_ms", wantOp: ">", wantVal: "100"},
		{name: "contains", input: "name contains foo", wantCol: "name", wantOp: "contains", wantVal: "foo"},
		{name: "invalid filter", input: "garbage", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filter, err := cmd.ParseFilter(test.input)
			if test.wantErr {
				is.True(t, err != nil)
				return
			}
			is.NotError(t, err)
			is.Equal(t, test.wantCol, filter.Column)
			is.Equal(t, test.wantOp, filter.Op)
			is.Equal(t, test.wantVal, fmt.Sprint(filter.Value))
		})
	}
}
