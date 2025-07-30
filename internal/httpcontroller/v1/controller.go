package v1

import (
	"archiveFiles/internal/httpcontroller/v1/api"
	"archiveFiles/internal/usecase"
	"embed"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=api/types.cfg.yaml api/openapi.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=api/server.cfg.yaml api/openapi.yaml

//go:embed api/openapi.yaml
var OpenApi embed.FS

var _ api.ServerInterface = (*V1)(nil)

type V1 struct {
	usecase usecase.ILinks
}
