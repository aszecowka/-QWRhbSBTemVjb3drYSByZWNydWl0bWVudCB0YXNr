// +build integration

package integration_test

import (
	"encoding/json"
	"fmt"
	"github.com/aszecowka/QWRhbSBTemVjb3drYSByZWNydWl0bWVudCB0YXNr/internal/weather"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestIntegration(t *testing.T) {
	url := os.Getenv("BULK_API_URL")
	require.NotEmpty(t, url)

	fullURL := fmt.Sprintf("%s/weather?city=Madrid&city=Barcelona", url)
	t.Log("Performing GET request: ", fullURL)
	resp, err := http.Get(fullURL)
	require.NoError(t, err)
	defer func() {
		_, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)
	}()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	out := map[string]*weather.OpenWeatherResponse{}

	err = json.NewDecoder(resp.Body).Decode(&out)
	require.NoError(t, err)
	require.Len(t, out, 2)
	require.NotNil(t, out["Madrid"])
	require.NotNil(t, out["Barcelona"])

}
