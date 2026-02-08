package server

import (
	"testing"

	"github.com/thewh1teagle/sona/internal/whisper"
)

func TestCsToSeconds(t *testing.T) {
	tests := []struct {
		cs   int64
		want float64
	}{
		{0, 0.0},
		{100, 1.0},
		{250, 2.5},
		{6000, 60.0},
	}
	for _, tt := range tests {
		got := csToSeconds(tt.cs)
		if got != tt.want {
			t.Errorf("csToSeconds(%d) = %f, want %f", tt.cs, got, tt.want)
		}
	}
}

func TestCsToSRTTime(t *testing.T) {
	tests := []struct {
		cs   int64
		want string
	}{
		{0, "00:00:00,000"},
		{250, "00:00:02,500"},
		{6150, "00:01:01,500"},
		{360000, "01:00:00,000"},
	}
	for _, tt := range tests {
		got := csToSRTTime(tt.cs)
		if got != tt.want {
			t.Errorf("csToSRTTime(%d) = %q, want %q", tt.cs, got, tt.want)
		}
	}
}

func TestCsToVTTTime(t *testing.T) {
	tests := []struct {
		cs   int64
		want string
	}{
		{0, "00:00:00.000"},
		{250, "00:00:02.500"},
	}
	for _, tt := range tests {
		got := csToVTTTime(tt.cs)
		if got != tt.want {
			t.Errorf("csToVTTTime(%d) = %q, want %q", tt.cs, got, tt.want)
		}
	}
}

func TestFormatSRT(t *testing.T) {
	segments := []whisper.Segment{
		{Start: 0, End: 250, Text: " Hello world"},
		{Start: 250, End: 510, Text: " How are you"},
	}
	got := formatSRT(segments)
	want := "1\n00:00:00,000 --> 00:00:02,500\nHello world\n\n2\n00:00:02,500 --> 00:00:05,100\nHow are you\n"
	if got != want {
		t.Errorf("formatSRT() =\n%q\nwant:\n%q", got, want)
	}
}

func TestFormatVTT(t *testing.T) {
	segments := []whisper.Segment{
		{Start: 0, End: 250, Text: " Hello world"},
		{Start: 250, End: 510, Text: " How are you"},
	}
	got := formatVTT(segments)
	want := "WEBVTT\n\n00:00:00.000 --> 00:00:02.500\nHello world\n\n00:00:02.500 --> 00:00:05.100\nHow are you\n"
	if got != want {
		t.Errorf("formatVTT() =\n%q\nwant:\n%q", got, want)
	}
}

func TestBuildVerboseJSON(t *testing.T) {
	segments := []whisper.Segment{
		{Start: 0, End: 250, Text: "Hello"},
		{Start: 250, End: 510, Text: " world"},
	}
	v := buildVerboseJSON(segments, nil)
	if v.Text != "Hello world" {
		t.Errorf("Text = %q, want %q", v.Text, "Hello world")
	}
	if len(v.Segments) != 2 {
		t.Fatalf("got %d segments, want 2", len(v.Segments))
	}
	if v.Segments[0].Start != 0.0 || v.Segments[0].End != 2.5 {
		t.Errorf("segment[0] times = (%f, %f), want (0.0, 2.5)", v.Segments[0].Start, v.Segments[0].End)
	}
}

func TestParseBoolFormValue(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"false", false},
		{"1", true},
		{"0", false},
		{"", false},
		{"yes", false},
	}
	for _, tt := range tests {
		got := parseBoolFormValue(tt.input)
		if got != tt.want {
			t.Errorf("parseBoolFormValue(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
