module github.com/haniz/byakugan/scanner

go 1.21

require (
	github.com/stretchr/testify v1.8.4
	google.golang.org/grpc v1.62.1
	google.golang.org/protobuf v1.36.5
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/haniz/byakugan/scanner/engine => ./engine
	github.com/haniz/byakugan/scanner/config => ./config
)
