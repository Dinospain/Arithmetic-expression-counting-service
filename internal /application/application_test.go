package application

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Dinospain/Arithmetic-expression-counting-service/pkg/calculation"
)

func TestConfigFromEnv(t *testing.T) {
	t.Run("Default port", func(t *testing.T) {
		os.Unsetenv("PORT") // Убедимся, что переменная окружения не задана
		config := ConfigFromEnv()
		if config.Addr != "8080" {
			t.Errorf("expected default port 8080, got %s", config.Addr)
		}
	})

	t.Run("Custom port", func(t *testing.T) {
		os.Setenv("PORT", "9090")
		config := ConfigFromEnv()
		if config.Addr != "9090" {
			t.Errorf("expected port 9090, got %s", config.Addr)
		}
	})
}

func TestCalcHandler(t *testing.T) {
	t.Run("Valid expression", func(t *testing.T) {
		reqBody := `{"expression": "2+2"}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
		rec := httptest.NewRecorder()

		CalcHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", res.StatusCode)
		}

		body, _ := io.ReadAll(res.Body)
		if !strings.Contains(string(body), "result: 4") {
			t.Errorf("unexpected response body: %s", body)
		}
	})

	t.Run("Invalid expression", func(t *testing.T) {
		reqBody := `{"expression": "2+"}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
		rec := httptest.NewRecorder()

		CalcHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", res.StatusCode)
		}

		body, _ := io.ReadAll(res.Body)
		if !strings.Contains(string(body), "err: calculation failed") {
			t.Errorf("unexpected response body: %s", body)
		}
	})

	t.Run("Malformed JSON", func(t *testing.T) {
		reqBody := `{"expr": "2+2"}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
		rec := httptest.NewRecorder()

		CalcHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", res.StatusCode)
		}
	})
}

func TestApplicationRun(t *testing.T) {
	// Сохраняем текущий stdin и восстанавливаем его после теста
	origStdin := os.Stdin
	defer func() { os.Stdin = origStdin }()

	// Имитация ввода в консоль
	input := "2+2\nexit\n"
	os.Stdin = io.NopCloser(strings.NewReader(input))

	app := New()
	if err := app.Run(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestApplicationRunServer(t *testing.T) {
	app := New()

	// Запуск сервера в отдельной горутине
	go func() {
		if err := app.RunServer(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	// Отправка запроса к серверу
	reqBody := `{"expression": "2+2"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	rec := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "result: 4") {
		t.Errorf("unexpected response body: %s", body)
	}
}

