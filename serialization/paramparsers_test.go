package serialization

import "testing"

func FuzzParsestring(f *testing.F) {
	f.Add([]byte("some input data"))

	f.Fuzz(func(t *testing.T, input []byte) {
		v, e := Parsestring(input)
		if e != nil {
			t.Errorf("string parsing value: %v, error: %v", v, e)
		}
	})
}

func FuzzParseint(f *testing.F) {
	f.Add([]byte("894375893234"))

	f.Fuzz(func(t *testing.T, input []byte) {
		_, e := Parseint(input)
		if e != nil {

		}
	})
}
