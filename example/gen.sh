scriptDir=$(
  cd $(dirname $0)
  pwd
)                             # 脚本路径
projDir=$(dirname $scriptDir) # 项目路径

protoDir=example
outDir=${projDir}/example # 生成代码路径
thirdPartyDir=${projDir}/example/third_party

protoc \
  -I ${projDir}/${protoDir} \
  -I ${thirdPartyDir} \
  -I ${projDir} \
  --go_out=${outDir} \
  --go_opt paths=source_relative \
  --go-grpc_out ${outDir} \
  --go-grpc_opt paths=source_relative \
  --dyn-gin_out ${outDir} \
  --dyn-gin_opt paths=source_relative \
  --dyn-resty_out ${outDir} \
  --dyn-resty_opt paths=source_relative \
  hello.proto
