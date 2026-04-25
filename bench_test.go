package govalid

import "testing"

// =============================================================================
// Benchmarks track performance of the hot paths so future changes can be
// compared with `benchstat`.
// =============================================================================

type benchSimple struct {
	Name string `valid:"required" label:"名称"`
	Age  int    `valid:"min:0;max:120" label:"年龄"`
}

func BenchmarkCheck_Simple_Valid(b *testing.B) {
	v := benchSimple{Name: "iwh", Age: 24}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Check(v)
	}
}

func BenchmarkCheck_Simple_Invalid(b *testing.B) {
	v := benchSimple{Name: "", Age: 200}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Check(v)
	}
}

type benchRich struct {
	Username string `valid:"required;username;minlen:3;maxlen:20"`
	Email    string `valid:"required;email"`
	Mobile   string `valid:"required;mobile"`
	Age      uint   `valid:"required;min:0;max:120"`
	Tags     []string
}

func BenchmarkCheck_Rich_Valid(b *testing.B) {
	v := benchRich{
		Username: "iwh",
		Email:    "i@example.com",
		Mobile:   "13888888888",
		Age:      24,
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Check(v)
	}
}

func BenchmarkParseRules(b *testing.B) {
	const raw = "required;min:0;max:100;list:a,b,c,d,e"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = parseRules(raw)
	}
}
