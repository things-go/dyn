#!/bin/bash

script_dir=$(
    cd $(dirname $0)
    pwd
)                                             # 脚本路径
project_dir=$(dirname $(dirname $script_dir)) # 项目路径

proto_dir=${project_dir}/examples/enums
out_dir=${project_dir}/examples/enums # 生成代码路径
third_party_dir=../../../examples/third_party

protoc \
    -I ${proto_dir} \
    -I ${third_party_dir} \
    --go_out=${out_dir} \
    --go_opt paths=source_relative \
    --dyn-enum_out ${out_dir} \
    --dyn-enum_opt paths=source_relative \
    nested.proto \
    non_nested.proto

protoc \
    -I ${proto_dir} \
    -I ${third_party_dir} \
    --dyn-enum_out ${out_dir} \
    --dyn-enum_opt suffix=".example.pb.go" \
    --dyn-enum_opt template=${proto_dir}/mapper_template.tpl \
    --dyn-enum_opt paths=source_relative \
    --dyn-enum_opt merge=true \
    --dyn-enum_opt filename=mapper \
    --dyn-enum_opt package=enums \
    --dyn-enum_opt go_package="github.com/things-go/examples/enums" \
    nested.proto \
    non_nested.proto
