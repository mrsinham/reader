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

	for i, b := range buf {
		// nothing to do
		// fmt.Println(offsetFound)
		if b != byte(r.search[offsetFound]) {
			res = append(res, b)
			counter++

			// reset
			if offsetFound > 0 {
				// if the offset found was not the last and it didnt match, reseting the changes
				if offsetFound < len(r.search)-1 {
					for t := 1; t < offsetFound-1; t++ {
						fmt.Println(string(res[i-t]), string(r.search[offsetFound-t]))
						res[i-t] = r.search[offsetFound-t]
					}
				}
				offsetFound = 0
			}
			continue
		}

		// fmt.Printf("found %v", string(b))
		// are we at the end ?
		if i == n-1 {
			// yes we need to load to read the len(search)-offsetFound char
			missingBytes := make([]byte, len(r.search)-offsetFound)
			var k int
			k, err = r.underlying.Read(missingBytes)
			if err != nil {
				// can't read the next bytes
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

		var replacementOffset int
		// if replacement is shorter than search
		if offsetFound >= len(r.replacement)-1 {
			replacementOffset = len(r.replacement) - 1
		} else {
			replacementOffset = offsetFound
		}

		// replacing the char
		res = append(res, r.replacement[replacementOffset])
		counter++
		offsetFound++

		// all replacement was done
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
