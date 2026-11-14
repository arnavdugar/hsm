module github.com/arnavdugar/hsm

go 1.27

require (
	github.com/microsoft/TypeScript/tsc/ast v0.0.0
	github.com/stretchr/testify v1.12.1
	go.uber.org/mock v0.6.0
	go.yaml.in/yaml/v3 v3.0.6
)

require (
	github.com/klauspost/compress v1.20.3 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/microsoft/TypeScript/tsc v0.0.0-20261116080028-5e32f3249529 // indirect
	github.com/zeebo/xxh3 v1.1.0 // indirect
	golang.org/x/mod v0.42.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/telemetry v0.0.0-20261010193011-f71dc92da45e // indirect
	golang.org/x/text v0.43.0 // indirect
	golang.org/x/tools v0.52.0 // indirect
	golang.org/x/vuln v1.9.0 // indirect
)

tool (
	go.uber.org/mock/mockgen
	golang.org/x/vuln/cmd/govulncheck
)

replace github.com/microsoft/TypeScript/tsc/ast => ./codegen/typescript/tsc
