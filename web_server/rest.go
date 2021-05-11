package web_server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// populateRESTRoutes populates the given router with the routes needed for REST.
func (s *WebServer) populateRESTRoutes(r *mux.Router) {
	r.HandleFunc("/seasons/{id:[0-9]+}", s.getSeasonByIdHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/seasons/{seasonId:[0-9]+}/episodes", s.getEpisodesOfSeason).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/podcasts/by-key/{key}", s.getPodcastByKeyHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/podcasts/{id}", s.getPodcastByIdHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/podcasts", s.getPodcastsHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/podcasts/{podcastId:[0-9]+}/seasons/last", s.getLastSeasonOfPodcastHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/podcasts/{podcastId:[0-9]+}/seasons/{seasonNum:[0-9]+}", s.getLastSeasonOfPodcastHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/podcasts/{podcastId:[0-9]+}/seasons", s.getSeasonsOfPodcastHandler).Methods(http.MethodGet, http.MethodOptions)
}

// getSeasonByIdHandler retrieves a season by id.
func (s *WebServer) getSeasonByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		writeString(w, http.StatusBadRequest, "invalid season id")
		return
	}
	season, err := s.stores.Seasons.ById(id)
	if err != nil {
		writeString(w, http.StatusInternalServerError, "could not load seasons")
		log.Printf("error while retrieving seasons: %v", err)
		return
	}
	writeJSON(w, season)
}

// getLastSeasonOfPodcastHandler retrieves the last season of a given podcast.
func (s *WebServer) getLastSeasonOfPodcastHandler(w http.ResponseWriter, r *http.Request) {
	podcastId, err := strconv.Atoi(mux.Vars(r)["podcastId"])
	if err != nil {
		writeString(w, http.StatusBadRequest, "invalid podcast id")
		return
	}
	seasons, err := s.stores.Seasons.ByPodcast(podcastId)
	if err != nil {
		writeString(w, http.StatusInternalServerError, "could not retrieve seasons by podcast")
		log.Printf("error while retrieving seasons by podcast %d: %v", podcastId, err)
		return
	}
	if len(seasons) == 0 {
		writeString(w, http.StatusNotFound, "podcast has no seasons")
		return
	}
	// Find the one with highest num.
	season := seasons[0]
	for _, s := range seasons {
		if s.Num > season.Num {
			season = s
		}
	}
	writeJSON(w, season)
}

// getPodcastByIdHandler retrieves a podcast by the given podcast id.
func (s *WebServer) getPodcastByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		writeString(w, http.StatusBadRequest, "id was empty")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeString(w, http.StatusBadRequest, "id was not numeric")
		return
	}
	podcast, err := s.stores.Podcasts.ById(id)
	if err != nil {
		writeString(w, http.StatusNotFound, "could not retrieve podcast")
		return
	}
	writeJSON(w, podcast)
}

// getPodcastByKeyHandler retrieves a podcast by a given key.
func (s *WebServer) getPodcastByKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	if key == "" {
		writeString(w, http.StatusBadRequest, "key was empty")
		return
	}
	podcast, err := s.stores.Podcasts.ByKey(key)
	if err != nil {
		writeString(w, http.StatusNotFound, "could not retrieve podcast")
		return
	}
	writeJSON(w, podcast)
}

// getPodcastsHandler retrieves all podcasts.
func (s *WebServer) getPodcastsHandler(w http.ResponseWriter, _ *http.Request) {
	podcasts, err := s.stores.Podcasts.All()
	if err != nil {
		writeString(w, http.StatusInternalServerError, "could not retrieve podcasts")
		log.Printf("error while retrieving podcasts: %v", err)
		return
	}
	writeJSON(w, podcasts)
}

// getSeasonOfPodcastHandler retrieves a season in a podcast by it's number.
func (s *WebServer) getSeasonOfPodcastHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	podcastId, err := strconv.Atoi(vars["podcastId"])
	if err != nil {
		writeString(w, http.StatusBadRequest, "invalid podcast id")
		return
	}
	seasonNum, err := strconv.Atoi(vars["seasonNum"])
	if err != nil {
		writeString(w, http.StatusBadRequest, "invalid season num")
		return
	}
	seasons, err := s.stores.Seasons.ByPodcast(podcastId)
	if err != nil {
		writeString(w, http.StatusInternalServerError, "could not retrieve seasons")
		log.Printf("error while retrieving seasons: %v", err)
		return
	}
	if len(seasons) == 0 {
		writeString(w, http.StatusNotFound, "not seasons found for given podcast")
		return
	}
	// Find the season with the given number.
	for _, season := range seasons {
		if season.Num == seasonNum {
			writeJSON(w, season)
			return
		}
	}
	// No season with given number found.
	writeString(w, http.StatusNotFound, "no season with given number found")
}

// getEpisodesOfSeason retrieves the episodes for the given season.
func (s *WebServer) getEpisodesOfSeason(w http.ResponseWriter, r *http.Request) {
	seasonId, err := strconv.Atoi(mux.Vars(r)["seasonId"])
	if err != nil {
		writeString(w, http.StatusBadRequest, "invalid season id")
		return
	}
	episodes, err := s.stores.Episodes.BySeason(seasonId)
	if err != nil {
		writeString(w, http.StatusInternalServerError, "could not retrieve episodes for given season")
		log.Printf("error while retrieving episodes by season %d: %v", seasonId, err)
		return
	}
	writeJSON(w, episodes)
}

// getSeasonsOfPodcastHandler retrieves the seasons for the given podcast.
func (s *WebServer) getSeasonsOfPodcastHandler(w http.ResponseWriter, r *http.Request) {
	podcastId, err := strconv.Atoi(mux.Vars(r)["podcastId"])
	if err != nil {
		writeString(w, http.StatusBadRequest, "invalid podcast id")
		return
	}
	seasons, err := s.stores.Seasons.ByPodcast(podcastId)
	if err != nil {
		writeString(w, http.StatusInternalServerError, "could not retrieve seasons for given podcast")
		log.Printf("error while retrieving seasons by podcast %d: %v", podcastId, err)
		return
	}
	writeJSON(w, seasons)
}

// writeJSON writes the given interface marshalled and with status code http.StatusOK.
func writeJSON(w http.ResponseWriter, response interface{}) {
	s, err := json.Marshal(response)
	if err != nil {
		log.Printf("could not marshal %v: %v", response, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	write(w, http.StatusOK, s)
}

// writeString writes the given string and converts it to bytes.
func writeString(w http.ResponseWriter, statusCode int, response string) {
	write(w, statusCode, []byte(response))
}

// write writes the given response and status code and logs a possible write error.
func write(w http.ResponseWriter, statusCode int, response []byte) {
	w.WriteHeader(statusCode)
	_, err := w.Write(response)
	if err != nil {
		log.Printf("could not write response: %v", err)
	}
}
