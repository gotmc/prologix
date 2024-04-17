# List the available justfile recipes.
@default:
  just --list

# Format, vet, and test Go code.
check:
	go fmt ./...
	go vet ./...
	GOEXPERIMENT=loopvar go test ./... -cover

# Verbosely format, vet, and test Go code.
checkv:
	go fmt ./...
	go vet ./...
	GOEXPERIMENT=loopvar go test -v ./... -cover

# Lint code using staticcheck.
lint:
	staticcheck -f stylish ./...

# Test and provide HTML coverage report.
cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# List the outdated go modules.
outdated:
  go list -u -m all

# Run the Prologix VCP GPIB Keysight E3631A example application.
key3631 port gpib:
  #!/usr/bin/env bash
  echo '# Prologix VCP GPIB Keysight E3631A Example Application'
  cd {{justfile_directory()}}/examples/vcp/e3631a
  env go build -o e3631a
  ./e3631a -port={{port}} -gpib={{gpib}}

# Run the Prologix VCP GPIB Keysight 33220A example application.
key33220 port gpib:
  #!/usr/bin/env bash
  echo '# Prologix VCP GPIB Keysight 33220A Example Application'
  cd {{justfile_directory()}}/examples/vcp/key33220a
  env go build -o key33220a
  ./key33220a -port={{port}} -gpib={{gpib}}

# Run the Prologix VCP GPIB SRS DS345 example application.
ds345 port gpib:
  #!/usr/bin/env bash
  echo '# Prologix VCP GPIB SRS DS345 Example Application'
  cd {{justfile_directory()}}/examples/vcp/ds345
  env go build -o ds345
  ./ds345 -port={{port}} -gpib={{gpib}}
