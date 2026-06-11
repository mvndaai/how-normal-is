package domain

import "testing"

func TestNormalizeTraitName(t *testing.T) {
	if got, want := NormalizeTraitName("  Blue Eyes  "), "blue-eyes"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestValidateExclusiveTraitGroups(t *testing.T) {
	selected := []Trait{{Name: "Blue eyes", Group: "eye_color"}, {Name: "Brown eyes", Group: "eye_color"}, {Name: "Woman", Group: "gender"}}
	conflicts := ValidateExclusiveTraitGroups(selected)
	if len(conflicts) != 1 || conflicts[0] != "eye_color" {
		t.Fatalf("unexpected conflicts: %#v", conflicts)
	}
}
