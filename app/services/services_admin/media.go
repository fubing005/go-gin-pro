package services_admin

import (
	"context"
	"errors"
	"path"
	"shalabing-gin/app/common/request/request_admin"
	"shalabing-gin/app/models"
	"shalabing-gin/core/trans"
	"shalabing-gin/global"
	"strconv"
	"time"

	"shalabing-gin/go-storage/storage"

	uuid "github.com/satori/go.uuid"
)

type mediaService struct {
}

var MediaService = new(mediaService)

type outPut struct {
	Id   int64  `json:"id"`
	Path string `json:"path"`
	Url  string `json:"url"`
}

const mediaCacheKeyPre = "media:"

func (mediaService *mediaService) makeFaceDir(business string) string {
	return global.App.Config.App.Env + "/" + business
}

func (mediaService *mediaService) HashName(fileName string) string {
	fileSuffix := path.Ext(fileName)
	return uuid.NewV4().String() + fileSuffix
}

// SaveImage 保存图片（公共读）
func (mediaService *mediaService) SaveImage(params request_admin.ImageUpload) (result outPut, err error) {
	file, err := params.Image.Open()

	if err != nil {
		err = errors.New(trans.Trans("common.上传失败"))
		return
	}
	defer file.Close()

	localPrefix := ""
	if global.App.Config.Storage.Default == storage.Local {
		localPrefix = "public" + "/"
	}
	key := mediaService.makeFaceDir(params.Business) + "/" + mediaService.HashName(params.Image.Filename)
	disk := global.App.Disk()
	err = disk.Put(localPrefix+key, file, params.Image.Size)
	if err != nil {
		return
	}

	image := models.Media{
		DiskType: string(global.App.Config.Storage.Default),
		SrcType:  1,
		Src:      key,
	}
	err = global.App.DB.Create(&image).Error
	if err != nil {
		return
	}

	result = outPut{int64(image.ID.ID), key, disk.Url(key)}
	return
}

func (mediaService *mediaService) GetUrlById(id int64) string {
	if id == 0 {
		return ""
	}

	var url string
	cacheKey := mediaCacheKeyPre + strconv.FormatInt(id, 10)

	exist := global.App.Redis.Exists(context.Background(), cacheKey).Val()
	if exist == 1 {
		url = global.App.Redis.Get(context.Background(), cacheKey).Val()
	} else {
		media := models.Media{}
		err := global.App.DB.First(&media, id).Error
		if err != nil {
			return ""
		}
		url = global.App.Disk(media.DiskType).Url(media.Src)
		global.App.Redis.Set(context.Background(), cacheKey, url, time.Second*3*24*3600)
	}

	return url
}
