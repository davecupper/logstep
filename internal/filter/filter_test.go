package filter_test

import (
	"testing"

	"github.com/yourorg/logstep/internal/filter"
)

func TestMatchNoRules(t *testing.T) {
	f := filter.New(nil)
	ok, err := f.Match([]byte(`{"level":"info","msg":"hello"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected match with no rules")
	}
}

func TestMatchEq(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Field: "level", Operator: "eq", Value: "error"},
	})

	ok, _ := f.Match([]byte(`{"level":"error","msg":"boom"}`))
	if !ok {
		t.Fatal("expected match for level=error")
	}

	ok, _ = f.Match([]byte(`{"level":"info","msg":"fine"}`))
	if ok {
		t.Fatal("expected no match for level=info")
	}
}

func TestMatchContains(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Field: "msg", Operator: "contains", Value: "timeout"},
	})

	ok, _ := f.Match([]byte(`{"msg":"connection timeout reached"}`))
	if !ok {
		t.Fatal("expected match containing 'timeout'")
	}

	ok, _ = f.Match([]byte(`{"msg":"all good"}`))
	if ok {
		t.Fatal("expected no match")
	}
}

func TestMatchExists(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Field: "trace_id", Operator: "exists"},
	})

	ok, _ := f.Match([]byte(`{"trace_id":"abc123","msg":"traced"}`))
	if !ok {
		t.Fatal("expected match when field exists")
	}

	ok, _ = f.Match([]byte(`{"msg":"no trace"}`))
	if ok {
		t.Fatal("expected no match when field absent")
	}
}

func TestMatchInvalidJSON(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Field: "level", Operator: "eq", Value: "info"},
	})
	_, err := f.Match([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
