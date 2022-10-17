package routes

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"gowebserver/src/models"
	"gowebserver/src/utils"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type BookGroup struct {
	DB          *sql.DB
	MinioClient *minio.Client
}

type BookCoverForm struct {
	Cover *multipart.FileHeader `form:"cover" binding:"required"`
}

func (b *BookGroup) create(c *gin.Context) {
	var book models.Book
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bk, err := models.CreateBook(b.DB, &book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, bk)
}

func (b *BookGroup) list(c *gin.Context) {
	books, err := models.AllBooks(b.DB)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, books)
}

func (b *BookGroup) uploadCover(c *gin.Context) {
	var form BookCoverForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	isbn := c.Param("isbn")
	_, err := models.FindBookByIsbn(b.DB, isbn)
	if err != nil {
		// TODO: Handle NotFound
		c.AbortWithError(500, err)
		return
	}

	err = b.uploadFile(
		form.Cover,
		fmt.Sprintf(`/images/book/%s/cover%s`, isbn, path.Ext(form.Cover.Filename)),
	)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, gin.H{})
}

func (b *BookGroup) uploadFile(file *multipart.FileHeader, filename string) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := utils.ResizeImage(f, 200, 200)
	if err != nil {
		return err
	}

	imgBytes, l, err := utils.ImageToBytes(img)
	if err != nil {
		return err
	}

	_, err = b.MinioClient.PutObject(
		context.Background(),
		"test",
		filename,
		bytes.NewReader(imgBytes),
		int64(l),
		minio.PutObjectOptions{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (b *BookGroup) registerGroup(rg *gin.RouterGroup) {
	group := rg.Group("/books")
	group.POST("/", b.create)
	group.GET("/", b.list)
	group.POST("/:isbn/cover", b.uploadCover)
}
