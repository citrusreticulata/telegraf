//go:generate ../../../tools/readme_config_includer/generator
package iotdb

// iotdb.go

import (
	_ "embed"
	"errors"
	"fmt"
	"math"

	// Register IoTDB go client
	"github.com/apache/iotdb-client-go/client"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
)

// DO NOT REMOVE THE NEXT TWO LINES! This is required to embed the sampleConfig data.
//go:embed sample.conf
var sampleConfig string

type IoTDB struct {
	Host            string `toml:"host"`
	Port            string `toml:"port"`
	User            string `toml:"user"`
	Password        string `toml:"password"`
	Timeout         int    `toml:"timeout"`
	ConvertUint64To string `toml:"convertUint64To"`
	TimeStampUnit   string `toml:"timeStampUnit"`
	session         *client.Session

	Log telegraf.Logger `toml:"-"`
}

func (*IoTDB) SampleConfig() string {
	return sampleConfig
}

// Init is for setup, and validating config.
func (s *IoTDB) Init() error {
	var errorMsg string
	if s.Timeout < 0 {
		errorMsg = fmt.Sprintf("IoTDB Config Warning: The value of 'timeout' is negative:%d. Now it's fixed to 0.", s.Timeout)
		s.Log.Warnf(errorMsg)
		// return errors.New(errorMsg)
		s.Timeout = 0
	}
	if !(s.ConvertUint64To == "ToInt64" ||
		s.ConvertUint64To == "ForceToInt64" ||
		s.ConvertUint64To == "Text") {
		errorMsg = fmt.Sprintf("IoTDB Config Warning: The value of 'ConvertUint64To' is invaild: %s. Now it's fixed to 'ToInt64'.", s.ConvertUint64To)
		s.Log.Warnf(errorMsg)
		s.ConvertUint64To = "ToInt64"
	}
	if !(s.TimeStampUnit == "second" ||
		s.TimeStampUnit == "millisecond" ||
		s.TimeStampUnit == "microsecond" ||
		s.TimeStampUnit == "nanosecond") {
		errorMsg = fmt.Sprintf("IoTDB Config Warning: The value of 'TimeStampUnit' is invaild: %s. Now it's fixed to 'nanosecond'.", s.TimeStampUnit)
		s.Log.Warnf(errorMsg)
		s.TimeStampUnit = "nanosecond"
	}
	return nil
}

func (s *IoTDB) Connect() error {
	// Make any connection required here
	// Check the configuration
	config := &client.Config{
		Host:     s.Host,
		Port:     s.Port,
		UserName: s.User,
		Password: s.Password,
	}
	var ss = client.NewSession(config)
	s.session = &ss
	if err := s.session.Open(false, s.Timeout); err != nil {
		s.Log.Errorf("IoTDB Connect Error: Fail to connect host:'%s', port:'%s', err:%v", s.Host, s.Port, err)
		return err
	}

	return nil
}

func (s *IoTDB) Close() error {
	// Close any connections here.
	// Write will not be called once Close is called, so there is no need to synchronize.
	_, err := s.session.Close()
	if err != nil {
		s.Log.Errorf("IoTDB Close Error: %v", err)
	}
	return nil
}

// Write should write immediately to the output, and not buffer writes
// (Telegraf manages the buffer for you). Returning an error will fail this
// batch of writes and the entire batch will be retried automatically.
func (s *IoTDB) Write(metrics []telegraf.Metric) error {

	var deviceId_list []string
	var measurements_list [][]string
	var values_list [][]interface{}
	var dataTypes_list [][]client.TSDataType
	var timestamp_list []int64

	for _, metric := range metrics {
		// write `metric` to the output sink here
		var keys []string
		var values []interface{}
		var dataTypes []client.TSDataType
		for _, tag := range metric.TagList() {
			datatype, value := s.getDataTypeAndValue(tag.Value)
			if datatype != client.UNKNOW {
				keys = append(keys, tag.Key)
				values = append(values, value)
				dataTypes = append(dataTypes, datatype)
			}
		}
		for _, field := range metric.FieldList() {
			//pk = append(pk, quoteIdent(tag.Key))
			datatype, value := s.getDataTypeAndValue(field.Value)
			if datatype != client.UNKNOW {
				keys = append(keys, field.Key)
				values = append(values, value)
				dataTypes = append(dataTypes, datatype)
			}
		}
		if s.TimeStampUnit == "second" {
			timestamp_list = append(timestamp_list, metric.Time().Unix())
		} else if s.TimeStampUnit == "millisecond" {
			timestamp_list = append(timestamp_list, metric.Time().UnixMilli())
		} else if s.TimeStampUnit == "microsecond" {
			timestamp_list = append(timestamp_list, metric.Time().UnixMicro())
		} else if s.TimeStampUnit == "nanosecond" {
			timestamp_list = append(timestamp_list, metric.Time().UnixNano())
		} else {
			var errorMsg string
			errorMsg = fmt.Sprintf("IoTDB Configuration Error: unknown TimeStampUnit: %s", s.TimeStampUnit)
			s.Log.Errorf(errorMsg)
			return errors.New(errorMsg)
		}
		// append other metric date to list
		deviceId_list = append(deviceId_list, metric.Name())
		measurements_list = append(measurements_list, keys)
		values_list = append(values_list, values)
		dataTypes_list = append(dataTypes_list, dataTypes)
	}
	// Wirte to client
	status, err := s.session.InsertRecords(deviceId_list, measurements_list, dataTypes_list, values_list, timestamp_list)
	if err != nil {
		s.Log.Errorf(err.Error())
	}
	if status != nil {
		if verifyResult := client.VerifySuccess(status); verifyResult != nil {
			s.Log.Info(verifyResult)
		}
	}
	return err
}

// Find out data type of the value and return it's id in TSDataType, and convert it if nessary.
func (s *IoTDB) getDataTypeAndValue(value interface{}) (client.TSDataType, interface{}) {
	switch v := value.(type) {
	case int32:
		return client.INT32, int32(v)
	case int64:
		return client.INT64, int64(v)
	case uint32:
		return client.INT64, int64(v)
	case uint64:
		if s.ConvertUint64To == "ToInt64" {
			if v <= uint64(math.MaxInt64) {
				return client.INT64, int64(v)
			} else {
				return client.INT64, int64(math.MaxInt64)
			}
		} else if s.ConvertUint64To == "ForceToInt64" {
			return client.INT64, int64(v)
		} else if s.ConvertUint64To == "Text" {
			return client.TEXT, fmt.Sprintf("%d", value)
		} else {
			s.Log.Errorf("unknown converstaion configuration of 'convertUint64To': %s", s.ConvertUint64To)
			return client.UNKNOW, int64(0)
		}
	case float64:
		return client.DOUBLE, float64(v)
	case string:
		return client.TEXT, v
	case bool:
		return client.BOOLEAN, v
	default:
		s.Log.Errorf("Unknown datatype: '%T' %v", value, value)
		return client.UNKNOW, int64(0)
	}
	// s.Log.Errorf("function IoTDB.deriveDatType reaches unreachable area.")
	// return client.UNKNOW, int64(0)
}

func init() {
	outputs.Add("iotdb", func() telegraf.Output { return newIoTDB() })
}

func newIoTDB() *IoTDB {
	return &IoTDB{}
}
