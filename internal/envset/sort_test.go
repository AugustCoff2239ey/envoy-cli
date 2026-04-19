package envset

import (
	"strings"
	"testing"
)

func baseSortSet() *EnvSet {
	es, _ := New("sorttest", "local")
	_ = es.Set("ZEBRA", "1")
	_ = es.Set("APPLE", "3")
	_ = es.Set("MANGO", "2")
	return es
}

func TestSort_AscendingByKey(t *testing.T) {
	es := baseSortSet()
	err := Sort(es, DefaultSortOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := SortedKeys(es)
	if keys[0] != "APPLE" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("expected ascending key order, got %v", keys)
	}
}

func TestSort_DescendingByKey(t *testing.T) {
	es := baseSortSet()
	err := Sort(es, SortOptions{Descending: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := SortedKeys(es)
	if keys[0] != "ZEBRA" || keys[1] != "MANGO" || keys[2] != "APPLE" {
		t.Errorf("expected descending key order, got %v", keys)
	}
}

func TestSort_ByValue(t *testing.T) {
	es := baseSortSet()
	err := Sort(es, SortOptions{ByValue: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := SortedKeys(es)
	// values: ZEBRA=1, MANGO=2, APPLE=3
	if keys[0] != "ZEBRA" || keys[1] != "MANGO" || keys[2] != "APPLE" {
		t.Errorf("expected value-ascending order, got %v", keys)
	}
}

func TestSort_SelectedKeys(t *testing.T) {
	es := baseSortSet()
	err := Sort(es, SortOptions{Keys: []string{"ZEBRA", "APPLE"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	order := es.Meta["_sort_order"]
	if !strings.Contains(order, "APPLE") || !strings.Contains(order, "ZEBRA") {
		t.Errorf("expected selected keys in order, got %q", order)
	}
	if strings.Contains(order, "MANGO") {
		t.Errorf("MANGO should not be in selected-key sort order")
	}
}

func TestSort_MissingKey(t *testing.T) {
	es := baseSortSet()
	err := Sort(es, SortOptions{Keys: []string{"NOTEXIST"}})
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestSort_NilEnvSet(t *testing.T) {
	err := Sort(nil, DefaultSortOptions())
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestSortedKeys_FallbackOrder(t *testing.T) {
	es := baseSortSet()
	keys := SortedKeys(es)
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("fallback order not lexicographic: %v", keys)
		}
	}
}
