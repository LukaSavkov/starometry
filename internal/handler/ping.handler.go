package handler

import (
	"net/http"

	"github.com/c12s/metrics/internal/utils"
)

type PingHandler struct{}

func (p PingHandler) Ping(rw http.ResponseWriter, h *http.Request) {
	utils.WriteResp("OK", 200, rw)
}
