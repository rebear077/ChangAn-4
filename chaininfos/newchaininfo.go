package chaininfo

import (
	"bytes"
	"fmt"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type ChaininfoFormatter struct{}

func NewChainlog() *logrus.Logger {
	NewChainInfo().initLogid()               //初始化，获得当前数据库日志中最大id
	mLog := logrus.New()                     //新建一个实例
	mLog.SetOutput(os.Stdout)                //设置输出类型
	mLog.SetReportCaller(true)               //开启返回函数名和行号
	mLog.SetFormatter(&ChaininfoFormatter{}) //设置自己定义的Formatter
	mLog.SetLevel(logrus.DebugLevel)         //设置最低的Level
	path := "chaininfos/chaininfo_files/output.log"
	/* 日志轮转相关函数
	`WithLinkName` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	WithMaxAge 和 WithRotationCount二者只能设置一个
	  `WithMaxAge` 设置文件清理前的最长保存时间
	  `WithRotationCount` 设置文件清理前最多保存的个数
	*/
	// 下面配置日志每隔 1 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(1200)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
	)
	mLog.SetOutput(writer) // 将文件设置为log输出的文件
	return mLog
}

// Format 实现Formatter(entry *logrus.Entry) ([]byte, error)接口
func (f *ChaininfoFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	//自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	log_num = log_num + 1
	// fmt.Println("log_num1: ", log_num)
	// fmt.Println("log_num: ", log_num)
	if entry.HasCaller() {
		//自定义输出格式
		fmt.Fprintf(b, "%d|[%s]|[%s]|[%s]\n", log_num, timestamp, entry.Level, entry.Message)
	} else {
		fmt.Fprintf(b, "%d|[%s]|[%s]|[%s]\n", log_num, timestamp, entry.Level, entry.Message)
	}
	return b.Bytes(), nil
}
