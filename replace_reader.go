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

	var res []byte = make([]byte, 0)
	var counter int
	var offsetFound int

	// testing travis

	// fmt.Println(string(r.search[6]))

	// os.Exit(-1)

	for i, b := range buf {

		// if offsetFound == len(r.search) {
		// 	if offsetFound < len(r.replacement)-1 {
		// 		// adding missing char
		// 		for _, o := range r.replacement[offsetFound:] {
		// 			res = append(res, o)
		// 			counter++
		// 		}
		// 	}
		// 	offsetFound = 0
		// 	// continue
		// }

		// if b != r.search[offsetFound] {
		// 	res = append(res, b)
		// 	counter++
		// 	continue
		// }

		// res = append(res, r.replacement[offsetFound])
		// offsetFound++
		// counter++

		// fmt.Println(string(b), string(r.search[offsetFound]), offsetFound)
		// fmt.Println(offsetFound)
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

		if b != r.search[offsetFound] {
			// fmt.Println(string(b))
			res = append(res, b)
			counter++

			// reset
			if offsetFound > 0 {
				if offsetFound == len(r.search)-1 {
					// if the offset found was not the last and it didnt match, resetting the changes
					fmt.Println("toto")
					for t := offsetFound; t >= 0; t-- {
						res[len(res)-1-t] = r.search[offsetFound-t]
					}
				}
				offsetFound = 0
			}

			// fmt.Println(b, r.search[offsetFound])

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
		if offsetFound == len(r.replacement)-1 {
			continue
		}

		// fmt.Println(string(r.replacement[offsetFound]))
		res = append(res, r.replacement[offsetFound])
		counter++

		offsetFound++

	}

	copy(p, res)

	return counter, err
}

func NewReplaceReader(underlying io.Reader, search string, replacement string) *ReplaceReader {
	return &ReplaceReader{underlying: underlying, search: []byte(search), replacement: []byte(replacement)}
}
