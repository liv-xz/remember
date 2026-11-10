package upload

import (
	"server/errcode"
	uploadsvc "server/service/upload"
	"server/utils/response"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploader uploadsvc.Uploader
}

func NewUploadHandler(uploader uploadsvc.Uploader) *UploadHandler {
	return &UploadHandler{uploader: uploader}
}

type UploadResp struct {
	URL string `json:"url"`
}

func (h *UploadHandler) Upload(c *gin.Context) {
	subDir := c.PostForm("subDir")
	if subDir == "" {
		subDir = "common"
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Fail(c, errcode.ParamInvalid)
		return
	}
	defer file.Close()

	url, err := h.uploader.Upload(file, header, subDir)
	if err != nil {
		errMsg := err.Error()
		switch {
		case errMsg == "文件大小超出限制":
			response.Fail(c, errcode.FileTooLarge)
		case errMsg == "不支持的文件类型":
			response.Fail(c, errcode.FileTypeDenied)
		default:
			response.Fail(c, errcode.FileUploadFail)
		}
		return
	}

	response.SuccessData(c, UploadResp{URL: url})
}
