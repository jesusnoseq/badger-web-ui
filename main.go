package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/gorilla/mux"
)

type App struct {
	db        *badger.DB
	templates *template.Template
}

type KeyValue struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

type Stats struct {
	NumKeys      int64 `json:"num_keys"`
	DatabaseSize int64 `json:"database_size"`
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func main() {
	dbPath := getEnv("BADGER_DB_PATH", "./badger-data")
	opts := badger.DefaultOptions(dbPath)
	if getEnv("BADGER_LOG", "false") != "true" {
		opts.Logger = nil // Disable logging for cleaner output
	}

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Parse templates
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Failed to parse templates:", err)
	}

	app := &App{
		db:        db,
		templates: templates,
	}

	// Setup routes
	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Main page
	r.HandleFunc("/", app.indexHandler).Methods("GET")

	// API routes
	r.HandleFunc("/api/keys", app.listKeysHandler).Methods("GET")
	r.HandleFunc("/api/keys", app.createKeyHandler).Methods("POST")
	r.HandleFunc("/api/keys/{key}", app.getKeyHandler).Methods("GET")
	r.HandleFunc("/api/keys/{key}", app.updateKeyHandler).Methods("PUT")
	r.HandleFunc("/api/keys/{key}", app.deleteKeyHandler).Methods("DELETE")
	r.HandleFunc("/api/stats", app.statsHandler).Methods("GET")
	r.HandleFunc("/api/search", app.searchKeysHandler).Methods("GET")

	port := getEnv("PORT", "8080")
	fmt.Printf("Server starting on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func (app *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	err := app.templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *App) listKeysHandler(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	keys := make([]KeyValue, 0)
	err := app.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		count := 0
		for it.Rewind(); it.Valid() && count < limit; it.Next() {
			item := it.Item()
			key := string(item.Key())

			err := item.Value(func(val []byte) error {
				keys = append(keys, KeyValue{
					Key:       key,
					Value:     string(val),
					CreatedAt: time.Unix(int64(item.Version()), 0),
				})
				return nil
			})
			if err != nil {
				return err
			}
			count++
		}
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(keys); err != nil {
		http.Error(w, "Failed to encode keys", http.StatusInternalServerError)
		return
	}
}

func (app *App) createKeyHandler(w http.ResponseWriter, r *http.Request) {
	var kv KeyValue
	if err := json.NewDecoder(r.Body).Decode(&kv); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if kv.Key == "" {
		http.Error(w, "Key cannot be empty", http.StatusBadRequest)
		return
	}

	err := app.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(kv.Key), []byte(kv.Value))
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	kv.CreatedAt = time.Now()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(kv); err != nil {
		http.Error(w, "Failed to encode kv", http.StatusInternalServerError)
		return
	}
}

func (app *App) getKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	var kv KeyValue
	err := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			kv = KeyValue{
				Key:       key,
				Value:     string(val),
				CreatedAt: time.Unix(int64(item.Version()), 0),
			}
			return nil
		})
	})

	if err == badger.ErrKeyNotFound {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(kv); err != nil {
		http.Error(w, "Failed to encode kv", http.StatusInternalServerError)
		return
	}
}

func (app *App) updateKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	var kv KeyValue
	if err := json.NewDecoder(r.Body).Decode(&kv); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := app.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(kv.Value))
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	kv.Key = key
	kv.CreatedAt = time.Now()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(kv); err != nil {
		http.Error(w, "Failed to encode kv", http.StatusInternalServerError)
		return
	}
}

func (app *App) deleteKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := app.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})

	if err == badger.ErrKeyNotFound {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *App) statsHandler(w http.ResponseWriter, r *http.Request) {
	var stats Stats

	err := app.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		count := int64(0)
		for it.Rewind(); it.Valid(); it.Next() {
			count++
		}
		stats.NumKeys = count
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get database size
	if info, err := os.Stat("./badger-data"); err == nil {
		stats.DatabaseSize = info.Size()
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Failed to encode stats", http.StatusInternalServerError)
		return
	}
}

func (app *App) searchKeysHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	keys := make([]KeyValue, 0)
	err := app.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())

			if strings.Contains(strings.ToLower(key), strings.ToLower(query)) {
				err := item.Value(func(val []byte) error {
					keys = append(keys, KeyValue{
						Key:       key,
						Value:     string(val),
						CreatedAt: time.Unix(int64(item.Version()), 0),
					})
					return nil
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(keys); err != nil {
		http.Error(w, "Failed to encode keys", http.StatusInternalServerError)
		return
	}
}
