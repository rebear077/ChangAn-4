package chaininfo

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

type ChainInfo struct {
	db          *sql.DB
	infopath    []string // 全局变量，存放日志名称（以时间命名）, 记录到mysql之后，从列表中剔除对应的日志
	currentTime time.Time
}

type ChainData struct {
	Id        string
	Timestamp string
	Type      string
	Info      string
}

// func init() {
// 	NewChainInfo()
// }

// 生成log数据库对象
func NewChainInfo() *ChainInfo {
	// db, err := sql.Open("mysql", "root:123456@/db_node0")
	//重要操作：注册文件！！！

	//之后写在配置文件里
	str := "root" + ":" + "123456" + "@/" + "chaininfo_test" + "?allowAllFiles=true"
	db, err := sql.Open("mysql", str)
	createTable(db)
	var pathstring []string

	currentTime := time.Now() //开始运行时间
	currentsecond := time.Now().Second()
	minussecond, _ := time.ParseDuration("-" + strconv.Itoa(currentsecond) + "s") //清除秒，从xx:xx:00开始
	currentTime = currentTime.Add(minussecond)

	pathstring = append(pathstring, TimeStringFormatter(currentTime))
	if err != nil {
		logrus.Fatalln(err)
	}
	return &ChainInfo{
		db:          db,
		infopath:    pathstring,
		currentTime: currentTime,
	}
}

func (c *ChainInfo) Start() {
	for {
		c.appendLogpath()
		// 如果当前存放日志文件名称的数组长度大于等于2
		// 说明至少有一个日志文件已经记录完毕，需要插入到mysql数据库中
		if len(c.infopath) >= 2 {
			c.InsertLogs()
		}
	}
}

func createTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `u_t_chaininfo`( `id` INT NOT NULL AUTO_INCREMENT, `timestamp` varchar(1000), `type` varchar(100) not null, `info` varchar(1000) not null,  primary key (`id`) )ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	if err != nil {
		fmt.Println(err)
	}
}

func (c *ChainInfo) initLogid() {
	id_str, err := c.maxId()
	if err != nil {
		fmt.Println(err)
	}
	id, err := strconv.Atoi(id_str)
	if err != nil {
		fmt.Println(err)
	}
	log_num = id
}

// 此函数的目的是考虑到，如果中途停止了程序，那么重新启动后，应该从当前数据库的最大id开始标注
func (c *ChainInfo) maxId() (string, error) {
	result, err := c.db.Prepare("select * from u_t_chaininfo where id = (select max(id) from u_t_chaininfo);")
	if err != nil {
		return "-1", err
	}
	rows, err := result.Query()
	if err != nil {
		return "-2", err
	}
	exist := false
	record := &ChainData{}
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
	return record.Id, nil
}

func (c *ChainInfo) appendLogpath() {
	m, _ := time.ParseDuration("+1m")
	checktime := time.Now()
	temptime := c.currentTime.Add(m)
	if checktime.After(temptime) {
		temptime_str := TimeStringFormatter(temptime)
		c.infopath = append(c.infopath, temptime_str)
		c.currentTime = temptime
	}
}

// 除了最后一个，插入日志数据至mysql
func (c *ChainInfo) InsertLogs() error {
	lasttime := c.infopath[len(c.infopath)-1]
	for i := 0; i < len(c.infopath)-1; i++ {
		path := "output.log" + "." + c.infopath[i]
		// fmt.Println("chaininfo path: ", path)
		_, err := c.db.Exec("load data local infile '/home/jackson/ChangAn-3/chaininfos/chaininfo_files/" + path + "' into table u_t_chaininfo CHARACTER SET utf8 fields terminated by '|' lines terminated by '\n'")
		if err != nil {
			fmt.Println(err)
			// return err
		}
	}
	c.infopath = nil
	c.infopath = append(c.infopath, lasttime)
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
