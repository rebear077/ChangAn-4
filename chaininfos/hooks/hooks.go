package hooks

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// FileHook
//
//	@Description: 文件hook
type FileHook struct {
	//
	//  FileName
	//  @Description:  文件名
	//
	FileName string
}

func (h FileHook) Levels() []log.Level {
	return log.AllLevels
}
func (h FileHook) Fire(entry *log.Entry) error {
	fomat := &log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2022-07-17 00:00:00.000",
	}
	//  从entry从获得日志内容
	msg, err := fomat.Format(entry)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(h.FileName, msg, 0644)
	if err != nil {
		return err
	}
	return nil
}
