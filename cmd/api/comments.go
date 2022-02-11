package main

import (
	"errors"
	"net/http"
	"sync"

	"github.com/SevereCloud/vksdk/v2/api"
	"vkparser.com/internal/data"
	"vkparser.com/internal/validator"
)


func (app *application) getComments(r *http.Request, comments *data.Comments, values api.Params) (int, error) {
	vk := app.contextGetVK(r)

	wait := make(chan struct{})

	posts, err := vk.WallGet(values)

	if errors.Is(err, api.ErrAuth) {
		return http.StatusUnauthorized, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	var wg sync.WaitGroup
	wg.Add(len(posts.Items))

	for _, post := range posts.Items {
		var appendComment data.Comment

		commentsOfPost, err := vk.WallGetComments(api.Params{"owner_id": post.OwnerID, "post_id": post.ID})

		if err != nil {
			return http.StatusBadRequest, err
		}
		for _, comment := range commentsOfPost.Items {
			if comment.Text != "" {
				appendComment.ID = comment.ID
				appendComment.GroupID = comment.OwnerID
				appendComment.Text = comment.Text

				comments.Items = append(comments.Items, appendComment)
			}
		}
		go timeoutPerRequest(wait)
		<-wait
	}
	return http.StatusOK, nil
}

func (app *application) getRatedComments(w http.ResponseWriter, r *http.Request) {

	var comments data.Comments

	query := r.URL.Query()

	validate, ok := validator.ValidateQuery(query, []string{"count", "type"})

	if !ok {
		app.queryMissingError(w, r, validate)
		return
	}

	status, err := app.getComments(r, &comments, app.ConvertQueryToInterface(query))

	if err != nil {
		app.errorResponse(w, r, status, err.Error())
		return
	}

	err = app.CheckToxicity(envelope{"count": len(comments.Items), "items": comments.Items}, &comments, "https://zubibubi.pythonanywhere.com/"+query.Get("type"))

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ratedComments := &data.Comments{
		Count:  comments.Count,
		Items:  comments.Items,		
	}

	app.writeJSON(w, http.StatusOK, envelope{"response": ratedComments}, nil)
}
