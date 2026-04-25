package govalid

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Check is safe to call concurrently from multiple goroutines for a fixed
// set of templates. We exercise this with the race detector enabled to
// surface any unsynchronized state access.
// =============================================================================

func Test_Check_ConcurrentSameStruct(t *testing.T) {
	type form struct {
		Name string `valid:"required;username" label:"名称"`
		Mail string `valid:"required;email" label:"邮箱"`
		Age  int    `valid:"min:0;max:120" label:"年龄"`
	}

	bad := form{Name: "1abc", Mail: "not-an-email", Age: 200}
	good := form{Name: "abc", Mail: "user@example.com", Age: 30}

	const goroutines = 16
	const iterations = 200

	var wg sync.WaitGroup
	var goodOk, badOk int64

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				if id%2 == 0 {
					_, ok := Check(good)
					if ok {
						atomic.AddInt64(&goodOk, 1)
					}
				} else {
					_, ok := Check(bad)
					if !ok {
						atomic.AddInt64(&badOk, 1)
					}
				}
			}
		}(g)
	}
	wg.Wait()

	assert.Equal(t, int64(goroutines/2*iterations), goodOk)
	assert.Equal(t, int64(goroutines/2*iterations), badOk)
}

// =============================================================================
// Concurrent Check + completely independent goroutines must not crash or
// surface partially-rendered errors when many distinct structs are
// validated at once.
// =============================================================================

func Test_Check_ConcurrentDifferentStructs(t *testing.T) {
	const goroutines = 32

	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			a := struct {
				A string `valid:"required" label:"a"`
			}{}
			b := struct {
				B int `valid:"min:0;max:10" label:"b"`
			}{B: 100}
			c := struct {
				C []int `valid:"minlen:3" label:"c"`
			}{C: []int{1}}

			for i := 0; i < 50; i++ {
				_, _ = Check(a)
				_, _ = Check(b)
				_, _ = Check(c)
			}
		}(g)
	}
	wg.Wait()
}
