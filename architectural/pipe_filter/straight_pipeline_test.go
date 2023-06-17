package pipe_filter

import "testing"

func TestStraightPipeline(t *testing.T) {
	split := NewSplitFilter(",")
	toInt := NewToIntFilter()
	sum := NewSumFilter()
	sp := NewStraightPipeline("p_01", split, toInt, sum)
	ret, err := sp.Process("1,2,3")
	if err != nil {
		t.Fatal(err)
	}
	if ret != 6 {
		t.Fatalf("The expected is 6, but the actual is %d", ret)
	}
	t.Log("Done!")
}
