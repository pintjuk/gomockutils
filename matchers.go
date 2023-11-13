package gomockutils

import (
	"fmt"
	"reflect"
	"regexp"

	"go.uber.org/mock/gomock"
)

var (
	patField  = regexp.MustCompile(`^\.(\w+)`)
	patForall =  regexp.MustCompile(`^\[\]`)
)


type filedMatcher struct {
	field   string
	matcher gomock.Matcher
}

// matches if x is an struct with that has a field f.field that matches f.matcher
func (f filedMatcher) Matches(x interface{}) bool {
	value := reflect.ValueOf(x)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return false
	}
	targetField := value.FieldByName(f.field)
	if !targetField.IsValid() {
		return false

	}
	return f.matcher.Matches(targetField.Interface())

}
func (f filedMatcher) String() string {
	return fmt.Sprintf(".%v %v", f.field, f.matcher.String())
}

// Field is a matcher that matches structs that contain a field, that matches matcher
func Field(field string, matcher gomock.Matcher) gomock.Matcher {
	return filedMatcher{field, matcher}
}

type forEachMatcher struct {
	matcher gomock.Matcher
}

func (f forEachMatcher) Matches(x interface{}) bool {
	value := reflect.ValueOf(x)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return false
	}

	for i := 0; i < value.Len(); i++ {

		if !f.matcher.Matches(value.Index(i).Interface()) {
			return false
		}
	}

	return true

}

func (f forEachMatcher) String() string {
	return fmt.Sprintf("[*] %v", f.matcher.String())
}

// ForEach is a matcher that matches lists in witch
// every element matches the suplied matcher
func ForEach(matcher gomock.Matcher) gomock.Matcher {
	return forEachMatcher{matcher: matcher}
}

// Query will apply the provided matcher to the nested struct field specified by a path
//
// For example StructQuery(".delivery[].addons", Eq([]string{"H9"}) will match any struct that has a field deliviery that is a list,
// in witch each element contains labels taht are qual to a list with a single "H9" string.
// This is the same as Field("delivery", ForEach(Filed("addons", Eq([]string{"H9"}))), but more readable
func SQuery(path string, matcher gomock.Matcher) gomock.Matcher {
	if len(path) == 0 {
		return matcher
	}

	if patField.MatchString(path) {
		indices := patField.FindStringIndex(path)
		fmt.Printf("is field %v, %v, %v\n", indices, path[indices[0]:indices[1]], path[indices[1]:])
		return Field(
			path[indices[0]+1:indices[1]],
			SQuery(path[indices[1]:], matcher),
		)
	} else if patForall.MatchString(path) {
		indices := patForall.FindStringIndex(path)
		return ForEach(
			SQuery(path[indices[1]:], matcher),
		)
	} else {
		panic(fmt.Sprintf("Invalid query: When parsing %v, Expected . or [ found %c", path, path[0]))
	}
}
