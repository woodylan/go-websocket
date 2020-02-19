package crypto

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	raw := []byte("password123")
	key := []byte("asdf1234qwer7894")
	str, err := Encrypt(raw, key)
	if err == nil {
		t.Log("suc", str)
	} else {
		t.Fatal("fail", err)
	}
}

func TestDncrypt(t *testing.T) {
	raw := "pqjPM0GJUjlgryzMaslqBAzIknumcdgey1MN+ylWHqY="
	key := []byte("asdf1234qwer7894")
	str, err := Decrypt(raw, key)
	if err == nil {
		t.Log("suc", str)
	} else {
		t.Fatal("fail", err)
	}
}

func Benchmark_CBCEncrypt(b *testing.B) {
	b.StopTimer()
	raw := []byte("146576885")
	key := []byte("gdgghf")

	b.StartTimer() //重新开始时间
	for i := 0; i < b.N; i++ {
		_, _ = Encrypt(raw, key)
	}
}

func Benchmark_CBCDecrypt(b *testing.B) {
	b.StopTimer()
	raw := "CDl4Uas8ZyaGXaoOhPZ9NLDvcsIkyvvd++TONd8UPZc="
	key := []byte("gdgghf")

	b.StartTimer() //重新开始时间
	for i := 0; i < b.N; i++ {
		_, _ = Decrypt(raw, key)
	}
}
