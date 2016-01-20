package main

import (
	"fmt"
    "net/url"
    "net/http"
    "bytes"
	"errors"
	"io/ioutil"
    "encoding/json"
)

//KnowYourMeme struct for their search JSON API. Has all available attributes.
type KnowYourMeme struct {
	Origin       string  `json:"origin"`
	Body         string  `json:"body"`
	DocID        string  `json:"docid"`
    Name         string  `json:"name"`
    Timestamp    string  `json:"timestamp"`
	Tags         string  `json:"tags"`
	URL          string  `json:"url"`
    IconURL      string  `json:"icon_url"`
	OriginDate   string  `json:"origin_date"`
	_ID          string  `json:"__id"`
	Summary      string  `json:"summary"`
	Status       string  `json:"status"`
	NSFW         string  `json:"nsfw"`
	Year         string  `json:"year"`
	_Type        string  `json:"__type"`
    Category     string  `json:"category"`
    Type         string  `json:"type"`
}

type request struct {
	url       string
    values    *url.Values
	cookie    *http.Cookie
}

//getResponse created by Jimmy Zelinskie in 2012
func (r request) getResponse() (*bytes.Buffer, error) {
	// Determine the HTTP action.
	var action, finalurl string
	if r.values == nil {
		action = "GET"
		finalurl = r.url
	} else {
		action = "POST"
		finalurl = r.url + "?" + r.values.Encode()
	}

	// Create a request and add the proper headers.
	req, err := http.NewRequest(action, finalurl, nil)
	if err != nil {
		return nil, err
	}
	if r.cookie != nil {
		req.AddCookie(r.cookie)
	}

	// Handle the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	respbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(respbytes), nil
}

// String returns the string representation of a submission.
func (h *KnowYourMeme) String() string {
	return fmt.Sprintf("%s - %s (%s)", h.URL, h.Name, h.Body)
}


//knowyourmeme
func getMemes(urlquery string) ([]*KnowYourMeme, error){
    
    req := request{
    url: fmt.Sprintf("http://rkgk.api.searchify.com/v1/indexes/kym_production/instantlinks?query=%s&fetch=*", urlquery),
	}
    //fmt.Printf("Meme: %s \n\n", req.url)
	body, err := req.getResponse()
	if err != nil {
        fmt.Printf("Something went wrong in getResponse: %s \n\n", err)
		return nil, err
	}

	type Response struct {
        Results []*KnowYourMeme
	}
    
    //Decodes JSON into Response struct
	r := new(Response)
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
        fmt.Printf("JSON Decode Error: %s \n\n", err)
		return nil, err
	}

    //Allocates a go array and puts JSON array members into it.
	knowyourmeme := make([]*KnowYourMeme, len(r.Results))
	for i, result := range r.Results {
		knowyourmeme[i] = result
	}
    return knowyourmeme, nil
}