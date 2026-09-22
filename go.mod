module github.com/arnavdugar/hsm

go 1.26.0

require (
	github.com/stretchr/testify v1.12.1
	go.uber.org/mock v0.6.0
	go.yaml.in/yaml/v3 v3.0.5
)

require (
	golang.org/x/mod v0.41.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/telemetry v0.0.0-20260908163034-4bcc4b2ee518 // indirect
	golang.org/x/tools v0.50.0 // indirect
	golang.org/x/vuln v1.8.0 // indirect
)

tool (
	go.uber.org/mock/mockgen
	golang.org/x/vuln/cmd/govulncheck
)
