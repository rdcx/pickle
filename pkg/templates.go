package pickle

import _ "embed"

//go:embed templates/docker-compose.yaml.template
var templateDockerComposeYaml string

//go:embed templates/mux/Dockerfile.template
var templateMuxDockerfile string

//go:embed templates/mux/go.mod.template
var templateMuxGoMod string

//go:embed templates/mux/main_test.go.template
var templateMuxMainTest string

//go:embed templates/mux/main.go.template
var templateMuxMain string

//go:embed templates/mux/redis/show.main.go.template
var templateRedisMainShow string

//go:embed templates/mux/redis/store.main.go.template
var templateRedisMainStore string

//go:embed templates/gateway/Dockerfile.template
var templateGatewayDockerfile string

//go:embed templates/gateway/go.mod.template
var templateGatewayGoMod string

//go:embed templates/gateway/main.go.template
var templateGatewayMain string
