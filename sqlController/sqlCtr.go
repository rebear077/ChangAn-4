package sql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rebear077/changan/conf"
	logloader "github.com/rebear077/changan/logs"
)

var logs = logloader.NewLog()

type SqlCtr struct {
	db        *sql.DB
	logDB     *sql.DB
	chainInfo *sql.DB
}

type LogData struct {
	Id        string
	Timestamp string
	Type      string
	Info      string
}

type ChainInfo struct {
	Id        string
	Timestamp string
	Type      string
	Title     string
	Info      string
}

func NewSqlCtr() *SqlCtr {
	configs, err := conf.ParseConfigFile("./configs/config.toml")
	if err != nil {
		// logrus.Fatalln(err)
		logs.Fatalln(err)
	}
	config := &configs[0]
	// db, err := sql.Open("mysql", "root:123456@/db_node0")
	str := config.MslUsername + ":" + config.MslPasswd + "@/" + config.MslName
	db, err := sql.Open("mysql", str)
	if err != nil {
		fmt.Println("err.Error()", err.Error())
		logs.Fatalln(err)
	}

	str1 := config.LogDBUsername + ":" + config.LogDBPasswd + "@/" + config.LogDBName
	logdb, err := sql.Open("mysql", str1)
	if err != nil {
		fmt.Println("err.Error()", err.Error())
		logs.Fatalln(err)
	}

	str2 := config.ChainInfoUsername + ":" + config.ChainInfoPasswd + "@/" + config.ChainInfoName
	chaininfo, err := sql.Open("mysql", str2)
	if err != nil {
		// logrus.Fatalln(err)
		logs.Fatalln(err)
	}

	return &SqlCtr{
		db:        db,
		logDB:     logdb,
		chainInfo: chaininfo,
	}
}

// 插入日志数据
func (s *SqlCtr) InsertLogs(Timestamp string, Type string, Info string) error {
	// _, err := s.logDB.Exec("CREATE TABLE IF NOT EXISTS `u_t_log`( `log_id` INT NOT NULL AUTO_INCREMENT, `timestamp` TIMESTAMP, `type` varchar(100) not null, `info` varchar(1000) not null,  primary key (`log_id`) )ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	// if err != nil {
	// 	// fmt.Printf("insert failed, err: %v\n", err)
	// 	return err
	// }
	// _, err = s.logDB.Exec("insert into u_t_log(timestamp, type, info) values (?,?,?)", Timestamp, Type, Info)
	// if err != nil {
	// 	// fmt.Printf("insert failed, err: %v\n", err)
	// 	return err
	// }
	return nil
}

// 插入区块连信息
func (s *SqlCtr) InserChainInfos(Timestamp string, Type string, Title string, Info string) error {
	// _, err := s.logDB.Exec("CREATE TABLE IF NOT EXISTS `u_t_chaininfo`( `log_id` INT NOT NULL AUTO_INCREMENT, `timestamp` TIMESTAMP, `type` varchar(100) not null, `title` varchar(1000) not null, `info` varchar(1000) not null,  primary key (`log_id`) )ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	// if err != nil {
	// 	// fmt.Printf("insert failed, err: %v\n", err)
	// 	return err
	// }
	// _, err = s.logDB.Exec("insert into u_t_chaininfo(timestamp, type, title, info) values (?,?,?, ?)", Timestamp, Type, Title, Info)
	// if err != nil {
	// 	// fmt.Printf("insert failed, err: %v\n", err)
	// 	return err
	// }
	return nil
}

// 查询日志数据
func (s *SqlCtr) QueryLogsByOrder(order string) interface{} {
	in_stmt, err := s.logDB.Prepare(order)
	if err != nil {
		// logrus.Panicln(err)
		return (err)
	}
	rows, err := in_stmt.Query()
	if err != nil {
		return err
	}
	switch order {
	case "select * from u_t_log":
		ret := make([]*LogData, 0)
		count := 0
		i := 0
		for rows.Next() {
			record := &LogData{}
			err = rows.Scan(&record.Id, &record.Timestamp, &record.Type, &record.Info)
			//fmt.Println(err)
			if err != nil {
				count++
				continue
			}
			ret = append(ret, record)
			fmt.Println(ret[i])
			i = i + 1
		}
		// logrus.Infof("select %d information from u_t_log and error counts is %x", len(ret), count)
		logs.Infof("select %d information from u_t_log and error counts is %x", len(ret), count)
		return ret
	}
	return nil
}

// 查询区块链信息
func (s *SqlCtr) QueryChaininfoByOrder(order string) interface{} {
	in_stmt, err := s.logDB.Prepare(order)
	if err != nil {
		// logrus.Panicln(err)
		logs.Panicln(err)
	}
	rows, err := in_stmt.Query()
	if err != nil {
		return err
	}
	switch order {
	case "select * from u_t_chaininfo":
		ret := make([]*ChainInfo, 0)
		count := 0
		i := 0
		for rows.Next() {
			record := &ChainInfo{}
			err = rows.Scan(&record.Id, &record.Timestamp, &record.Type, &record.Title, &record.Info)
			//fmt.Println(err)
			if err != nil {
				count++
				continue
			}
			ret = append(ret, record)
			fmt.Println(ret[i])
			i = i + 1
		}
		// logrus.Infof("select %d information from u_t_chainifo and error counts is %x", len(ret), count)
		logs.Infof("select %d information from u_t_chainifo and error counts is %x", len(ret), count)
		return ret
	}
	return nil
}

// 根据时间戳，获取当前时间前一小时以内的数据
func (s *SqlCtr) ExtractLogByHour() (interface{}, int) {
	currentTime := time.Now()         //当前时间
	m, _ := time.ParseDuration("-1h") //当前时间往前推1小时
	checktime := currentTime.Add(m)
	checktime_str := checktime.String()
	timeTemplate1 := "2006-01-02 15:04:05"                                                //改变时间戳模板， 方便与下面的recordtime比较
	checktime1, _ := time.ParseInLocation(timeTemplate1, checktime_str[0:19], time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
	// fmt.Println("checktime is: ", checktime1)
	in_stmt, err := s.logDB.Prepare("select * from u_t_log")
	if err != nil {
		// logrus.Panicln(err)
		logs.Panicln(err)
	}
	rows, err := in_stmt.Query()
	if err != nil {
		return err, 0
	}
	ret := make([]*LogData, 0)
	count := 0
	i := 0
	fmt.Printf("最近一小时以内的数据: \n")
	for rows.Next() {
		record := &LogData{}
		err = rows.Scan(&record.Id, &record.Timestamp, &record.Type, &record.Info)
		//fmt.Println(err)
		if err != nil {
			count++
			continue
		}
		recordtime, _ := time.ParseInLocation(timeTemplate1, record.Timestamp, time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
		// fmt.Println("recordtime is: ", recordtime)
		if recordtime.After(checktime1) {
			ret = append(ret, record)
			fmt.Println(ret[i])
			i = i + 1
		}
	}
	return ret, i
}
