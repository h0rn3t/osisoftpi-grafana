package plugin

import "testing"

func baseStreamableQuery() Query {
	enable := true
	return Query{
		Pi: PIWebAPIQuery{
			IsPiPoint: true,
			EnableStreaming: &struct {
				Enable *bool `json:"enable"`
			}{Enable: &enable},
		},
	}
}

func TestIsStreamable_PIPointEnabled(t *testing.T) {
	q := baseStreamableQuery()
	if !q.isStreamable() {
		t.Fatal("expected PI point query with streaming enabled to be streamable")
	}
}

func TestIsStreamable_ExcludesSummary(t *testing.T) {
	q := baseStreamableQuery()
	summaryEnable := true
	basis := "TimeWeighted"
	types := []SummaryType{{Value: SummaryTypeValue{Value: "Average"}}}
	q.Pi.Summary = &QuerySummary{
		Enable: &summaryEnable,
		Basis:  &basis,
		Types:  &types,
	}
	if q.isStreamable() {
		t.Fatal("expected summary query to be excluded from streaming")
	}
}

func TestIsStreamable_ExcludesInterpolated(t *testing.T) {
	q := baseStreamableQuery()
	q.Pi.Interpolate.Enable = true
	if q.isStreamable() {
		t.Fatal("expected interpolated query to be excluded from streaming")
	}
}

func TestIsStreamable_ExcludesRecordedValues(t *testing.T) {
	q := baseStreamableQuery()
	enable := true
	q.Pi.RecordedValues = &struct {
		Enable       *bool   `json:"enable"`
		MaxNumber    *int    `json:"maxNumber"`
		BoundaryType *string `json:"boundaryType"`
	}{Enable: &enable}
	if q.isStreamable() {
		t.Fatal("expected recorded values query to be excluded from streaming")
	}
}

func TestIsStreamable_ExcludesExpression(t *testing.T) {
	q := baseStreamableQuery()
	q.Pi.Expression = "Tag1 + Tag2"
	if q.isStreamable() {
		t.Fatal("expected expression query to be excluded from streaming")
	}
}

func TestIsStreamable_ExcludesAFAttribute(t *testing.T) {
	q := baseStreamableQuery()
	q.Pi.IsPiPoint = false
	if q.isStreamable() {
		t.Fatal("expected AF attribute (non-PI-point) query to be excluded from streaming")
	}
}

func TestIsUsingStreaming_WithoutExperimental(t *testing.T) {
	ds := newTestDatasource()
	streaming := true
	experimental := false
	ds.dataSourceOptions.UseStreaming = &streaming
	ds.dataSourceOptions.UseExperimental = &experimental
	if !ds.isUsingStreaming() {
		t.Fatal("expected streaming enabled without experimental mode")
	}
}
