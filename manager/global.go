package manager

//go:generate protoc --go_out=. --go_opt=Mmanager.proto=../manager manager.proto
//go:generate protoc --go_out=. --go_opt=Mmanager.proto=../manager manager_list.proto
