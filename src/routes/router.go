package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type Env struct {
	Db          *sql.DB
	MinioClient *minio.Client
}

func NewRouter(env *Env) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")

	bookGroup := BookGroup{env.Db, env.MinioClient}
	bookGroup.registerGroup(v1)

	return r
}
