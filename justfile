export BINARY_NAME := "golox"

build:
    mkdir -p bin
    go build -o bin/${BINARY_NAME} main.go
    export PATH="$PATH:/Users/quang-dang/Workspaces/golox/bin"

test:
    go test ./internal/lox/

testv:
    go test -v ./internal/lox/

# Live reload > export PATH="$PATH:/Users/quang-dang/Workspaces/golox/tmp/bin"
watch:
    go run github.com/air-verse/air@v1.52.3 \
    --build.cmd "mkdir -p tmp/bin/ && go build -o tmp/bin/${BINARY_NAME}" \
    --build.bin "tmp/bin/${BINARY_NAME}" \
    --build.delay "100" \
    --build.exclude_dir "node_modules,sql,scripts,tests" \
    --build.include_ext "go" \
    --build.stop_on_error "false" \
    --misc.clean_on_exit true

gen-ast:
    go run scripts/generate_ast.go internal/lox/
    gofmt -w .
