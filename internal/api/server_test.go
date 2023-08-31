package api

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServerStart(t *testing.T) {
	server := NewServer("8081", &Handler{})
	t.Cleanup(func() {
		err := server.Shutdown(context.Background())
		require.NoError(t, err)
	})

	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			require.Error(t, err)
		}
	}()

	time.Sleep(time.Millisecond * 100)

	_, err := http.Get("http://localhost:8081")
	require.NoError(t, err)
}
