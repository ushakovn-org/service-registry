package github

import (
  "context"
  "encoding/base64"
  "fmt"
  "net/http"
  "strings"

  "github.com/ushakovn-org/service-registry/internal/pkg/unmarshal"
)

type githubResponse struct {
  Content  string `json:"content"`
  Encoding string `json:"encoding"`
}

type responseProtoImport string

type responseProtoDep struct {
  Import responseProtoImport `yaml:"import"`
}

type responseProtoDeps struct {
  ExternalDeps []responseProtoDep `yaml:"external_deps"`
}

type parsedProtoImport struct {
  Owner   string
  Repo    string
  Path    string
  Package string
  Commit  string
}

func (r *githubResponse) Validate() error {
  if r.Content == "" {
    return fmt.Errorf("has blank content")
  }
  if r.Encoding != "base64" {
    return fmt.Errorf("has unsupported encoding: %s", r.Encoding)
  }
  return nil
}

func (r *githubResponse) Decode() ([]byte, error) {
  if err := r.Validate(); err != nil {
    return nil, fmt.Errorf("invalid response: %w", err)
  }
  decoded, err := base64.StdEncoding.DecodeString(r.Content)
  if err != nil {
    return nil, fmt.Errorf("base64.StdEncoding.DecodeString: %w", err)
  }
  return decoded, nil
}

func (p responseProtoImport) String() string {
  return string(p)
}

func (p responseProtoImport) Parse() (*parsedProtoImport, error) {
  path := string(p)

  // github.com/<owner>/<repo>/<path>.proto@<commit>
  partsImport := strings.SplitN(path, "/", 4)

  if len(partsImport) != 4 {
    return nil, fmt.Errorf("invalid path: %s. expected pattern: github.com/<owner>/<repo>/<path>.proto@<commit>", path)
  }
  ownerPart := partsImport[1]
  repoPart := partsImport[2]
  pathPart := partsImport[3]

  // <path>.proto@<commit>
  partsImport = strings.SplitN(pathPart, "@", 2)

  if len(partsImport) != 2 {
    return nil, fmt.Errorf("invalid path: %s. expected pattern: github.com/<owner>/<repo>/<path>.proto@<commit>", path)
  }
  pathPart = partsImport[0]
  commitPart := partsImport[1]

  packageParts := strings.Split(pathPart, "/")

  if len(partsImport) == 0 {
    return nil, fmt.Errorf("invalid path: %s. expected pattern: github.com/<owner>/<repo>/<path>.proto@<commit>", path)
  }
  var packagePart string

  switch {
  // <path>.proto = <package>/<file>.proto
  case len(partsImport) >= 2:
    packagePart = packageParts[len(packageParts)-2]

  // <path>.proto = <file>.proto
  case len(partsImport) >= 1:
    packagePart = packageParts[len(packageParts)-1]
    packagePart = strings.TrimSuffix(packagePart, ".proto")
  }

  return &parsedProtoImport{
    Owner:   ownerPart,
    Repo:    repoPart,
    Path:    pathPart,
    Package: packagePart,
    Commit:  commitPart,
  }, nil
}

func (r *Repository) doRequest(ctx context.Context, url string) (*githubResponse, error) {
  authToken := r.authToken.Provide().String()

  resp, err := r.httpClient.R().
    SetAuthToken(authToken).
    SetContext(ctx).
    Get(url)

  if err != nil {
    return nil, fmt.Errorf("client.Get: %s: %w", url, err)
  }

  if resp.StatusCode() != http.StatusOK || resp.Body() == nil {
    return nil, fmt.Errorf("client.Get: invalid response: resp.Status: %s", resp.Status())
  }

  out, err := unmarshal.JSON[githubResponse](resp.Body())
  if err != nil {
    return nil, fmt.Errorf("unmarshal.JSON: %w", err)
  }
  return out, nil
}
