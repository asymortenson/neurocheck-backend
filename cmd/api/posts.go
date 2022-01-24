package main

import (
	"fmt"
	"net/http"
	"strings"

	"vkparser.com/internal/data"
)





func (app *application) getPosts(w http.ResponseWriter, r *http.Request) {
	var requestData data.RequestData
	var fields []string


	err := app.readJSON(w, r, &requestData)

	for _, field := range requestData.Fields {
		fields = append(fields,field.Name + "=" + field.Value)
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	url := fmt.Sprintf("https://api.vk.com/method/wall.get?%s&access_token=%s&v=5.131", strings.Join(fields, "&"), requestData.AccessToken)

    comments, err := app.fetchPosts(w, r, url, requestData, 1)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, comments, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	}
