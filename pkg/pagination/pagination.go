package pagination

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type PaginatedFeedQuery struct {
	Limit    int    `json:"limit" validate:"gte=1,lte=20"`
	Page     int    `json:"page" validate:"gte=1"`
	Sort     string `json:"sort" validate:"oneof=asc desc"`
	Search   string `json:"search" validate:"max=100"`
	Status   string `json:"status" validate:"oneof=Active Inactive Maintenance ''"`
	Since    string `json:"since"`
	Until    string `json:"until"`
	Location string `json:"location"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	limitStr := qs.Get("limit")
	if limitStr == "" {
		fq.Limit = 10
	} else {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			return fq, fmt.Errorf("invalid 'limit' parameter: %w", err)
		}
		fq.Limit = l
	}

	pageStr := qs.Get("page")
	if pageStr == "" {
		fq.Page = 1
	} else {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			return fq, fmt.Errorf("invalid 'page' parameter: %w", err)
		}
		if p < 1 {
			p = 1
		}
		fq.Page = p
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	location := qs.Get("location")
	if location != "" {
		fq.Location = location
	}

	status := qs.Get("status")
	if status != "" {
		fq.Status = status
	}

	since := qs.Get("since")
	if since != "" {
		fq.Since = parseTime(since)
	}

	until := qs.Get("until")
	if until != "" {
		fq.Until = parseTime(until)
	}

	return fq, nil
}

func parseTime(s string) string {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}

type PaginatedResponse[T any] struct {
	Data     []T    `json:"data"`
	Count    int64  `json:"count"`
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}
