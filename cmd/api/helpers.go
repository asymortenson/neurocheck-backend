package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"vkparser.com/internal/data"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func(app *application) readServerJSON(w http.ResponseWriter, res *http.Response, dst interface{}) error {
	
	dec := json.NewDecoder(res.Body)
	err := dec.Decode(dst)

	fmt.Println(err)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmrshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmrshalTypeError): 
			if unmrshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmrshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmrshalTypeError.Offset)
		case errors.Is(err, io.EOF): 
			return nil
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
		
		}

		err = dec.Decode(&struct{}{})
		if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
		}
	return nil

}



func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmrshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmrshalTypeError): 
			if unmrshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmrshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmrshalTypeError.Offset)
		case errors.Is(err, io.EOF): 
			return errors.New("body must not be empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: uknown field")
			return fmt.Errorf("body contains uknown key %s", fieldName)
		case err.Error() == "http: request body too large": 
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		default:
			return err
		}
		
		}

		err = dec.Decode(&struct{}{})
		if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
		}


		return nil
}

func(app *application) fetchPosts(w http.ResponseWriter, r *http.Request, url string, requestData data.RequestData, offset int32) (interface{}, error) {
	var posts data.Response

	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}	
	
	defer res.Body.Close()

	
	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	err = app.readServerJSON(w, res, &posts)

	if err != nil {
		return nil, err
	}

	var comments []interface{}
	var comment data.Response

	for _, post := range posts.Response.Items {
		url := fmt.Sprintf("https://api.vk.com/method/wall.getComments?owner_id=%d&post_id=%d&offset=%d&access_token=%s&v=5.131", post.OwnerId, post.Id, offset, requestData.AccessToken)
		res, err := http.Get(url)

		if err != nil {
			return nil, err
		}	

		defer res.Body.Close()

		err = app.readServerJSON(w, res, &comment)

		if err != nil {
			return nil, err
		}	
		comments = append(comments, comment.Response.Items)
	}

	if err != nil {
		return nil, err
	}

	return comments, nil
}