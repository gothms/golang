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

//func BenchmarkProcessRequestAPI(b *testing.B) {
//	reqs := make([]string, 0)
//	reqs = append(reqs, createRequest())
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_ = processRequestAPI(reqs)
//	}
//	b.StopTimer()
//}

func BenchmarkProcessRequestAPI_1(b *testing.B) {
	reqs := make([]string, 0)
	reqs = append(reqs, createRequest())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processRequestAPI_1(reqs)
	}
	b.StopTimer()
}

//func BenchmarkProcessRequestAPI_Range(b *testing.B) {
//	reqs := make([]string, 0)
//	reqs = append(reqs, createRequest())
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_ = processRequestAPI_Range(reqs)
//	}
//	b.StopTimer()
//}
