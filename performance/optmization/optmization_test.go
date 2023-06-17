package optmization

import "testing"

func TestCreateRequest(t *testing.T) {
	str := createRequest()
	t.Log(str)
}

func TestProcessRequestEasyJson(t *testing.T) {
	reqs := make([]string, 0)
	reqs = append(reqs, createRequest())
	reps := precessRequestEasyJson(reqs)
	t.Log(reps[0])
}

/*
easyjson builder make(0,n)
BenchmarkProcessRequest-8                 171765              6582 ns/op           4664 B/op	20 allocs/op
easyjson buff make(n,n)
BenchmarkProcessRequestEasyJson-8        2804962               436.9 ns/op          256 B/op	5 allocs/op
api + make(n,n) i
BenchmarkProcessRequestAPI-8              738747              1557 ns/op            536 B/op	10 allocs/op
easyjson + make(n,n) i
BenchmarkProcessRequestAPI_1-8           2558658               429.7 ns/op          256 B/op	5 allocs/op
api + make(n,n) range
BenchmarkProcessRequestAPI_Range-8        732148              1611 ns/op            536 B/op	10 allocs/op
*/
func BenchmarkProcessRequest(b *testing.B) {
	reqs := make([]string, 0)
	reqs = append(reqs, createRequest())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processRequest(reqs)
	}
	b.StopTimer()
}

func BenchmarkProcessRequestEasyJson(b *testing.B) {
	reqs := make([]string, 0)
	reqs = append(reqs, createRequest())
	//var ret []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = precessRequestEasyJson(reqs)
		//ret = precessRequestEasyJson(reqs)
	}
	b.StopTimer()
	//b.Log(len(ret))
}

func BenchmarkProcessRequestAPI(b *testing.B) {
	reqs := make([]string, 0)
	reqs = append(reqs, createRequest())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processRequestAPI(reqs)
	}
	b.StopTimer()
}

func BenchmarkProcessRequestAPI_1(b *testing.B) {
	reqs := make([]string, 0)
	reqs = append(reqs, createRequest())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processRequestAPI_1(reqs)
	}
	b.StopTimer()
}

func BenchmarkProcessRequestAPI_Range(b *testing.B) {
	reqs := make([]string, 0)
	reqs = append(reqs, createRequest())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processRequestAPI_Range(reqs)
	}
	b.StopTimer()
}
