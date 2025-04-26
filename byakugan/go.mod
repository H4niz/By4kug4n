module github.com/haniz/byakugan/byakugan

go 1.21

replace (
	github.com/haniz/byakugan => ../
	github.com/haniz/byakugan/cmd/cli => ./cmd/cli
	github.com/haniz/byakugan/core => ./core
	github.com/haniz/byakugan/proto => ./proto
	github.com/haniz/byakugan/scanner => ./scanner
)
