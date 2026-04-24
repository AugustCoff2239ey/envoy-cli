package envset

import (
	"testing"
)

func baseSignatureSet() *EnvSet {
	es, _ := New("sigtest", "staging")
	_ = es.Set("DB_HOST", "localhost")
	_ = es.Set("DB_PORT", "5432")
	_ = es.Set("API_KEY", "secret")
	return es
}

func TestSign_Valid(t *testing.T) {
	es := baseSignatureSet()
	if err := Sign(es, "passphrase123"); err != nil {
		t.Fatalf("Sign returned unexpected error: %v", err)
	}
	if !HasSignature(es) {
		t.Fatal("expected signature to be present after Sign")
	}
}

func TestVerifySignature_Valid(t *testing.T) {
	es := baseSignatureSet()
	_ = Sign(es, "passphrase123")
	if err := VerifySignature(es, "passphrase123"); err != nil {
		t.Fatalf("VerifySignature returned unexpected error: %v", err)
	}
}

func TestVerifySignature_WrongPassphrase(t *testing.T) {
	es := baseSignatureSet()
	_ = Sign(es, "passphrase123")
	err := VerifySignature(es, "wrongpass")
	if err == nil {
		t.Fatal("expected error for wrong passphrase, got nil")
	}
	if err != ErrSignatureMismatch {
		t.Fatalf("expected ErrSignatureMismatch, got %v", err)
	}
}

func TestVerifySignature_TamperedValue(t *testing.T) {
	es := baseSignatureSet()
	_ = Sign(es, "passphrase123")
	_ = es.Set("DB_HOST", "evil-host")
	err := VerifySignature(es, "passphrase123")
	if err != ErrSignatureMismatch {
		t.Fatalf("expected ErrSignatureMismatch after tamper, got %v", err)
	}
}

func TestVerifySignature_NotFound(t *testing.T) {
	es := baseSignatureSet()
	err := VerifySignature(es, "passphrase123")
	if err != ErrSignatureNotFound {
		t.Fatalf("expected ErrSignatureNotFound, got %v", err)
	}
}

func TestSign_EmptyPassphrase(t *testing.T) {
	es := baseSignatureSet()
	if err := Sign(es, ""); err != ErrEmptyPassphrase {
		t.Fatalf("expected ErrEmptyPassphrase, got %v", err)
	}
}

func TestSign_NilEnvSet(t *testing.T) {
	if err := Sign(nil, "pass"); err == nil {
		t.Fatal("expected error for nil envset")
	}
}

func TestClearSignature(t *testing.T) {
	es := baseSignatureSet()
	_ = Sign(es, "passphrase123")
	ClearSignature(es)
	if HasSignature(es) {
		t.Fatal("expected signature to be absent after ClearSignature")
	}
}

func TestSign_ReplacesPreviousSignature(t *testing.T) {
	es := baseSignatureSet()
	_ = Sign(es, "first")
	sig1 := es.Meta[signatureMetaKey]
	_ = Sign(es, "second")
	sig2 := es.Meta[signatureMetaKey]
	if sig1 == sig2 {
		t.Fatal("expected different signatures for different passphrases")
	}
}
