package ioutil

import "io"

func WriteAll(w io.Writer, p []byte) error {
	for n := 0; n < len(p); {
		nn, err := w.Write(p[n:])
		n += nn
		if err != nil {
			return err
		}
	}
	return nil
}
