package patch_data

//go:generate protoc --proto_path=./ --go_out=./ --go_opt=paths=source_relative ./primary-entities.proto
//go:generate protoc --proto_path=./ --go_out=./ --go_opt=paths=source_relative ./enums.proto
//go:generate protoc --proto_path=./ --go_out=./ --go_opt=paths=source_relative ./entities.proto
//go:generate protoc --proto_path=./ --go_out=./ --go_opt=paths=source_relative ./envelopes.proto
