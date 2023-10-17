cd mpbdefinition;
protoc --gofast_out=../ ./db*.proto;
protoc --go_out=../mpb --go_opt=paths=source_relative  ./cmd.proto
protoc --go_out=../mpb --go_opt=paths=source_relative  ./aptos.proto
protoc --go_out=../mpb --go_opt=paths=source_relative  ./common.proto;
protoc --go_out=../mpb --go_opt=paths=source_relative  --go-grpc_out=../mpb --go-grpc_opt=paths=source_relative ./grpc*.proto;
protoc --go_out=../mpb --go_opt=paths=source_relative  ./errors.proto;
protoc --go_out=../mpb --go_opt=paths=source_relative  ./http*.proto;
protoc --go_out=../mpb --go_opt=paths=source_relative  ./resource.proto
cd ../;
