// Package server exposes HTTP handlers and wiring for the subscriptions
// service. It registers routes for creating, listing, retrieving, updating,
// deleting subscriptions and calculating aggregated values (total cost).
//
// The package also contains swagger annotations used by swag to generate API
// documentation.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"task_effective_mobile/internal/config"
	"task_effective_mobile/internal/repositories"
	"task_effective_mobile/pkg/logger"
)

// @title Subscriptions API
// @version 1.0
// @description API for managing user subscriptions.
// @host localhost:8080
// @BasePath /

// @Summary Create subscription
// @Description Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body object true "Subscription to create"
// @Success 201 {object} map[string]int
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions [post]
func createSubscriptionsDoc() {}

// @Summary List subscriptions
// @Description Get list of subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} object
// @Failure 500 {string} string
// @Router /subscriptions [get]
func listSubscriptionsDoc() {}

// subscriptionsHandler returns an http.HandlerFunc that handles requests to
// the /subscriptions endpoint. It supports POST for creating a subscription
// and GET for listing all subscriptions.
func subscriptionsHandler(ctx context.Context, repo *repositories.SubscriptionsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger(ctx)
		log.Info("subscriptionsHandler:", "Received", r.Method, "request at", r.URL.Path)
		switch r.Method {
		case http.MethodPost:
			defer func() { _ = r.Body.Close() }()
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "failed to read request body", http.StatusBadRequest)
				log.Error("subscriptionsHandler: failed to read request body", "error", err)
				return
			}
			var req struct {
				ServiceName string `json:"service_name"`
				Price       int    `json:"price"`
				UserID      string `json:"user_id"`
				StartDate   string `json:"start_date"`
				EndDate     string `json:"end_date"`
			}
			if err := json.Unmarshal(body, &req); err != nil {
				http.Error(w, "invalid json body", http.StatusBadRequest)
				log.Error("subscriptionsHandler: failed to unmarshal request body", "error", err)
				return
			}

			if req.ServiceName == "" || req.UserID == "" || req.StartDate == "" {
				http.Error(w, "missing required fields", http.StatusBadRequest)
				log.Error("subscriptionsHandler: Missing required fields", "error", err)
				return
			}
			if req.Price < 0 {
				http.Error(w, "price must be non-negative", http.StatusBadRequest)
				log.Error("subscriptionsHandler: Price must be non-negative", req.Price)
				return
			}

			id, err := repo.CreateSub(r.Context(), req.ServiceName, req.Price, req.UserID, req.StartDate, req.EndDate)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to create subscription: %v", err), http.StatusInternalServerError)
				log.Error("subscriptionsHandler: failed to create subscription", "error", err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]int{"id": id})
			log.Info("subscriptionsHandler: Created subscription", "id", id)

		case http.MethodGet:
			subs, err := repo.GetSubsList(r.Context())
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to get subscriptions: %v", err), http.StatusInternalServerError)
				log.Error("subscriptionsHandler: failed to get subscriptions", "error", err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(subs); err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
				log.Error("subscriptionsHandler: failed to encode subscriptions response", "error", err)
				return
			}
			log.Info("subscriptionsHandler: Returned subscriptions list", "count", len(subs))

		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			log.Error("subscriptionsHandler: Unsupported method", r.Method)
		}
	}
}

// @Summary Get subscription by id
// @Description Get subscription by id
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} object
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [get]
func getSubscriptionsDoc() {}

// @Summary Update subscription by id
// @Description Update subscription fields partially
// @Tags subscriptions
// @Accept json
// @Param id path int true "Subscription ID"
// @Param subscription body object true "Fields to update"
// @Success 204 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [put]
func updateSubscriptionsDoc() {}

// @Summary Delete subscription by id
// @Description Delete subscription
// @Tags subscriptions
// @Param id path int true "Subscription ID"
// @Success 204 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [delete]
func deleteSubscriptionsDoc() {}

// @Summary Get total cost
// @Description Calculate total sum of subscription prices for the given filters and period (optional filters: user_id, service_name, start_date, end_date in MM-YYYY)
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service name"
// @Param start_date query string false "Period start in MM-YYYY"
// @Param end_date query string false "Period end in MM-YYYY"
// @Success 200 {object} map[string]int
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/total [get]
func subscriptionsTotalDoc() {}

func subscriptionsTotalHandler(ctx context.Context, repo *repositories.SubscriptionsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger(ctx)
		log.Info("subscriptionsTotalHandler: Received", r.Method, "request at", r.URL.Path)
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			log.Error("subscriptionsTotalHandler: Unsupported method", r.Method)
			return
		}

		q := r.URL.Query()
		var userIDPtr *string
		if v := q.Get("user_id"); v != "" {
			userIDPtr = &v
		}
		var serviceNamePtr *string
		if v := q.Get("service_name"); v != "" {
			serviceNamePtr = &v
		}
		var startPtr *string
		if v := q.Get("start_date"); v != "" {
			startPtr = &v
		}
		var endPtr *string
		if v := q.Get("end_date"); v != "" {
			endPtr = &v
		}

		total, err := repo.GetTotalCost(r.Context(), userIDPtr, serviceNamePtr, startPtr, endPtr)
		if err != nil {
			if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "cannot be empty") {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Error("subscriptionsTotalHandler: bad request", "error", err)
				return
			}
			http.Error(w, fmt.Sprintf("failed to calculate total: %v", err), http.StatusInternalServerError)
			log.Error("subscriptionsTotalHandler: failed to calculate total", "error", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]int{"total": total}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			log.Error("subscriptionsTotalHandler: failed to encode response", "error", err)
			return
		}
		log.Info("subscriptionsTotalHandler: Returned total", "total", total)
	}
}

// subscriptionsIDHandler returns an http.HandlerFunc that handles GET, PUT
// and DELETE for the /subscriptions/{id} endpoint. It supports retrieving
// a single subscription, performing partial updates, and deleting the
// subscription by id.
func subscriptionsIDHandler(ctx context.Context, repo *repositories.SubscriptionsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger(ctx)
		log.Info("subscriptionsIDHandler: Received ", r.Method, " request at ", r.URL.Path)
		idPart := strings.TrimPrefix(r.URL.Path, "/subscriptions/")
		idPart = strings.Trim(idPart, "/")
		if idPart == "" {
			http.Error(w, "missing id in path", http.StatusBadRequest)
			log.Error("subscriptionsIDHandler: Missing id in path", "path", r.URL.Path)
			return
		}
		id, err := strconv.Atoi(idPart)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			log.Error("subscriptionsIDHandler: Invalid id in path", "path", r.URL.Path)
			return
		}

		switch r.Method {
		case http.MethodGet:
			sub, err := repo.GetSub(r.Context(), id)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					http.Error(w, "not found", http.StatusNotFound)
					log.Error("subscriptionsIDHandler: Not Found", "id", id)
					return
				}
				http.Error(w, fmt.Sprintf("failed to get subscription: %v", err), http.StatusInternalServerError)
				log.Error("subscriptionsIDHandler: Failed to get subscription", "id", id)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(sub)
			log.Info("subscriptionsIDHandler: Get subscription", "id", id)

		case http.MethodPut:
			defer func() { _ = r.Body.Close() }()
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "failed to read request body", http.StatusBadRequest)
				log.Error("subscriptionsIDHandler: Failed to read request body", "error", err)
				return
			}
			var req struct {
				ServiceName *string `json:"service_name"`
				Price       *int    `json:"price"`
				UserID      *string `json:"user_id"`
				StartDate   *string `json:"start_date"`
				EndDate     *string `json:"end_date"`
			}
			if err := json.Unmarshal(body, &req); err != nil {
				http.Error(w, "invalid json body", http.StatusBadRequest)
				log.Error("subscriptionsIDHandler: Failed to unmarshal request body", "error", err)
				return
			}

			if req.Price != nil && *req.Price < 0 {
				http.Error(w, "price must be non-negative", http.StatusBadRequest)
				log.Error("subscriptionsIDHandler: Price must be non-negative", req.Price)
				return
			}

			if err := repo.UpdateSub(r.Context(), id, req.ServiceName, req.Price, req.UserID, req.StartDate, req.EndDate); err != nil {
				if strings.Contains(err.Error(), "not found") {
					http.Error(w, "not found", http.StatusNotFound)
					log.Error("subscriptionsIDHandler: Not Found", "id", id)
					return
				}
				if strings.Contains(err.Error(), "no fields to update") || strings.Contains(err.Error(), "invalid") {
					http.Error(w, err.Error(), http.StatusBadRequest)
					log.Error("subscriptionsIDHandler: Failed to update subscription", "id", id)
					return
				}
				http.Error(w, fmt.Sprintf("failed to update subscription: %v", err), http.StatusInternalServerError)
				log.Error("subscriptionsIDHandler: Failed to update subscription", "id", id)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			log.Info("subscriptionsIDHandler: Updated subscription", "id", id)

		case http.MethodDelete:
			if err := repo.DeleteSub(r.Context(), id); err != nil {
				if strings.Contains(err.Error(), "not found") {
					http.Error(w, "not found", http.StatusNotFound)
					log.Error("subscriptionsIDHandler: Not Found", "id", id)
					return
				}
				http.Error(w, fmt.Sprintf("failed to delete subscription: %v", err), http.StatusInternalServerError)
				log.Error("subscriptionsIDHandler: Failed to delete subscription", "id", id)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			log.Info("subscriptionsIDHandler: Deleted subscription", "id", id)

		default:
			w.Header().Set("Allow", "GET, PUT, DELETE")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			log.Error("subscriptionsIDHandler: Unsupported method", r.Method)
		}
	}
}

// Start initializes the server routing and starts the HTTP server.
//
// It reads configuration using the internal config package, creates a
// SubscriptionsRepository and registers handlers on a new ServeMux. The
// function blocks until the HTTP server exits or returns an error.
func Start(ctx context.Context) error {
	log := logger.GetLogger(ctx)
	mux := http.NewServeMux()
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("start: failed to load config: %w", err)
	}
	repo, err := repositories.NewSubscriptionsRepository(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("start: failed to create subscriptions repository: %w", err)
	}

	mux.HandleFunc("/subscriptions", subscriptionsHandler(ctx, repo))
	mux.HandleFunc("/subscriptions/", subscriptionsIDHandler(ctx, repo))
	mux.HandleFunc("/subscriptions/total", subscriptionsTotalHandler(ctx, repo))

	log.Info("Starting server on port", "port", cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), mux)
	if err != nil {
		return fmt.Errorf("start: error while starting http server: %w", err)
	}
	return nil
}
