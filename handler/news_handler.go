package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"news-retrieval/constant"
	errconst "news-retrieval/error"
	"news-retrieval/model"
	"news-retrieval/service"
)

type NewsHandler struct {
	news service.NewsService
}

func NewNewsHandler(n service.NewsService) *NewsHandler {
	return &NewsHandler{news: n}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{
		"error": msg,
	})
}

func (h *NewsHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("query"))
	if q == "" {
		writeError(w, http.StatusBadRequest, errconst.ErrQueryParamRequired.Error())
		return
	}

	pgn, err := parsePaginationOptions(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.news.Search(r.Context(), q, pgn)
	if err != nil {
		writeError(w, http.StatusInternalServerError, errconst.ErrGettingNews.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *NewsHandler) Category(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("query"))
	if query == "" {
		writeError(w, http.StatusBadRequest, errconst.ErrQueryParamRequired.Error())
		return
	}

	pgn, err := parsePaginationOptions(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.news.Category(r.Context(), query, pgn)
	if err != nil {
		writeError(w, http.StatusInternalServerError, errconst.ErrGettingCategoryNews.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *NewsHandler) Source(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("query"))
	if query == "" {
		writeError(w, http.StatusBadRequest, errconst.ErrQueryParamRequired.Error())
		return
	}

	pgn, err := parsePaginationOptions(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.news.Source(r.Context(), query, pgn)
	if err != nil {
		writeError(w, http.StatusInternalServerError, errconst.ErrGettingSourceNews.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *NewsHandler) Score(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		writeError(w, http.StatusBadRequest, errconst.ErrQueryParamRequired.Error())
		return
	}

	min, err := strconv.ParseFloat(query, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid query value")
		return
	}

	pgn, err := parsePaginationOptions(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.news.Score(r.Context(), min, pgn)
	if err != nil {
		writeError(w, http.StatusInternalServerError, errconst.ErrGettingScoreNews.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *NewsHandler) Nearby(w http.ResponseWriter, r *http.Request) {

	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	if latStr == "" || lonStr == "" {
		writeError(w, http.StatusBadRequest, "lat and lon parameters are required")
		return
	}

	opts, err := parseSearchOptions(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	pgn, err := parsePaginationOptions(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	lat := *opts.Lat
	lon := *opts.Lon
	radius := opts.Radius

	res, err := h.news.Nearby(r.Context(), lat, lon, radius, pgn)
	if err != nil {
		writeError(w, http.StatusInternalServerError, errconst.ErrGettingNearbyNews.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)
}

// --------------------
// Semantic (LLM-powered)
// --------------------

func (h *NewsHandler) Semantic(w http.ResponseWriter, r *http.Request) {

	q := strings.TrimSpace(r.URL.Query().Get("query"))
	if q == "" {
		writeError(w, http.StatusBadRequest, errconst.ErrQueryParamRequired.Error())
		return
	}

	opts, err := parseSearchOptions(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	pgn, err := parsePaginationOptions(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.news.GetNews(r.Context(), q, opts, pgn)
	if err != nil {
		writeError(w, http.StatusInternalServerError, errconst.ErrGettingNews.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func parseSearchOptions(r *http.Request) (model.SearchOptions, error) {
	q := r.URL.Query()
	opts := model.SearchOptions{}

	// ---------- Location ----------
	if latStr := q.Get("lat"); latStr != "" {
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil || lat < -90 || lat > 90 {
			return opts, errconst.ErrInvalidLatitude
		}
		opts.Lat = &lat
	}

	if lonStr := q.Get("lon"); lonStr != "" {
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil || lon < -180 || lon > 180 {
			return opts, errconst.ErrInvalidLongitude
		}
		opts.Lon = &lon
	}

	// Radius (only meaningful for nearby)
	if rStr := q.Get("radius"); rStr != "" {
		radius, err := strconv.ParseFloat(rStr, 64)
		if err != nil || radius <= 0 {
			return opts, errconst.ErrInvalidRadius
		}
		if radius > constant.MaxRadius {
			return opts, errconst.ErrRadiusExceed
		}
		opts.Radius = radius
	} else {
		opts.Radius = constant.DefaultRadius
	}
	return opts, nil
}

func parsePaginationOptions(r *http.Request) (model.Pagination, error) {
	q := r.URL.Query()

	pg := model.Pagination{
		Page:  1,
		Limit: 5,
	}

	if pStr := q.Get("page"); pStr != "" {
		p, err := strconv.Atoi(pStr)
		if err != nil || p <= 0 {
			return pg, errconst.ErrInvalidPage
		}
		pg.Page = p
	}

	if lStr := q.Get("limit"); lStr != "" {
		l, err := strconv.Atoi(lStr)
		if err != nil || l <= 0 {
			return pg, errconst.ErrInvalidLimit
		}
		if l > 50 {
			return pg, errconst.ErrLimitExceed
		}
		pg.Limit = l
	}
	return pg, nil
}
