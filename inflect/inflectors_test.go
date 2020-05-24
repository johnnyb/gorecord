package inflect

import (
	"testing"
)

func expect(t *testing.T, cond bool, msg string, args ...interface{}) {
	if !cond {
		t.Errorf(msg, args...)
	}
}

func TestComponentize(t *testing.T) {
	vals := Componentize("This IsA strange-looking_string")
	expect(t, len(vals) == 6, "Wrong number of components: %d", len(vals))
	expect(t, vals[0] == "this", "Wrong first component: %s", vals[0])
	expect(t, vals[4] == "looking", "Wrong next-to-last component: %s", vals[4])
	expect(t, vals[5] == "string", "Wrong last component: %s", vals[5])

	vals = Componentize("")
	expect(t, len(vals) == 0, "Should have no components of a zero-length string")
}

func TestUnderscore(t *testing.T) {
	val := Underscore("This is AnotherTest")
	expect(t, val == "this_is_another_test", "Did not underscore string properly (%s)", val)
}

func TestSingularization(t *testing.T) {
	s := Singularize("chairs")
	expect(t, s == "chair", "Wrong singular: %s", s)
}

func TestPluralization(t *testing.T) {
	p := Pluralize("chair")
	expect(t, p == "chairs", "Wrong plural: %s", p)
	p = Pluralize("dog")
	expect(t, p == "dogs", "Wrong plural: %s", p)
	p = Pluralize("the dog")
	expect(t, p == "the dogs", "Wrong plural: %s", p)
	p = Pluralize("Dog")
	expect(t, p == "Dogs", "Wrong plural: %s", p)
	p = Pluralize("The Dog")
	expect(t, p == "The Dogs", "Wrong plural: %s", p)
	p = Pluralize("person")
	expect(t, p == "people", "Wrong plural: %s", p)
	p = Pluralize("dogs")
	expect(t, p == "dogs", "Wrong (already) plural: %s", p)
	p = Pluralize("The Dogs")
	expect(t, p == "The Dogs", "Wrong (already) plural: %s", p)
}
