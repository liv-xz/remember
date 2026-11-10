package errcode

const (
	FileTooLarge   int = 6001 // 文件大小超出限制
	FileTypeDenied int = 6002 // 文件类型不允许
	FileUploadFail int = 6003 // 文件上传失败
	FileDeleteFail int = 6004 // 文件删除失败
)
