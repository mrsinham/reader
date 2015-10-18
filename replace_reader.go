package reader

import (
	"fmt"
	"io"
)

type ReplaceReader struct {
	underlying     io.Reader
	search         []byte
	replacement    []byte
	firstOccurence map[rune]bool
}

func (r *ReplaceReader) Read(p []byte) (int, error) {
	buf := make([]byte, len(p))
	n, err := r.underlying.Read(buf)

	// reducing to the good size
	buf = buf[:n]

	var res []byte
	var counter int
	var offsetFound int

	// testing travis

	for i, b := range buf {
		if b != byte(r.search[offsetFound]) {
			// reset
			if offsetFound > 0 && offsetFound < len(r.search) {
				// if the offset found was not the last and it didnt match, resetting the changes
				if offsetFound < len(r.search)-1 {
					for t := offsetFound; t > 0; t-- {
						res[len(res)-t] = r.search[offsetFound-t]
					}
				}
				offsetFound = 0
			}
			fmt.Println(string(b))
			res = append(res, b)
			counter++

			continue
		}

		// are we at the end ?
		if i == n-1 {
			// yes we need to load to read the len(search)-offsetFound char
			missingBytes := make([]byte, len(r.search)-offsetFound)
			var k int
			k, err = r.underlying.Read(missingBytes)
			if err != nil && io.EOF != nil {
				// can't read the next bytes
				return counter, err
			}
			if k < len(r.search)-1 {
				// no need to read anymore, copy what we received to keep it
				res = append(res, b)
				res = append(res, missingBytes[:k]...)
				// increment of k more read byte + 1
				counter = counter + k + 1
				break
			}

			// we have more char, appending to buf to see if it's ok
			buf = append(buf, missingBytes[:k]...)
		}

		// if replacement is shorter than search
		if offsetFound >= len(r.replacement)-1 {
			continue
		}

		res = append(res, r.replacement[offsetFound])
		counter++
		offsetFound++

		// all replacement were done
		if offsetFound == len(r.search)-1 {
			// replacement is more long than the search
			if offsetFound < len(r.replacement)-1 {
				// adding missing char
				for _, o := range r.replacement[offsetFound:] {
					res = append(res, o)
					counter++
				}
			}

			offsetFound = 0
		}

	}

	copy(p, res)

	return counter, err
}

func NewReplaceReader(underlying io.Reader, search string, replacement string) *ReplaceReader {
	return &ReplaceReader{underlying: underlying, search: []byte(search), replacement: []byte(replacement)}
}
