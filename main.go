package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
)

// TODO(sneha): make backends configurable.
var (
	backends = []string{"cnn.com", "bbc.co.uk", "msn.com"}
)

func main() {

	// TODO(sneha) Validate backends/transform into valid list.

	// Create global client.
	// TODO(sneha): configure with more options.
	client := http.Client{}

	// HTTP handler and server.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Randomly select from list of backends.
		n := rand.Intn(len(backends))

		// TODO(sneha): provide hostname and port, scheme for backends
		// How do real lbs handle this?
		// HTTP client is limiting us but realy want to demarcate
		// a - which host to send the request to vs.
		// b - which host is in the header that we want to maintain
		fmt.Println(r)
		fmt.Println(r.URL.String())
		r.URL.Host = backends[n]
		r.URL.Scheme = "https"
		fmt.Println(r.URL.String())
		req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
		if err != nil {
			// TODO(sneha): fix how this returns later.
			http.Error(w, "cannot process request", http.StatusBadGateway)
			return
		}

		for key, vals := range r.Header {
			for _, val := range vals {
				req.Header.Add(key, val)
			}
		}

		res, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		defer res.Body.Close()

		if res.StatusCode/100 == 5 {
			http.Error(w, fmt.Sprintf("backend returns status code: %s", res.Status), http.StatusBadGateway)
			return
		}

		for key, vals := range res.Header {
			for _, val := range vals {
				w.Header().Add(key, val)
			}
		}

		w.WriteHeader(res.StatusCode)

		// TODO(sneha): Split into ioutil.Readall and therefore be able to
		// differentiate and clearly demarcate what the error is.
		_, err = io.Copy(w, res.Body)
		if err != nil {
			log.Printf("error writing response to client: %v", err)
		}

	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
