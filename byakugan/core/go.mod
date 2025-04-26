module github.com/haniz/byakugan/core

go 1.21

require (
    github.com/haniz/byakugan/proto v0.0.0
    google.golang.org/grpc v1.62.1
)

replace github.com/haniz/byakugan/proto => ../proto