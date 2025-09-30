generate-proto-metrics:
	@mkdir -p ./pb/metrics
	@protoc -I ./proto/metrics \
        --go_out=./pb/metrics  --go_opt=paths=source_relative \
        --go-grpc_out=./pb/metrics  --go-grpc_opt=paths=source_relative \
        proto/metrics/*.proto
	@echo 'Success generate proto Metrics'