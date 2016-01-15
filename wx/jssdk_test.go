package wx

import "testing"

type jsSdkTest struct {
	sdk JsSdk
	t   *testing.T
}

func (sdkTest *jsSdkTest) generateNonceStr() {
	expectLen := 32
	str := sdkTest.sdk.GenerateNoncestr(expectLen)
	if len(str) != expectLen {
		sdkTest.t.Fatalf("result string length is not equal %v", expectLen)
	}
}

func TestRun(t *testing.T) {
	sdkTest := &jsSdkTest{
		sdk: NewJsSdk(),
		t:   t,
	}
	sdkTest.generateNonceStr()
}

func Benchmark_GenerateNonceStr(b *testing.B) {
	sdk := NewJsSdk()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sdk.GenerateNoncestr(32)
	}
}
