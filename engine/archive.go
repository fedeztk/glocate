package engine

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/opencoff/go-walk"
	"github.com/pierrec/lz4"
)

// compressPipe all entries from the channel to a lz4 file separated by \n
func compressPipe(r chan walk.Result, w io.Writer) error {
	zw := lz4.NewWriter(w)
	defer zw.Close()

	for path := range r {
		_, err := zw.Write([]byte(path.Path + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

// decompressPipe reads the lz4 file and sends each line to the channel
func decompressPipe() chan string {
	f, err := os.Open(dbFileName)
	if err != nil {
		log.Fatalf("%s. Please run glocate -i to create the database", err)
	}

	zr := lz4.NewReader(f)

	r := make(chan string)

	go func() {
		buf := bytes.Buffer{}
		_, err := io.Copy(&buf, zr)
		if err != nil {
			log.Fatal(err)
		}

		for _, path := range strings.Split(buf.String(), "\n") {
			if path != "" {
				r <- path
			}
		}
		close(r)
		f.Close()
	}()

	return r
}
