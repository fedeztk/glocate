package engine

import (
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
		log.Fatal(err)
	}

	zr := lz4.NewReader(f)

	r := make(chan string)

	// read the file line by line and send it to the channel
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := zr.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			paths := strings.Split(string(buf[:n]), "\n")
			for _, path := range paths {
				if path != "" {
					r <- path
				}
			}
		}
		close(r)
		f.Close()
	}()

	return r
}
