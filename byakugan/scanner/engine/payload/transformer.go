package payload

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// Transformer handles payload transformations
type Transformer struct {
	encoders map[string]EncoderFunc
}

// EncoderFunc defines a function that encodes a payload
type EncoderFunc func(string) string

// NewTransformer creates a new payload transformer
func NewTransformer() *Transformer {
	t := &Transformer{
		encoders: make(map[string]EncoderFunc),
	}

	// Register default encoders
	t.RegisterEncoder("base64", base64Encode)
	t.RegisterEncoder("url", urlEncode)
	t.RegisterEncoder("hex", hexEncode)

	return t
}

// RegisterEncoder adds a new encoder
func (t *Transformer) RegisterEncoder(name string, encoder EncoderFunc) {
	t.encoders[name] = encoder
}

// Transform applies transformations to a payload
func (t *Transformer) Transform(p *Payload, transformations []string) error {
	value := p.Value

	for _, trans := range transformations {
		encoder, exists := t.encoders[trans]
		if !exists {
			return fmt.Errorf("unknown transformation: %s", trans)
		}
		value = encoder(value)
	}

	p.Value = value
	p.Encoded = len(transformations) > 0
	return nil
}

// Built-in encoders
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func urlEncode(s string) string {
	return url.QueryEscape(s)
}

func hexEncode(s string) string {
	var sb strings.Builder
	for _, c := range s {
		sb.WriteString(fmt.Sprintf("\\x%02x", c))
	}
	return sb.String()
}
