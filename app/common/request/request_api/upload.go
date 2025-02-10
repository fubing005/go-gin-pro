package request_api

import (
	"mime/multipart"
	"shalabing-gin/app/common/request"
	"shalabing-gin/core/trans"
)

type ImageUpload struct {
	Business string                `form:"business" json:"business" binding:"required"`
	Image    *multipart.FileHeader `form:"image" json:"image" binding:"required"`
}

func (imageUpload ImageUpload) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"business.required": trans.Trans("common.业务类型不能为空"),
		"image.required":    trans.Trans("common.请选择图片"),
	}
}
