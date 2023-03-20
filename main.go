package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type post struct {
	Method string `json:"method"`
	Params struct {
		Urls string `json:"urls"`
	} `json:"params"`
}

type claim struct {
	Result struct {
		Key struct {
			Address string `json:"address"`

			// useful fields I suppose..
			Meta struct {
				Claims int `json:"claims_in_channel"`
			}

			Name         string
			PermamentURL string `json:"permanent_url"`
			ShortURL     string `json:"short_url"`
			Type         string
			Value        struct {
				Title       string
				Description string
				Thumbnail   struct {
					URL string
				}
			}

			ValueType string
		} `json:"@stavi"`
	}
}

func resolveClaim(p post) (*claim, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(data)
	resp, err := http.Post("http://localhost:5279", "application/json", r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	v := claim{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func main() {
	p := post{Method: "resolve", Params: struct {
		Urls string "json:\"urls\""
	}{"@stavi"}}

	handler := func(w http.ResponseWriter, r *http.Request) {
		claim, err := resolveClaim(p)
		if err != nil {
			panic(err)
		}

		io.WriteString(w, claim.Result.Key.Name)
		io.WriteString(w, "\n")
		io.WriteString(w, claim.Result.Key.Value.Description)
		io.WriteString(w, "\n")
		io.WriteString(w, claim.Result.Key.ShortURL)
		io.WriteString(w, "\n")
		_, err = fmt.Fprintf(w,
			"%s has %d videos currently",
			claim.Result.Key.Name,
			claim.Result.Key.Meta.Claims,
		)
		if err != nil {
			panic(err)
		}
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
