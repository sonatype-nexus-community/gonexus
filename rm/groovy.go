package nexusrm

import (
	"crypto/sha1"
	"encoding/hex"
)

func newAnonGroovyScript(content string) (s Script) {
	h := sha1.New()
	h.Write([]byte(content))

	s.Name = hex.EncodeToString(h.Sum(nil))
	s.Content = content
	s.Type = "groovy"
	return
}
