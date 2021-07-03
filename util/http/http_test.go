package http

import "testing"

func TestConvertToQueryParams(t *testing.T) {
    url := convertToQueryParams(map[string]interface{}{
        "a": 1,
        "b":"bb?bb",
    })
    t.Log(url)
}
