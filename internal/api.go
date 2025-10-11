package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)


// Claims struct for JWT
type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// ================== Handlers ==================

// SignupHandler
// Correct signature for SignupHandler
func SignupHandler(store *Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
            Role     Role   `json:"role"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid request", http.StatusBadRequest)
            return
        }

        hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
        u := User{
            Email:        req.Email,
            PasswordHash: string(hashed),
            Role:         req.Role,
        }

        if err := store.CreateUser(r.Context(), &u); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(u)
    }
}
// Correct signature for LoginHandler
func LoginHandler(store *Store, jwtSecret []byte) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req loginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid request", http.StatusBadRequest)
            return
        }

        user, err := store.GetUserByEmail(r.Context(), req.Email)
        if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
            http.Error(w, "invalid credentials", http.StatusUnauthorized)
            return
        }

        claims := Claims{
            UserID: user.ID,
            Role:   string(user.Role),
            RegisteredClaims: jwt.RegisteredClaims{
                ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenStr, _ := token.SignedString(jwtSecret)

        json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
    }
}

// Same pattern for CreateOrderHandler, CancelOrderHandler, etc.
// All handlers must return `http.HandlerFunc`, do NOT directly accept `w, r`.

// CreateOrderHandler
func CreateOrderHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value("user").(*Claims)
		var o Order
		json.NewDecoder(r.Body).Decode(&o)
		o.CustomerID = claims.UserID
		o.Status = "created"

		err := store.CreateOrder(r.Context(), &o)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		store.CacheOrder(r.Context(), &o)
		go StartOrderTracker(store, o.ID)

		json.NewEncoder(w).Encode(o)
	}
}

// CancelOrderHandler
func CancelOrderHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value("user").(*Claims)
		idStr := mux.Vars(r)["id"]
		id, _ := strconv.Atoi(idStr)

		order, err := store.GetOrderByID(r.Context(), id)
		if err != nil {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		if claims.Role != "admin" && order.CustomerID != claims.UserID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		err = store.CancelOrder(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "cancelled"})
	}
}

// GetAllOrdersHandler
func GetAllOrdersHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value("user").(*Claims)
		if claims.Role != string(Admin) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		rows, err := store.DB.QueryContext(r.Context(),
			`SELECT id, customer_id, items, status, created_at, updated_at, cancelled_at FROM orders`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var orders []Order
		for rows.Next() {
			var o Order
			if err := rows.Scan(&o.ID, &o.CustomerID, &o.Items, &o.Status, &o.CreatedAt, &o.UpdatedAt, &o.CancelledAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			orders = append(orders, o)
		}

		json.NewEncoder(w).Encode(orders)
	}
}

// TrackOrderHandler
func TrackOrderHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		id, _ := strconv.Atoi(idStr)

		order, err := store.GetCachedOrder(r.Context(), id)
		if err != nil {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(order)
	}
}

// ================== Middleware ==================
func AuthMiddleware(jwtSecret []byte) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}