package mkvparse

import (
	"log"
	"testing"
	"time"
)

type BenchmarkParser struct {
	DefaultHandler
}

func (p *BenchmarkParser) HandleMasterBegin(id ElementID, info ElementInfo) (bool, error) {
	return true, nil
}

func (p *BenchmarkParser) HandleMasterEnd(id ElementID, info ElementInfo) error {
	return nil
}

func (p *BenchmarkParser) HandleString(id ElementID, value string, info ElementInfo) error {
	return nil
}

func (p *BenchmarkParser) HandleInteger(id ElementID, value int64, info ElementInfo) error {
	return nil
}

func (p *BenchmarkParser) HandleFloat(id ElementID, value float64, info ElementInfo) error {
	return nil
}

func (p *BenchmarkParser) HandleDate(id ElementID, value time.Time, info ElementInfo) error {
	return nil
}

func (p *BenchmarkParser) HandleBinary(id ElementID, value []byte, info ElementInfo) error {
	return nil
}

// https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
// This has been added in order to avoid the benchmark being optimized away.
var benchmarkError error

func BenchmarkParse(b *testing.B) {
	var err error
	for n := 0; n < b.N; n++ {
		handler := BenchmarkParser{}
		if err = ParsePath("benchmark-data.mkv", &handler); err != nil {
			log.Fatalf("%v", err)
		}
	}
	benchmarkError = err
}
