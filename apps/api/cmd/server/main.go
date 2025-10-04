package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/go-libsql"

	"github.com/ryu1013-job/next-go-template/apps/api/gen"
	"github.com/ryu1013-job/next-go-template/apps/api/internal/features/todo/controller"
	"github.com/ryu1013-job/next-go-template/apps/api/internal/features/todo/repository"
	"github.com/ryu1013-job/next-go-template/apps/api/internal/features/todo/usecase"
	"github.com/ryu1013-job/next-go-template/apps/api/internal/infra/db"
	"github.com/ryu1013-job/next-go-template/apps/api/internal/infra/migrate"
)

func main() {
	// --- .env ロード（ローカル開発用。存在しなければ何もしない）---
	loadDotenv()

	ctx := context.Background()

	dsn := buildTursoDSN(
		must("TURSO_DATABASE_URL"),
		os.Getenv("TURSO_AUTH_TOKEN"),
	)

	sqlDB, err := sql.Open("libsql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	migDir := filepath.Join("..", "..", "infra", "db", "migrations")
	if err := migrate.Run(ctx, sqlDB, migDir); err != nil {
		log.Fatal(err)
	}

	queries := db.New(sqlDB)
	todoRepo := repository.NewTodoRepository(queries)
	todoUsecase := usecase.NewTodoUsecase(todoRepo)
	todoController := controller.NewTodoController(todoUsecase)

	mux := http.NewServeMux()
	handler := gen.HandlerFromMux(todoController, mux)

	port := env("API_PORT", "8080")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("listen :" + port)
	log.Fatal(srv.ListenAndServe())
}

// --- helpers ---

// ローカル開発時だけ .env を試し読みする
func loadDotenv() {
	// apps/api 直下 or リポジトリルートの .env を順に試す
	candidates := []string{
		".env",
		filepath.Join("..", "..", ".env"),
		filepath.Join("apps", "api", ".env"),
	}
	for _, p := range candidates {
		if err := godotenv.Load(p); err == nil {
			log.Printf("loaded .env: %s\n", p)
			return
		}
	}
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s required", k)
	}
	return v
}

// TURSO_AUTH_TOKEN があれば dsn に authToken= を付与する（既に付いていればそのまま）
func buildTursoDSN(raw, token string) string {
	if token == "" || strings.Contains(raw, "authToken=") {
		return raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return raw // 失敗時はそのまま
	}
	q := u.Query()
	q.Set("authToken", token)
	u.RawQuery = q.Encode()
	return u.String()
}
