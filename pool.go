package hl7

import (
	"sync"
)

const (
	defaultSegmentCapacity = 32
	defaultFieldCapacity   = 64
)

var messagePool = sync.Pool{
	New: func() interface{} {
		return &Message{
			segments: make([]Segment, 0, defaultSegmentCapacity),
			delims:   DefaultDelimiters(),
		}
	},
}

var segmentPool = sync.Pool{
	New: func() interface{} {
		return &Segment{
			fields: make([]string, 0, defaultFieldCapacity),
		}
	},
}

func GetMessage() *Message {
	msg := messagePool.Get().(*Message)
	msg.segments = msg.segments[:0]
	msg.delims = DefaultDelimiters()
	return msg
}

func PutMessage(msg *Message) {
	if msg == nil {
		return
	}
	for i := range msg.segments {
		putSegment(&msg.segments[i])
	}
	msg.segments = msg.segments[:0]
	messagePool.Put(msg)
}

func getSegment() *Segment {
	return segmentPool.Get().(*Segment)
}

func putSegment(seg *Segment) {
	if seg == nil {
		return
	}
	seg.name = ""
	seg.fields = seg.fields[:0]
	segmentPool.Put(seg)
}

type MessageBuilder struct {
	msg *Message
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		msg: GetMessage(),
	}
}

func (b *MessageBuilder) WithDelimiters(d Delimiters) *MessageBuilder {
	b.msg.delims = d
	return b
}

func (b *MessageBuilder) AddSegment(name string, fields ...string) *MessageBuilder {
	seg := getSegment()
	seg.name = name
	seg.fields = append(seg.fields[:0], fields...)
	b.msg.segments = append(b.msg.segments, *seg)
	return b
}

func (b *MessageBuilder) Build() *Message {
	return b.msg
}

func (b *MessageBuilder) Release() {
	PutMessage(b.msg)
	b.msg = nil
}
