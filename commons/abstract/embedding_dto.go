package abstract

import (
	"errors"
	"strings"
	"time"
)

// Embedding ... a named feature vector to be used for similarity search
type Embedding struct {
	Name       string    `json:"name,omitempty"`
	InsertedAt time.Time `json:"inserted_at,omitempty"`
	Vector     []float32 `json:"vector,omitempty"`
}

// Validate ... validate an embedding
func (em *Embedding) Validate() error {
	// the name should not be empty or we may not be able to retrieve the embedding
	if len(strings.TrimSpace(em.Name)) == 0 {
		return errors.New("Embedding Name is undefined")
	}

	return nil
}
