package struc

import (
	"bytes"
	"testing"
)

type stringArray struct {
	Names [32]string `struc:"[32][256]byte"`
}

func TestPackStringArray(t *testing.T) {
	s := stringArray{}
	for i := 0; i < 32; i++ {
		s.Names[i] = "name" + string(rune('A'+i))
	}

	buf := &bytes.Buffer{}
	err := Pack(buf, &s)
	if err != nil {
		t.Fatalf("Failed to pack: %v", err)
	}

	expected := 32 * 256
	if buf.Len() != expected {
		t.Errorf("Expected buffer length %d, got %d", expected, buf.Len())
	}

	s2 := stringArray{}
	err = Unpack(bytes.NewReader(buf.Bytes()), &s2)
	if err != nil {
		t.Fatalf("Failed to unpack: %v", err)
	}

	for i, name := range s.Names {
		expected := name
		if len(name) > 256 {
			expected = name[:256]
		}
		actual := s2.Names[i]
		if idx := bytes.IndexByte([]byte(actual), 0); idx >= 0 {
			actual = actual[:idx]
		}
		if actual != expected {
			t.Errorf("String %d: expected '%s', got '%s'", i, expected, actual)
		}
	}
}

func BenchmarkStringArrayPack(b *testing.B) {
	s := stringArray{}
	for i := 0; i < 32; i++ {
		s.Names[i] = "name" + string(rune('A'+i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := &bytes.Buffer{}
		if err := Pack(buf, &s); err != nil {
			b.Fatalf("Failed to pack: %v", err)
		}
	}
}

func BenchmarkStringArrayUnpack(b *testing.B) {
	s := stringArray{}
	for i := 0; i < 32; i++ {
		s.Names[i] = "name" + string(rune('A'+i))
	}

	buf := &bytes.Buffer{}
	if err := Pack(buf, &s); err != nil {
		b.Fatalf("Failed to pack: %v", err)
	}
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s2 := stringArray{}
		if err := Unpack(bytes.NewReader(data), &s2); err != nil {
			b.Fatalf("Failed to unpack: %v", err)
		}
	}
}
