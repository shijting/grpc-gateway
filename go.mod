module github.com/showiot/camera

go 1.15

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.572
	github.com/aliyun/aliyun-oss-go-sdk v2.1.5+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/envoyproxy/protoc-gen-validate v0.4.1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-pg/pg/v10 v10.5.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.0.1
	github.com/jinzhu/copier v0.0.0-20201025035756-632e723a6687
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/rakyll/statik v0.1.7
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.7.1
	github.com/zwczou/copier v0.1.0 // indirect
	golang.org/x/crypto v0.0.0-20201012173705-84dcc777aaee
	google.golang.org/genproto v0.0.0-20201021134325-0d71844de594
	google.golang.org/grpc v1.33.1
	google.golang.org/protobuf v1.25.0
	gopkg.in/ini.v1 v1.62.0 // indirect
)

replace github.com/jinzhu/copier v0.0.0-20201025035756-632e723a6687 => github.com/zwczou/copier v0.1.0
