// Copyright 2022 Tatris Project Authors. Licensed under Apache-2.0.

package handler

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestIngestHandler(t *testing.T) {

	t.Run("ingest", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := &http.Request{
			URL:    &url.URL{},
			Header: make(http.Header),
		}
		c.Request = req
		p := gin.Params{}
		p = append(p, gin.Param{Key: "index", Value: "storage_product"})
		c.Params = p
		c.Request.Header.Set("Content-Type", "application/json;charset=utf-8")
		c.Request.Body = io.NopCloser(bytes.NewBufferString("{\"documents\":[{\"name\":\"tatris\",\"desc\":\"Time-aware storage and search system\"},{\"name\":\"mysql\",\"desc\":\"Relational database\"}]}"))
		IngestHandler(c)
		fmt.Println(w)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}