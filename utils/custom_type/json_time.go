package custom_type

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

//JsonTime 适用用于json序列化的自定义时间类型数据
type JsonTime struct {
	formatString string
	time.Time
}

//MarshalJSON 实现json.Marshaler 的接口，将time.Time类型解析为指定格式的字符串
func (t JsonTime) MarshalJSON() ([]byte, error) {
	if t.formatString == "" {
		t.formatString = `"` + TimeFormat + `"`
	}
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return []byte(""), nil
	}
	tune := t.Local().Format(t.formatString)
	return []byte(tune), nil
}

//UnmarshalJSON 实现json.UnmarshalJSON 的接口，将指定格式的字符串解析为time.Time类型
func (t *JsonTime) UnmarshalJSON(data []byte) error {
	if t.formatString == "" {
		t.formatString = `"` + TimeFormat + `"`
	}

	if len(data) == 0 {
		*t = JsonTime{
			Time: time.Time{},
		}
		return nil
	}
	now, err := time.ParseInLocation(t.formatString, string(data), time.Local)
	*t = JsonTime{
		Time: now,
	}
	return err
}

//SetFormat 设定进行json化时需要的时间格式
func (t *JsonTime) SetFormat(format string) {
	t.formatString = fmt.Sprintf("\"%s\"", format)
}

//Value 实现driver.Valuer接口，用于数据存入数据库
func (t JsonTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	//如果时间是零值，则存入数据库的值为null
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

//Scan 实现sql.Scanner接口， 从数据库中取出数据
func (t *JsonTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JsonTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
