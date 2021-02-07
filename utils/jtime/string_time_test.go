package jtime

import (
	"encoding/json"
	"testing"
	"time"
)

type demoStruct struct {
	CreatedAt JsonTime `json:"createdAt"`
}

func TestJsonTimeToJson(t *testing.T) {
	nowtime := time.Now()
	t.Logf("createdAt:%s", nowtime.Format(TimeFormat))
	demo := demoStruct{
		CreatedAt: JsonTime{Time: nowtime},
	}
	timeJSON, err := json.Marshal(demo)
	if err != nil {
		t.Errorf("Marshal err:%s \n", err)
		return
	}
	t.Logf(" json.Marshal SUCCESS %s \n", string(timeJSON))
	var timeJsonToMap map[string]string
	err = json.Unmarshal(timeJSON, &timeJsonToMap)
	if err != nil {
		t.Errorf("Unmarshal err:%s \n", err)
		return
	}
	timeStr := timeJsonToMap["createdAt"]
	nowParse := nowtime.Format(TimeFormat)
	if nowParse != timeStr {
		t.Error("TestJsonTimeToJson FAIL")
		return
	}
	t.Log("TestJsonTimeToJson SUCCESS")
	return
}

func TestJsonToJsonTime(t *testing.T) {
	jsonData := `{"createdAt":"2021-02-07 14:33:47"}`
	createdtime, _ := time.ParseInLocation(TimeFormat, "2021-02-07 14:33:47", time.Local)
	var (
		err  error
		demo demoStruct
	)
	err = json.Unmarshal([]byte(jsonData), &demo)
	if err != nil {
		t.Error(err)
		return
	}
	if demo.CreatedAt.Unix() != createdtime.Unix() {
		t.Error("json.Unmarshal FAIL ")
		return
	}
	t.Logf("json.Unmarshal SUCCESS %+v \n", demo)
	return
}

func TestJsonTimeSetFormate(t *testing.T) {
	nowtime := time.Now()
	newFormat := "2006-01-02"
	t.Logf("createdAt:%s", nowtime.Format(newFormat))
	demo := demoStruct{
		CreatedAt: JsonTime{Time: nowtime},
	}
	demo.CreatedAt.SetFormat(newFormat)
	timeJSON, err := json.Marshal(demo)
	if err != nil {
		t.Errorf("Marshal err:%s \n", err)
		return
	}
	t.Logf(" json.Marshal SUCCESS %s \n", string(timeJSON))
	var timeJsonToMap map[string]string
	err = json.Unmarshal(timeJSON, &timeJsonToMap)
	if err != nil {
		t.Errorf("Unmarshal err:%s \n", err)
		return
	}
	timeStr := timeJsonToMap["createdAt"]
	nowParse := nowtime.Format(newFormat)
	if nowParse != timeStr {
		t.Error("TestJsonTimeSetFormate FAIL")
		return
	}
	t.Log("TestJsonTimeSetFormate SUCCESS")
	return
}
