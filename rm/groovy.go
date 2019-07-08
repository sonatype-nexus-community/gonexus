package nexusrm

import (
	"crypto/sha1"
	"encoding/hex"
)

func newAnonGroovyScript(content string) Script {
	h := sha1.New()
	h.Write([]byte(content))

	return Script{hex.EncodeToString(h.Sum(nil)), content, "groovy"}
}
