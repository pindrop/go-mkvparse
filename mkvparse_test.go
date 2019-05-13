package mkvparse

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"
)

func TestReadVarInt_2Encodings(t *testing.T) {
	testEncodings := [][]byte{
		{0x82},
		{0x40, 0x02},
		{0x20, 0x00, 0x02},
		{0x10, 0x00, 0x00, 0x02},
	}
	for _, encoding := range testEncodings {
		encoding = append(encoding, 0xde, 0xad, 0xbe, 0xef)
		reader := bytes.NewReader(encoding)
		result, count, err := readVarInt(reader)
		if err != nil {
			t.Errorf("%x: %v", encoding, err)
		}
		if count != int64(len(encoding))-4 {
			t.Errorf("%x: %d != %d", encoding, count, len(encoding)-4)
		}
		if result != 0x2 {
			t.Errorf("%x: %d != %d", encoding, result, 10)
		}
	}
}

func TestReadElementID(t *testing.T) {
	testIDs := map[ElementID][]byte{
		TimecodeElement:        {0xE7},
		EBMLVersionElement:     {0x42, 0x86},
		DefaultDurationElement: {0x23, 0xE3, 0x83},
		EBMLElement:            {0x1A, 0x45, 0xDF, 0xA3},
	}
	for id, encoding := range testIDs {
		encoding = append(encoding, 0xde, 0xad, 0xbe, 0xef)
		reader := bytes.NewReader(encoding)
		result, count, err := readElementID(reader)
		if err != nil {
			t.Errorf("%x: %v", encoding, err)
		}
		if count != int64(len(encoding))-4 {
			t.Errorf("%x: %d != %d", encoding, count, len(encoding)-4)
		}
		if result != id {
			t.Errorf("%x: %x != %x", encoding, result, id)
		}
	}
}

//////////////////////////////////////////////////////////////////////

type ParseEvent struct {
	id    ElementID
	info  ElementInfo
	value interface{}
}

type MasterBeginEvent struct{}
type MasterEndEvent struct{}

type ParseHandler struct {
	DefaultHandler

	events []ParseEvent
}

func (p *ParseHandler) HandleMasterBegin(id ElementID, info ElementInfo) (bool, error) {
	fmt.Printf("%s\n", NameForElementID(id))
	p.events = append(p.events, ParseEvent{id, info, MasterBeginEvent{}})
	return true, nil
}

func (p *ParseHandler) HandleMasterEnd(id ElementID, info ElementInfo) error {
	p.events = append(p.events, ParseEvent{id, info, MasterEndEvent{}})
	return nil
}

func (p *ParseHandler) HandleString(id ElementID, value string, info ElementInfo) error {
	p.events = append(p.events, ParseEvent{id, info, value})
	return nil
}

func (p *ParseHandler) HandleInteger(id ElementID, value int64, info ElementInfo) error {
	p.events = append(p.events, ParseEvent{id, info, value})
	return nil
}

func (p *ParseHandler) HandleFloat(id ElementID, value float64, info ElementInfo) error {
	p.events = append(p.events, ParseEvent{id, info, value})
	return nil
}

func (p *ParseHandler) HandleDate(id ElementID, value time.Time, info ElementInfo) error {
	p.events = append(p.events, ParseEvent{id, info, value})
	return nil
}

func (p *ParseHandler) HandleBinary(id ElementID, value []byte, info ElementInfo) error {
	p.events = append(p.events, ParseEvent{id, info, value})
	return nil
}

func (p *ParseHandler) HandleParseError(reader io.Reader, err *ParseError) *ParseError {
	_, ok := err.Err.(*InvalidElementError)
	if ok {
		return nil
	}

	return err
}

type ParseTest struct {
	data  []byte
	event ParseEvent
}

func TestParseElement(t *testing.T) {
	tests := map[string]ParseTest{
		"time before millenium": {
			[]byte{0x44, 0x61, 0x88, 0xf6, 0xd3, 0xc2, 0xb9, 0x1b, 0xee, 0x28, 0x00},
			ParseEvent{
				DateUTCElement,
				ElementInfo{
					Offset: 3,
					Size:   8,
					Level:  0,
				},
				time.Date(1980, time.January, 21, 21, 03, 0, 0, time.UTC),
			},
		},
	}

	for name, test := range tests {
		reader := bytes.NewReader(test.data)

		handler := ParseHandler{}

		err := Parse(reader, &handler)
		if err != nil {
			t.Errorf("%s: %v", name, err)
		}

		if len(handler.events) != 1 {
			t.Errorf("%s: Invalid event count: %d", name, len(handler.events))
		}

		if test.event != handler.events[0] {
			t.Errorf("%s: Invalid event: %v != %v", name, test.event, handler.events[0])
		}
	}
}
