package secretrabbit

import (
	"reflect"
	"strings"
	"testing"
)

// Most tests borrowed/adapted from https://github.com/dh1tw/gosamplerate

var invalidConverter = Converter(5)

func TestGetConverterName(t *testing.T) {
	got := Linear.String()
	want := "Linear Interpolator"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGetConverterNameError(t *testing.T) {
	got := invalidConverter.String()
	want := "unknown samplerate converter (5)"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGetConverterDescription(t *testing.T) {
	got := Linear.Description()
	want := "Linear interpolator, very fast, poor quality."
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGetConverterDescriptionError(t *testing.T) {
	got := invalidConverter.Description()
	want := "unknown samplerate converter"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGetVersion(t *testing.T) {
	version := Version()
	if !strings.Contains(version, "libsamplerate-") {
		t.Fatalf("got %q, want something containing %q", version, "libsamplerate-")
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

func TestFullNewClose(t *testing.T) {
	src, err := New(SincFastest, 2)
	if err != nil {
		t.Fatal(err)
	}
	src.Close()
	src.Close() // second close should not panic
}

func TestInvalidSrcObject(t *testing.T) {
	_, err := New(invalidConverter, 2)
	want := "Bad converter number"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("got %q, want something containing %q", err.Error(), want)
	}
}

func TestProcessWithEndOfInputFlagSet(t *testing.T) {
	nChannels := 2
	src, err := New(SincFastest, nChannels)
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()

	in := []float32{0.1, -0.5, 0.2, -0.3}
	out := make([]float32, 10)
	nIn, nOut, err := src.Process(in, out, 2.0, true)
	if err != nil {
		t.Fatal(err)
	}
	if want := len(in) / nChannels; want != nIn {
		t.Fatalf("consumed %d frames, want %d", nIn, want)
	}
	if want := 4; want != nOut {
		t.Fatalf("wrote %d frames, want %d", nOut, want)
	}
	out = out[:nOut*nChannels]
	wantOut := []float32{
		0.11488709,
		-0.46334597, 0.18373828, -0.48996875, 0.1821644,
		-0.32879135, 0.10804618, -0.11150829,
	}

	if !reflect.DeepEqual(out, wantOut) {
		t.Log("input:", in)
		t.Log("output:", out)
		t.Logf("expected output: %v", wantOut)
		t.Fatal("unexpected output")
	}
}

func TestProcessErrorWithInvalidRatio(t *testing.T) {
	src, err := New(Linear, 1)
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()

	in := make([]float32, 100)
	out := make([]float32, 100)
	_, _, err = src.Process(in, out, -5, true)
	got := err.Error()
	want := "SRC ratio outside [1/256, 256] range."
	if !strings.Contains(got, want) {
		t.Fatalf("got %q, want something containing %q", got, want)
	}
}
