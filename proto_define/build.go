package proto_define

//go:generate protoc --proto_path=./ --go_out=./ --go_opt=paths=source_relative ./common_data.proto
//go:generate protoc --proto_path=./ --go_out=./ --go_opt=paths=source_relative ./pgsql_service.proto
//go:generate protoc --proto_path=./ --go_out=./ --go_opt=paths=source_relative ./user_service.proto
