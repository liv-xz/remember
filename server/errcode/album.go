package errcode

const (
	AlbumNotFound   int = 5001 // 相册不存在
	AlbumCreateFail int = 5002 // 创建相册失败
	AlbumDeleteFail int = 5003 // 删除相册失败
	PhotoNotFound   int = 5011 // 照片不存在
	PhotoUploadFail int = 5012 // 照片上传失败
	PhotoDeleteFail int = 5013 // 照片删除失败
	PhotoOverLimit  int = 5014 // 照片数量已达上限
)
