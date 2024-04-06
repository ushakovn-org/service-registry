package graph

import (
  "fmt"
  "net/http"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
  result, err := h.getUseCase(r.Context())
  if err != nil {
    msg := fmt.Sprintf("getUseCase: %v", err)
    http.Error(w, msg, http.StatusInternalServerError)
    return
  }

  if _, err = w.Write(result.RenderedGraph); err != nil {
    msg := fmt.Sprintf("http.ResponseWriter.Write: %v", err)
    http.Error(w, msg, http.StatusInternalServerError)
    return
  }
}
