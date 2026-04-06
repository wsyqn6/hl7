package hl7

import (
	"context"
	"strings"
	"sync"
	"testing"
)

func TestParseParallel(t *testing.T) {
	ctx := context.Background()
	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
EVN|A01|202401151200|||
PID|1||12345^^^MRN^MR^N||Smith^John^A||19800115|M|||123 Main St^^New York^NY^10001||(555)123-4567|||
PV1|1|I|||Bed 1^Room 100^^Facility 1|||Dr. Johnson^John^M^MD|||ED|||||||||||`)

	msg, err := ParseParallel(ctx, data)
	if err != nil {
		t.Fatalf("ParseParallel failed: %v", err)
	}

	if len(msg.segments) != 4 {
		t.Errorf("expected 4 segments, got %d", len(msg.segments))
	}
}

func TestParseParallelWithWorkers(t *testing.T) {
	ctx := context.Background()
	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
EVN|A01|202401151200|||
PID|1||12345^^^MRN^MR^N||Smith^John^A||19800115|M|||123 Main St^^New York^NY^10001||(555)123-4567|||
PV1|1|I|||Bed 1^Room 100^^Facility 1|||Dr. Johnson^John^M^MD|||ED|||||||||||`)

	msg, err := ParseParallelWithWorkers(ctx, data, 4, 2)
	if err != nil {
		t.Fatalf("ParseParallelWithWorkers failed: %v", err)
	}

	if len(msg.segments) != 4 {
		t.Errorf("expected 4 segments, got %d", len(msg.segments))
	}
}

func TestParseParallelEmpty(t *testing.T) {
	ctx := context.Background()
	_, err := ParseParallel(ctx, []byte{})
	if err == nil {
		t.Error("expected error for empty data")
	}
}

func TestParseParallelContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Create large message that would take time to parse
	var sb strings.Builder
	sb.WriteString("MSH|^~\\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5\r\n")
	for i := 0; i < 100; i++ {
		sb.WriteString("PID|1||12345|||Test||\r\n")
	}

	_, err := ParseParallel(ctx, []byte(sb.String()))
	if err == nil {
		t.Error("expected context error")
	}
}

func TestParseParallelConcurrent(t *testing.T) {
	ctx := context.Background()
	data := []byte(`MSH|^~\&|SENDING|FACILITY|||202401151200||ADT^A01|MSG001|P|2.5
EVN|A01|202401151200|||
PID|1||12345|||Smith||`)

	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			msg, err := ParseParallel(ctx, data)
			if err != nil {
				errors <- err
				return
			}
			if len(msg.segments) != 3 {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("concurrent parse error: %v", err)
	}
}

func TestParseParallelLargeMessage(t *testing.T) {
	ctx := context.Background()

	var sb strings.Builder
	sb.WriteString("MSH|^~\\&|SENDING|FACILITY|||202401151200||ORU^R01|MSG001|P|2.5\r\n")
	sb.WriteString("PID|1||12345|||Test||\r\n")

	for i := 0; i < 100; i++ {
		sb.WriteString("OBX|")
		sb.WriteRune(rune('1' + i%9))
		sb.WriteString("|NM|WBC^WBC^LN||")
		sb.WriteRune(rune('0' + i%10))
		sb.WriteString(".5|x10^3/uL||4.5-11.0|N|||F\r\n")
	}

	msg, err := ParseParallel(ctx, []byte(sb.String()))
	if err != nil {
		t.Fatalf("ParseParallel failed: %v", err)
	}

	expectedSegs := 102 // MSH + PID + 100 OBX
	if len(msg.segments) != expectedSegs {
		t.Errorf("expected %d segments, got %d", expectedSegs, len(msg.segments))
	}
}

func BenchmarkParseParallel(b *testing.B) {
	ctx := context.Background()

	var sb strings.Builder
	sb.WriteString("MSH|^~\\&|SENDING|FACILITY|||202401151200||ORU^R01|MSG001|P|2.5\r\n")
	sb.WriteString("PID|1||12345|||Test||\r\n")

	for i := 0; i < 100; i++ {
		sb.WriteString("OBX|1|NM|WBC^WBC^LN||5.5|x10^3/uL||4.5-11.0|N|||F\r\n")
	}

	data := []byte(sb.String())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg, _ := ParseParallel(ctx, data)
		PutMessage(msg)
	}
}

func BenchmarkParseSerial(b *testing.B) {
	var sb strings.Builder
	sb.WriteString("MSH|^~\\&|SENDING|FACILITY|||202401151200||ORU^R01|MSG001|P|2.5\r\n")
	sb.WriteString("PID|1||12345|||Test||\r\n")

	for i := 0; i < 100; i++ {
		sb.WriteString("OBX|1|NM|WBC^WBC^LN||5.5|x10^3/uL||4.5-11.0|N|||F\r\n")
	}

	data := []byte(sb.String())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Parse(data)
	}
}
