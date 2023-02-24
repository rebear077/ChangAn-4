package loader

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// 记录mysql中的日志数量
var log_num int

type Loader struct {
	sql         *sql.DB
	logpath     []string // 全局变量，存放日志名称（以时间命名）, 记录到mysql之后，从列表中剔除对应的日志
	currentTime time.Time
}
type LogData struct {
	Id        string
	Timestamp string
	Type      string
	Info      string
}

// 生成log数据库对象
func NewLoader() *Loader {
	// db, err := sql.Open("mysql", "root:123456@/db_node0")
	//重要操作：注册文件！！！

	//之后写在配置文件里
	str := "root" + ":" + "123456" + "@/" + "log_test" + "?allowAllFiles=true"
	db, err := sql.Open("mysql", str)
	createTable(db)
	var pathstring []string

	currentTime := time.Now() //开始运行时间
	currentsecond := time.Now().Second()
	minussecond, _ := time.ParseDuration("-" + strconv.Itoa(currentsecond) + "s") //清除秒，从xx:xx:00开始
	currentTime = currentTime.Add(minussecond)

	pathstring = append(pathstring, TimeStringFormatter(currentTime))
	if err != nil {
		fmt.Println(2)
		logrus.Fatalln(err)
	}
	return &Loader{
		sql:         db,
		logpath:     pathstring,
		currentTime: currentTime,
	}
}

func createTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `u_t_log`( `log_id` INT NOT NULL AUTO_INCREMENT, `timestamp` varchar(1000), `type` varchar(100) not null, `info` varchar(1000) not null,  primary key (`log_id`) )ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	if err != nil {
		fmt.Println(err)
	}
}

func (l *Loader) Start() {
	for {
		l.appendLogpath()
		// 如果当前存放日志文件名称的数组长度大于等于2
		// 说明至少有一个日志文件已经记录完毕，需要插入到mysql数据库中
		if len(l.logpath) >= 2 {
			l.InsertLogs()
		}
	}
}

func (l *Loader) initLogid() {
	id_str, err := l.maxId()
	// fmt.Println("id_str", id_str)
	if err != nil {
		fmt.Println(err)
	}
	id, err := strconv.Atoi(id_str)
	// fmt.Println("id", id)
	if err != nil {
		fmt.Println(err)
	}
	log_num = id
	// fmt.Println(log_num)
}

// 此函数的目的是考虑到，如果中途停止了程序，那么重新启动后，应该从当前数据库的最大id开始标注
func (l *Loader) maxId() (string, error) {
	result, err := l.sql.Prepare("select * from u_t_log where log_id = (select max(log_id) from u_t_log);")
	if err != nil {
		return "-1", err
	}
	rows, err := result.Query()
	if err != nil {
		return "-2", err
	}
	exist := false
	record := &LogData{}
	for rows.Next() {
		exist = true
		err = rows.Scan(&record.Id, &record.Timestamp, &record.Type, &record.Info)
		if err != nil {
			return "-3", err
		}
	}
	if !exist {
		return "0", nil
	}
	// fmt.Println(record.Id)
	return record.Id, nil
}

// 根据时间,向logpath中添加日志文件名称
func (l *Loader) appendLogpath() {
	m, _ := time.ParseDuration("+1m")
	checktime := time.Now()
	temptime := l.currentTime.Add(m)
	if checktime.After(temptime) {
		temptime_str := TimeStringFormatter(temptime)
		l.logpath = append(l.logpath, temptime_str)
		l.currentTime = temptime
	}
}

// 除了最后一个，插入日志数据至mysql
func (l *Loader) InsertLogs() error {
	// if len(logpath) >= 2 {
	// fmt.Println("l.logpath: ", l.logpath)
	lasttime := l.logpath[len(l.logpath)-1]
	// fmt.Println("lasttime: ", lasttime)
	for i := 0; i < len(l.logpath)-1; i++ {
		path := "output.log" + "." + l.logpath[i]
		// fmt.Println("path: ", path)
		l.sql.Exec("load data local infile '/home/jackson/ChangAn-1/logs/log_files/" + path + "' into table u_t_log CHARACTER SET utf8 fields terminated by '|' lines terminated by '\n'")
	}
	l.logpath = nil
	l.logpath = append(l.logpath, lasttime)
	return nil
}

// 时间戳转格式化字符串
func TimeStringFormatter(temptime time.Time) string {
	currentTime_str := temptime.String()
	timeTemplate1 := "2006-01-02 15:04:05"                                                  //改变时间戳模板， 方便与下面的recordtime比较
	checktime1, _ := time.ParseInLocation(timeTemplate1, currentTime_str[0:19], time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
	checktime2 := checktime1.String()[0:4] + checktime1.String()[5:7] + checktime1.String()[8:10] + checktime1.String()[11:13] + checktime1.String()[14:16]
	return checktime2
}
