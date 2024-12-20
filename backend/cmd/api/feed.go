package main

import (
	"net/http"

	"github.com/AdmFjalar/CS301.3-Time-Tracker/internal/store"
)

// getUserFeedHandler godoc
//
//	@Summary		Fetches the user feed
//	@Description	Fetches the user feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]store.Timestamp
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize the query parameters with default values
	fq := store.Query{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Search: "",
	}

	// Parse the query parameters from the request
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the parsed query parameters
	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Retrieve the user from the request context
	ctx := r.Context()
	user := getUserFromContext(r)

	// Fetch the user feed from the store
	feed, err := app.store.Timestamps.GetUserFeed(ctx, user.ID, fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Send the user feed as a JSON response
	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
