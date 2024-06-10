protoc --proto_path=./ \
        --go_out=../ \
        --go_opt=paths=source_relative \
        --go_opt=Mmetrics_model.proto=github.com/c12s/metrics/pkg/api \
        metrics_model.proto

protoc --proto_path=./ \
        --go_out=../ \
        --go_opt=paths=source_relative \
        --go-grpc_out=../ \
        --go-grpc_opt=paths=source_relative \
        --go_opt=Mmetrics.proto=github.com/c12s/metrics/pkg/api \
        --go-grpc_opt=Mmetrics.proto=github.com/c12s/metrics/pkg/api \
        -I ./metrics_model.proto \
        metrics.proto
