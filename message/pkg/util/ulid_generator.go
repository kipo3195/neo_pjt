package util

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"io"
	mathRand "math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

type ULIDGenerator struct {
	mu      sync.Mutex
	entropy io.Reader
}

func NewULIDGenerator() (*ULIDGenerator, error) {
	var seedBytes [8]byte
	if _, err := cryptoRand.Read(seedBytes[:]); err != nil {
		return nil, err
	}
	seed := int64(binary.LittleEndian.Uint64(seedBytes[:]))
	src := mathRand.New(mathRand.NewSource(seed))
	monotonic := ulid.Monotonic(src, 0)
	return &ULIDGenerator{entropy: monotonic}, nil
}

func (g *ULIDGenerator) New() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	id, _ := ulid.New(ulid.Timestamp(time.Now().UTC()), g.entropy)
	return id.String()
}
