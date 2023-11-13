package gomockutils

import (
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestSQuery_Should_EquelNestedFieldMatcheres_When_ValidPathIsProvided(t *testing.T) {
	// Act
	matcher := SQuery(".Deliveries[].Addons[].Code", gomock.Eq("12"))

	//Assert
	expected := filedMatcher{
		field: "Deliveries",
		matcher: forEachMatcher{
			matcher: filedMatcher{
				field: "Addons",
				matcher: forEachMatcher{
					matcher: filedMatcher{
						field:   "Code",
						matcher: gomock.Eq("12")}}}}}

	if matcher != expected {
		t.Fatalf("Squery generated matcher should return %#v, but was %#v", expected, matcher)
	}
	fmt.Printf("%#v", matcher)
}

// A think paniking in this case is good because we will only use this in tests,
// we want the tests to fail eraly if an invalid path is provided,
// and we do not want to check errors of util functios in tests
func TestSQuery_Should_Panic_When_InvalidPathIsProvided(t *testing.T) {
	//Assert
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	// Act
	SQuery("Deliveries[].Addons[].Code", gomock.Eq("12"))

}

