script_dir=$(
  cd $(dirname $0)
  pwd
)                                  # 脚本路径
project_dir=$(dirname $script_dir) # 项目路径

proto_dir=proto
out_dir=${project_dir}/gen # 生成代码路径
third_party_dir=${project_dir}/third_party

protos=$(find ${project_dir}/${proto_dir} -type f -name '*.proto')
protoc \
  -I ${project_dir}/${proto_dir} \
  -I ${third_party_dir} \
  -I ${project_dir} \
  --go_out=${out_dir} \
  --go_opt paths=source_relative \
  --go-grpc_out ${out_dir} \
  --go-grpc_opt paths=source_relative \
  --dyn-gin_out ${out_dir} \
  --dyn-gin_opt paths=source_relative \
  --dyn-gin_opt use_encoding=true \
  --dyn-gin_opt enable_metadata=true \
  --dyn-resty_out ${out_dir} \
  --dyn-resty_opt paths=source_relative \
  --dyn-errno_out ${out_dir} \
  --dyn-errno_opt paths=source_relative \
  $protos
