package github

import (
  "context"
  "fmt"

  "github.com/go-resty/resty/v2"
  "github.com/ushakovn-org/service-registry/internal/pkg/models/proto_dep"
  "github.com/ushakovn-org/service-registry/internal/pkg/unmarshal"
  "github.com/ushakovn-org/service-registry/internal/pkg/usecases/graph/create"
  "github.com/ushakovn/boiler/pkg/config/types"
)

type Repository struct {
  authToken  AuthTokenProvider
  httpClient *resty.Client
}

type AuthTokenProvider interface {
  Provide() (authToken types.Value)
}

type Config struct {
  HttpClient        *resty.Client
  AuthTokenProvider AuthTokenProvider
}

func NewRepository(config Config) *Repository {
  return &Repository{
    httpClient: config.HttpClient,
    authToken:  config.AuthTokenProvider,
  }
}

func (r *Repository) Load(ctx context.Context, params create.LoadParams) (*proto_dep.ProtoDep, error) {
  url := buildGithubUrl(params)

  resp, err := r.doRequest(ctx, url)
  if err != nil {
    return nil, fmt.Errorf("http.Get: %w", err)
  }
  decoded, err := resp.Decode()
  if err != nil {
    return nil, fmt.Errorf("resp.Decode: %w", err)
  }
  respDeps, err := unmarshal.YAML[responseProtoDeps](decoded)
  if err != nil {
    return nil, fmt.Errorf("unmarshal.YAML: %w", err)
  }
  protoDepTargets := make([]proto_dep.ProtoDepTarget, 0, len(respDeps.ExternalDeps))

  for _, respProtoDep := range respDeps.ExternalDeps {
    protoImport := respProtoDep.Import

    var parsedImport *parsedProtoImport

    parsedImport, err = protoImport.Parse()
    if err != nil {
      return nil, fmt.Errorf("respProtoDep.Import.Parse: %w", err)
    }
    protoDepTargets = append(protoDepTargets, proto_dep.ProtoDepTarget{
      Source: parsedImport.Repo,
      Import: protoImport.String(),
      Commit: parsedImport.Commit,
    })
  }

  return &proto_dep.ProtoDep{
    Source:  params.Repo,
    Targets: protoDepTargets,
  }, nil
}

func buildGithubUrl(params create.LoadParams) string {
  return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s",
    params.Owner,
    params.Repo,
    params.Path,
  )
}
