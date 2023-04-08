package serialization

//go:generate protoc --go_out=. generictypes.proto
//go:generate protoc-go-inject-tag -input generictypes.pb.go
