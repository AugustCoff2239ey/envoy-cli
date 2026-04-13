package envset

import (
	"testing"
)

func TestNewAuditLog(t *testing.T) {
	log := NewAuditLog("myapp")
	if log.Name != "myapp" {
		t.Errorf("expected name 'myapp', got %q", log.Name)
	}
	if len(log.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(log.Entries))
	}
}

func TestAuditLog_Record(t *testing.T) {
	log := NewAuditLog("myapp")
	log.Record(AuditActionSet, "production", "DB_URL", "set DB_URL")
	log.Record(AuditActionDelete, "staging", "OLD_KEY", "removed OLD_KEY")

	if len(log.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(log.Entries))
	}

	if log.Entries[0].Action != AuditActionSet {
		t.Errorf("expected action 'set', got %q", log.Entries[0].Action)
	}
	if log.Entries[0].Key != "DB_URL" {
		t.Errorf("expected key 'DB_URL', got %q", log.Entries[0].Key)
	}
	if log.Entries[1].Environment != "staging" {
		t.Errorf("expected env 'staging', got %q", log.Entries[1].Environment)
	}
}

func TestAuditLog_Filter_ByAction(t *testing.T) {
	log := NewAuditLog("myapp")
	log.Record(AuditActionSet, "production", "KEY1", "set KEY1")
	log.Record(AuditActionDelete, "production", "KEY2", "deleted KEY2")
	log.Record(AuditActionSet, "staging", "KEY3", "set KEY3")

	setEntries := log.Filter(AuditActionSet)
	if len(setEntries) != 2 {
		t.Errorf("expected 2 'set' entries, got %d", len(setEntries))
	}

	delEntries := log.Filter(AuditActionDelete)
	if len(delEntries) != 1 {
		t.Errorf("expected 1 'delete' entry, got %d", len(delEntries))
	}
}

func TestAuditLog_Filter_Empty(t *testing.T) {
	log := NewAuditLog("myapp")
	log.Record(AuditActionImport, "local", "", "imported from file")
	log.Record(AuditActionSync, "production", "", "synced to production")

	all := log.Filter("")
	if len(all) != 2 {
		t.Errorf("expected 2 entries with empty filter, got %d", len(all))
	}
}

func TestAuditLog_Summary(t *testing.T) {
	log := NewAuditLog("myapp")
	log.Record(AuditActionMerge, "staging", "", "merged from local")

	summary := log.Summary()
	expected := "AuditLog[myapp]: 1 entries"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}

func TestAuditLog_TimestampSet(t *testing.T) {
	log := NewAuditLog("myapp")
	log.Record(AuditActionSet, "local", "FOO", "set FOO")

	entry := log.Entries[0]
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
