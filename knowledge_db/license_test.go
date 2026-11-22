package knowledge

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetDetailsFinishesWhenMoreThanTheWorkerLimitFail(t *testing.T) {
	failing := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "down", http.StatusServiceUnavailable)
	}))
	defer failing.Close()

	licenses := make([]License, 150)
	for i := range licenses {
		licenses[i].DetailsURL = failing.URL
	}

	done := make(chan []License, 1)
	go func() { done <- GetDetails(licenses) }()

	select {
	case got := <-done:
		if len(got) != len(licenses) {
			t.Fatalf("got %d licenses, want %d", len(got), len(licenses))
		}
	case <-time.After(20 * time.Second):
		t.Fatal("GetDetails did not return when every details fetch failed")
	}
}
