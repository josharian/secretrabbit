package secretrabbit

import (
	"reflect"
	"strings"
	"testing"
)

// Most tests borrowed/adapted from https://github.com/dh1tw/gosamplerate

func TestGetConverterName(t *testing.T) {
	got := Linear.String()
	want := "Linear Interpolator"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGetConverterNameError(t *testing.T) {
	got := Converter(5).String()
	want := "unknown samplerate converter"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestSimple(t *testing.T) {
	input := []float32{0.1, -0.5, 0.3, 0.4, 0.1}
	expectedOutput := []float32{0.1, 0.1, -0.10000001, -0.5, 0.033333343, 0.33333334, 0.4, 0.2}

	output, err := Simple(input, 1.5, 1, Linear)
	if err != nil {
		t.Fatal(err)
	}

	if !closeEnough(output, expectedOutput) {
		t.Log("input", input)
		t.Log("output", output)
		t.Log("expectedOutput", expectedOutput)
		t.Fatal("unexpected output")
	}
}

func TestSimpleLessThanOne(t *testing.T) {
	var input []float32
	for i := 0; i < 10; i++ {
		input = append(input, 0.1, -0.5, 0.3, 0.4, 0.1)
	}
	expectedOutput := []float32{0.1, -0.5, 0.4, 0.1, 0.3, 0.1, -0.5, 0.4, 0.1, 0.3, 0.1, -0.5, 0.4, 0.1, 0.3, 0.1, -0.5, 0.4, 0.1, 0.3, 0.1, -0.5, 0.4, 0.1, 0.3}

	output, err := Simple(input, 0.5, 1, Linear)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Log("input", input)
		t.Log("output", output)
		t.Fatal("unexpected output")
	}
}

func TestSimpleError(t *testing.T) {
	input := []float32{0.1, 0.9}
	var invalidRatio float64 = -5.3

	_, err := Simple(input, invalidRatio, 1, Linear)
	if err == nil {
		t.Fatal("expected Error")
	}
	got := err.Error()
	want := "SRC ratio outside [1/256, 256] range."
	if !strings.Contains(got, want) {
		t.Fatalf("got %q, want something containing %q", got, want)
	}
}

func closeEnough(a, b []float32) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v-b[i] > 0.00001 {
			return false
		}
	}
	return true
}
