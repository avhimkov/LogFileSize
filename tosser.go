package LogFileSize

type TosserStat struct {
	Dates map[string]map[string]*dirStatInfo
	ConfigName string
}
type dirStatInfo struct {
	//количество файлов
	Count int64
	//дата последней передачи файла
	LastProcessingDate int64
	//общий размер файлов
	TotalSize int64
}