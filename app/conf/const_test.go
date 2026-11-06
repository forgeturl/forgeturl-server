package conf

import "testing"

func TestPageIDValidation(t *testing.T) {
	for _, id := range []string{"", "X2cGtMA", "E", "E../../x", "R0bad", "E2cGtMA ", "E123456789ABCDE"} {
		if ValidPageID(id) {
			t.Errorf("accepted malformed id %q", id)
		}
	}
	for _, id := range []string{"E2cGtMA", "R2cGtMA", "A2cGtMA", "O2cGtMA"} {
		if !ValidPageID(id) {
			t.Errorf("rejected valid id %q", id)
		}
	}
	if ParseIdType("") != InvalidPage || ParseIdType("X2cGtMA") != InvalidPage {
		t.Fatal("invalid IDs must not panic or resolve as owned pages")
	}
}
