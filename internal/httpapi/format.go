package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"townbasketball/internal/domain"
)

type listResponse[T any] struct {
	Items  []T `json:"items"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

func parsePage(r *http.Request) domain.Page {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	return domain.Page{Offset: offset, Limit: limit}.Normalize()
}

func parseBoolQuery(r *http.Request, name string) (*bool, error) {
	value := strings.TrimSpace(r.URL.Query().Get(name))
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("%s must be boolean", name)
	}
	return &parsed, nil
}

func encodePretty(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(value)
}

func methodGuard(next http.Handler, methods ...string) http.Handler {
	allowed := map[string]struct{}{}
	for _, method := range methods {
		allowed[method] = struct{}{}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := allowed[r.Method]; !ok {
			w.Header().Set("Allow", strings.Join(methods, ", "))
			writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
