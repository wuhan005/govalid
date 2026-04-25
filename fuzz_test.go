package govalid

import (
	"reflect"
	"strings"
	"testing"
)

var reflectStringType = reflect.TypeOf("")

// FuzzCheck stress-tests Check with arbitrary tag values to make sure
// no random combination of rules can crash the validator.
//
// Run with:
//
//	go test -run=^$ -fuzz=FuzzCheck -fuzztime=10s
func FuzzCheck(f *testing.F) {
	// Seed with a representative set of rules so the fuzzer starts with
	// realistic shapes and can mutate from there.
	seeds := []string{
		"required",
		"required;min:0;max:100",
		"minlen:5;maxlen:10",
		"alpha",
		"alphanumeric",
		"alphadash",
		"username",
		"email",
		"ipv4",
		"mobile",
		"tel",
		"phone",
		"idcard",
		"equal:Other",
		"list:a,b,c",
		"unknown:foo,bar",
		"required;email",
		"min:abc",
		"maxlen:-1",
		"",
		";;;",
		"::::",
		"required:",
		":foo",
	}
	for _, s := range seeds {
		f.Add(s, "value")
	}

	f.Fuzz(func(t *testing.T, rule, value string) {
		// Drop any null bytes — Go struct tags can't contain them safely
		// and we'd be testing parsing of malformed tags rather than
		// validator behavior.
		if strings.ContainsRune(rule, 0) || strings.ContainsRune(value, 0) {
			t.Skip()
		}
		// Backtick or quote in the tag would terminate it. Skip those
		// so reflect doesn't observe a malformed struct tag.
		if strings.ContainsRune(rule, '`') || strings.ContainsRune(rule, '"') {
			t.Skip()
		}

		// We can't programmatically embed a fuzzed rule string into a
		// struct tag at runtime, so we exercise parseRules + the Checkers
		// map directly with synthetic CheckerContexts. This still covers
		// the full evaluation surface that struct-tag-driven validation
		// would reach.
		rules := parseRules(rule)
		for _, r := range rules {
			fn, ok := Checkers[r.checker]
			if !ok {
				continue
			}

			// Synthesize a real struct so equal/cross-field checkers don't
			// blow up — they need a valid StructValue to look siblings up.
			holder := struct {
				Field string
				Other string
			}{Field: value, Other: value}
			sv := reflect.ValueOf(holder)

			func() {
				defer func() {
					if rec := recover(); rec != nil {
						t.Fatalf("checker %q panicked on params=%v value=%q: %v",
							r.checker, r.params, value, rec)
					}
				}()
				_ = fn(CheckerContext{
					FieldName:   "Field",
					FieldLabel:  "Field",
					FieldValue:  value,
					FieldType:   reflectStringType,
					StructValue: sv,
					Rule:        r,
				})
			}()
		}
	})
}

// FuzzParseRules just hammers the rule parser. Combined with FuzzCheck it
// gives full coverage of the rule grammar.
func FuzzParseRules(f *testing.F) {
	for _, s := range []string{
		"", ";", "::", "a:b:c", "a:b,c,d,e", "x;y;z", "a:1,2;b:3",
		"required", "min:0;max:100", "list:a,b,c",
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("parseRules panicked on %q: %v", raw, rec)
			}
		}()
		_ = parseRules(raw)
	})
}
