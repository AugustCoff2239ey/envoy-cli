package envset

import (
	"testing"
)

func baseDigestSet() *EnvSet {
	es, _ := New("digest-test", "local")
	_ = es.Set("APP_HOST", "localhost")
	_ = es.Set("APP_PORT", "8080")
	_ = es.Set("APP_SECRET", "s3cr3t")
	return es
}

func TestDigest_AllKeys(t *testing.T) {
	es := baseDigestSet()
	res, err := Digest(es, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Digest == "" {
		t.Fatal("expected non-empty digest")
	}
	if len(res.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(res.Entries))
	}
}

func TestDigest_SelectedKeys(t *testing.T) {
	es := baseDigestSet()
	res, err := Digest(es, []string{"APP_HOST", "APP_PORT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
}

func TestDigest_Deterministic(t *testing.T) {
	es := baseDigestSet()
	r1, _ := Digest(es, nil)
	r2, _ := Digest(es, nil)
	if r1.Digest != r2.Digest {
		t.Fatal("digest should be deterministic")
	}
}

func TestDigest_ChangedValue(t *testing.T) {
	es1 := baseDigestSet()
	es2 := baseDigestSet()
	_ = es2.Set("APP_PORT", "9090")

	r1, _ := Digest(es1, nil)
	r2, _ := Digest(es2, nil)
	if r1.Digest == r2.Digest {
		t.Fatal("digests should differ after value change")
	}
}

func TestDigest_MissingKey(t *testing.T) {
	es := baseDigestSet()
	_, err := Digest(es, []string{"DOES_NOT_EXIST"})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestDigest_NilEnvSet(t *testing.T) {
	_, err := Digest(nil, nil)
	if err == nil {
		t.Fatal("expected error for nil envset")
	}
}

func TestDigestsMatch_Equal(t *testing.T) {
	a := baseDigestSet()
	b := baseDigestSet()
	match, err := DigestsMatch(a, b, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected digests to match")
	}
}

func TestDigestsMatch_Different(t *testing.T) {
	a := baseDigestSet()
	b := baseDigestSet()
	_ = b.Set("APP_HOST", "remotehost")
	match, err := DigestsMatch(a, b, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if match {
		t.Fatal("expected digests to differ")
	}
}
