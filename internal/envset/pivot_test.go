package envset

import (
	"testing"
)

func basePivotSets(t *testing.T) []*EnvSet {
	t.Helper()
	make := func(name, region string) *EnvSet {
		s, _ := New(name, "local")
		if region != "" {
			_ = s.Set("REGION", region)
		}
		_ = s.Set("APP", "envoy")
		return s
	}
	return []*EnvSet{
		make("alpha", "us-east-1"),
		make("beta", "eu-west-1"),
		make("gamma", "us-east-1"),
		make("delta", ""),
	}
}

func TestPivot_GroupsByValue(t *testing.T) {
	sets := basePivotSets(t)
	res, err := Pivot("REGION", sets...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Key != "REGION" {
		t.Errorf("expected key REGION, got %q", res.Key)
	}
	if len(res.Groups["us-east-1"]) != 2 {
		t.Errorf("expected 2 sets in us-east-1, got %d", len(res.Groups["us-east-1"]))
	}
	if len(res.Groups["eu-west-1"]) != 1 {
		t.Errorf("expected 1 set in eu-west-1, got %d", len(res.Groups["eu-west-1"]))
	}
	if len(res.Groups[""]) != 1 {
		t.Errorf("expected 1 set with missing key, got %d", len(res.Groups[""]))
	}
}

func TestPivot_PivotKeys_Sorted(t *testing.T) {
	sets := basePivotSets(t)
	res, _ := Pivot("REGION", sets...)
	keys := res.PivotKeys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 pivot keys, got %d", len(keys))
	}
	if keys[0] != "" {
		t.Errorf("expected first key to be empty string, got %q", keys[0])
	}
}

func TestPivot_InvalidKey(t *testing.T) {
	s, _ := New("x", "local")
	_, err := Pivot("123invalid", s)
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
}

func TestPivot_NoSets(t *testing.T) {
	_, err := Pivot("REGION")
	if err == nil {
		t.Fatal("expected error when no sets provided")
	}
}

func TestPivot_NilSetsSkipped(t *testing.T) {
	s, _ := New("only", "local")
	_ = s.Set("REGION", "ap-south-1")
	res, err := Pivot("REGION", nil, s, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Groups["ap-south-1"]) != 1 {
		t.Errorf("expected 1 set, got %d", len(res.Groups["ap-south-1"]))
	}
}
