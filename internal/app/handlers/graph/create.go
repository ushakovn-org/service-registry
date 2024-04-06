package graph

import (
  "fmt"
  "io"
  "net/http"
  "regexp"
  "strings"

  validation "github.com/go-ozzo/ozzo-validation"
  "github.com/ushakovn-org/service-registry/internal/pkg/unmarshal"
  "github.com/ushakovn-org/service-registry/internal/pkg/usecases/graph/create"
)

var (
  regexRepo = regexp.MustCompile(`^[a-zA-Z0-9_-]+\/[a-zA-Z0-9_-]+$`)
  regexYaml = regexp.MustCompile(`^.+\.yaml$`)
)

type CreateRequest struct {
  GithubRepoName string `json:"github_repo_name"`
  ProtoDepsPath  string `json:"proto_deps_path"`
}

func (r *CreateRequest) Validate() error {
  return validation.ValidateStruct(r,
    validation.Field(&r.GithubRepoName, validation.Required, validation.Match(regexRepo)),
    validation.Field(&r.ProtoDepsPath, validation.Required, validation.Match(regexYaml)),
  )
}

func (r *CreateRequest) toParams() create.Params {
  parts := strings.SplitN(r.GithubRepoName, "/", 2)
  var (
    owner = parts[0]
    repo  = parts[1]
  )
  return create.Params{
    Owner: owner,
    Repo:  repo,
    Path:  r.ProtoDepsPath,
  }
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
  buf, err := io.ReadAll(r.Body)
  if err != nil {
    msg := fmt.Sprintf("http.Request.Body: io.ReadAll: %v", err)
    http.Error(w, msg, http.StatusInternalServerError)
    return
  }

  req, err := unmarshal.JSON[CreateRequest](buf)
  if err != nil {
    msg := fmt.Sprintf("invalid request: unmarshal.JSON: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }

  if err = req.Validate(); err != nil {
    msg := fmt.Sprintf("invalid request: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }

  result, err := h.createUseCase(r.Context(), req.toParams())
  if err != nil {
    msg := fmt.Sprintf("createUseCase: %v", err)
    http.Error(w, msg, http.StatusInternalServerError)
    return
  }

  if _, err = w.Write(result.RenderedGraph); err != nil {
    msg := fmt.Sprintf("http.ResponseWriter.Write: %v", err)
    http.Error(w, msg, http.StatusInternalServerError)
    return
  }
}
