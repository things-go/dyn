scriptDir=$(cd $(dirname $0);pwd) # 脚本路径
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
  --go-gin_out ${outDir} \
  --go-gin_opt paths=source_relative \
  --go-gin_opt rpc_mode=official \
  --go-grpc_out ${outDir} \
  --go-grpc_opt paths=source_relative \
  --go-enum_out ${outDir} \
  --go-enum_opt paths=source_relative \
  --go-enum_opt merge=true \
  --go-enum_opt filename=typing \
  --go-enum_opt package=examples \
  --go-enum_opt go_package="github.com/things-go/examples" \
  hello.proto enum.proto

protoc \
  -I ${projDir}/${protoDir} \
  -I ${thirdPartyDir} \
  -I ${projDir} \
  --go-enum_out ${outDir} \
  --go-enum_opt suffix=".validate.pb.go" \
  --go-enum_opt template=${projDir}/${protoDir}/custom_template.tpl \
  --go-enum_opt paths=source_relative \
  --go-enum_opt merge=true \
  --go-enum_opt filename=typing \
  --go-enum_opt package=examples \
  --go-enum_opt go_package="github.com/things-go/examples" \
  hello.proto enum.proto