package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/paoloposso/manhattan-points/points"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func main() {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	registerRoutes(router)

	errs := make(chan error, 2)

	port := getHttpPort()

	go func() {
		fmt.Println("Listening on port ", port)
		errs <- http.ListenAndServe(port, router)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Println("Terminated ", <-errs)
}

func getHttpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func registerRoutes(router *chi.Mux) {
	baseUrl := "/api"
	router.Get(baseUrl+"/points", getPoints)
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	x, invalidX := strconv.ParseFloat(r.URL.Query().Get("x"), 32)
	y, invalidY := strconv.ParseFloat(r.URL.Query().Get("y"), 32)
	dist, invalidDist := strconv.ParseFloat(r.URL.Query().Get("distance"), 64)

	if invalidX != nil || invalidY != nil || invalidDist != nil {
		code := 500
		msg := fmt.Sprintf("{ \"message\": \"%s\" }", "Invalid parameters")
		w.WriteHeader(code)
		w.Write([]byte(msg))
		return
	}

	p, err := points.GetPointsInsideManhattanDistance(points.Point{X: x, Y: y}, dist)

	if err != nil {
		code := 500
		msg := formatError(err, fmt.Sprintf("%b", code))
		w.WriteHeader(code)
		w.Write([]byte(msg))
		return
	}

	res, err := json.Marshal(p)

	if err != nil {
		code := 500
		msg := formatError(err, fmt.Sprintf("%b", code))
		w.WriteHeader(code)
		w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func formatError(err error, code string) string {
	return fmt.Sprintf("{ \"message\": \"%s\" }", err)
}
