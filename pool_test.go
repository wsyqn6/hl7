package hl7

import (
	"runtime"
	"sync"
	"testing"
)

func TestGetPutMessage(t *testing.T) {
	msg := GetMessage()
	if msg == nil {
		t.Fatal("GetMessage returned nil")
	}
	if cap(msg.segments) != defaultSegmentCapacity {
		t.Errorf("expected segment capacity %d, got %d", defaultSegmentCapacity, cap(msg.segments))
	}

	msg.segments = append(msg.segments, Segment{name: "MSH", fields: []string{""}})
	PutMessage(msg)

	msg2 := GetMessage()
	if len(msg2.segments) != 0 {
		t.Errorf("expected empty segments after reset, got %d", len(msg2.segments))
	}
	PutMessage(msg2)
}

func TestMessageBuilder(t *testing.T) {
	builder := NewMessageBuilder()
	builder.AddSegment("MSH", "|", "^~\\&", "SEND", "RECV")
	builder.AddSegment("PID", "1", "12345", "DOE^JOHN")

	msg := builder.Build()
	if len(msg.segments) != 2 {
		t.Errorf("expected 2 segments, got %d", len(msg.segments))
	}
	if msg.segments[0].name != "MSH" {
		t.Errorf("expected first segment MSH, got %s", msg.segments[0].name)
	}
	if msg.segments[1].name != "PID" {
		t.Errorf("expected second segment PID, got %s", msg.segments[1].name)
	}

	builder.Release()
}

func TestMessageBuilderWithDelimiters(t *testing.T) {
	delims := Delimiters{
		Field:        '#',
		Component:    '*',
		Repetition:   '~',
		Escape:       '\\',
		SubComponent: '&',
	}

	builder := NewMessageBuilder()
	builder.WithDelimiters(delims)

	msg := builder.Build()
	if msg.delims.Field != '#' {
		t.Errorf("expected field delimiter #, got %c", msg.delims.Field)
	}
	if msg.delims.Component != '*' {
		t.Errorf("expected component delimiter *, got %c", msg.delims.Component)
	}

	builder.Release()
}

func TestSegmentPool(t *testing.T) {
	seg := getSegment()
	if seg == nil {
		t.Fatal("getSegment returned nil")
	}

	// Pool may return a recycled segment with different capacity
	// Just verify it's a valid segment
	if cap(seg.fields) < 4 {
		t.Errorf("expected field capacity at least 4, got %d", cap(seg.fields))
	}

	seg.name = "TEST"
	seg.fields = append(seg.fields, "field1", "field2")
	putSegment(seg)

	seg2 := getSegment()
	if seg2.name != "" {
		t.Errorf("expected empty name after reset, got %s", seg2.name)
	}
	if len(seg2.fields) != 0 {
		t.Errorf("expected empty fields after reset, got %d", len(seg2.fields))
	}
	putSegment(seg2)
}

func BenchmarkMessagePool(b *testing.B) {
	var msg *Message
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		msg = GetMessage()
		msg.segments = append(msg.segments, Segment{name: "MSH", fields: []string{"a", "b"}})
		msg.segments = append(msg.segments, Segment{name: "PID", fields: []string{"1", "2", "3"}})
		PutMessage(msg)
	}

	runtime.KeepAlive(msg)
}

func BenchmarkMessageAlloc(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		msg := NewMessage()
		msg.segments = append(msg.segments, Segment{name: "MSH", fields: []string{"a", "b"}})
		msg.segments = append(msg.segments, Segment{name: "PID", fields: []string{"1", "2", "3"}})
	}

	b.StopTimer()
}

func BenchmarkMessageBuilder(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		builder := NewMessageBuilder()
		builder.AddSegment("MSH", "|", "^~\\&", "SEND", "RECV")
		builder.AddSegment("PID", "1", "12345", "DOE^JOHN")
		msg := builder.Build()
		_ = msg
		builder.Release()
	}
}

func TestPoolConcurrent(t *testing.T) {
	const numGoroutines = 100
	const numIterations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				msg := GetMessage()
				for k := 0; k < 10; k++ {
					seg := getSegment()
					seg.name = "SEG"
					msg.segments = append(msg.segments, *seg)
				}
				PutMessage(msg)
			}
		}()
	}

	wg.Wait()
}

func TestPoolNilHandling(t *testing.T) {
	PutMessage(nil)
	putSegment(nil)
}
