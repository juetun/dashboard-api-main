// Package web /**
package web

const (
	RedisCacheKeyPrefixExport = "export:"
)

const (
	// RunProgressInit 导出任务特殊进度标记
	RunProgressInit     = 1
	RunProgressStart    = 2  // 开始获取第一页的进度
	RunProgressFistPage = 8  // 第一页数据获取成功后的进度
	RunProgressMax      = 98 // 获取完所有数据时的进度

	ExportLoading = 0 // 导出任务进行中
	ExportFailure = 2 // 导出失败
	ExportSuccess = 1 // 导出成功
	ExportExpire  = 3 // 导出超时
)
