// Package general defines general utils.
package general

import (
	"io"
	"log"
)

// CloseWithCheck safe closer.
func CloseWithCheck(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Println("Error in close", "err", err)
	}
}
