protoc --proto_path=./ \
        --go_out=../ \
        --go_opt=paths=source_relative \
        --go_opt=Mexternal_metrics_model.proto=github.com/c12s/metrics/pkg/external \
        external_metrics_model.proto

protoc --proto_path=./ \
        --go_out=../ \
        --go_opt=paths=source_relative \
        --go-grpc_out=../ \
        --go-grpc_opt=paths=source_relative \
        --go_opt=Mexternal_metrics.proto=github.com/c12s/metrics/pkg/external \
        --go-grpc_opt=Mexternal_metrics.proto=github.com/c12s/metrics/pkg/external \
        -I ./external_metrics.proto \
        external_metrics.proto
