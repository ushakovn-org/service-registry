package main

import (
  "context"

  "github.com/go-resty/resty/v2"
  "github.com/ushakovn-org/service-registry/internal/app"
  "github.com/ushakovn-org/service-registry/internal/app/handlers/graph"
  "github.com/ushakovn-org/service-registry/internal/config"
  "github.com/ushakovn-org/service-registry/internal/pkg/graphviz"
  "github.com/ushakovn-org/service-registry/internal/pkg/repository/proto_dep/etcd"
  "github.com/ushakovn-org/service-registry/internal/pkg/repository/proto_dep/github"
  "github.com/ushakovn-org/service-registry/internal/pkg/usecases/graph/create"
  "github.com/ushakovn-org/service-registry/internal/pkg/usecases/graph/get"
  boilerapp "github.com/ushakovn/boiler/pkg/app"
  boilerconfig "github.com/ushakovn/boiler/pkg/config"
  boileretcd "github.com/ushakovn/boiler/pkg/etcd"
)

func main() {
  ctx := context.Background()

  configClient := boilerconfig.ContextClient(ctx)
  configAppInfo := configClient.GetAppInfo()

  githubAuthTokenProvider := config.NewProvider(ctx, config.GithubAuthToken)

  etcdClient := boileretcd.NewClient(configAppInfo.Name)
  httpClient := resty.New()

  etcdProtoDepsRepo := etcd.NewRepository(etcd.Config{
    EtcdClient: etcdClient,
  })

  githubProtoDepsRepo := github.NewRepository(github.Config{
    HttpClient:        httpClient,
    AuthTokenProvider: githubAuthTokenProvider,
  })

  graphvizRenderer := graphviz.NewRenderer()

  getGraphUseCase := get.NewUseCase(get.Config{
    DepsRepo:      etcdProtoDepsRepo,
    GraphRenderer: graphvizRenderer,
  })

  createGraphUseCase := create.NewUseCase(create.Config{
    DepsLoader: githubProtoDepsRepo,
    DepsSaver:  etcdProtoDepsRepo,
    GetUseCase: getGraphUseCase,
  })

  graphHandler := graph.NewHandler(graph.Config{
    CreateUseCase: createGraphUseCase,
    GetUseCase:    getGraphUseCase,
  })

  serviceRegistry := app.NewService(app.Config{
    Graph: graphHandler,
  })

  boilerapp.NewApp().Run(serviceRegistry)
}
