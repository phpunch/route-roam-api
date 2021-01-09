module github.com/phpunch/route-roam-api

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/gomodule/redigo v1.8.3 // indirect
	github.com/jackc/pgx/v4 v4.10.1
	github.com/lib/pq v1.3.0
	github.com/minio/minio-go/v7 v7.0.6
	github.com/spf13/viper v1.7.1
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/postgres v1.0.6
	gorm.io/gorm v1.20.9
)
