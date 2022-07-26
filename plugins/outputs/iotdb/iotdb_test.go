//go:generate ../../../tools/readme_config_includer/generator
package iotdb

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/apache/iotdb-client-go/client"
	"github.com/stretchr/testify/require"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/testutil"
)

var (
	target_host     = "192.168.134.130" // The server's ip that you want to connect to.
	target_port     = "6667"            // The server's port that you want to connect to.
	target_user     = "root"
	target_password = "root"
)

func TestConnectAndClose(t *testing.T) {
	test_client := &IoTDB{
		Host:     target_host,
		Port:     target_port,
		User:     target_user,
		Password: target_password,
	}
	test_client.Log = testutil.Logger{}

	var err error
	err = test_client.Connect()
	require.NoError(t, err)
	err = test_client.Close()
	require.NoError(t, err)
}

func TestInitAndConnect(t *testing.T) {
	var test_client = &IoTDB{
		Host:     target_host,
		Port:     target_port,
		User:     target_user,
		Password: target_password,
	}
	test_client.Log = testutil.Logger{}

	var err error
	err = test_client.Init()
	require.NoError(t, err)
	err = test_client.Connect()
	require.NoError(t, err)
	err = test_client.Close()
	require.NoError(t, err)
}

func generateTestMetric(
	name string,
	tags []telegraf.Tag,
	fields []telegraf.Field,
	timestamp time.Time,
) telegraf.Metric {
	m := metric.New(name, map[string]string{}, map[string]interface{}{}, timestamp)
	for _, tag := range tags {
		m.AddTag(tag.Key, tag.Value)
	}
	for _, field := range fields {
		m.AddField(field.Key, field.Value)
	}
	return m
}

var (
	const_TestTimestamp = time.Date(2022, time.July, 20, 12, 25, 33, 44, time.UTC)
	testMetrics         = []telegraf.Metric{
		generateTestMetric(
			"root.computer.fan",
			[]telegraf.Tag{
				{Key: "price", Value: "expensive"},
				{Key: "owner", Value: "cpu"},
			},
			[]telegraf.Field{
				{Key: "temperature", Value: float64(42.55)},
				{Key: "counter", Value: int64(987654321)},
			},
			const_TestTimestamp,
		),
		generateTestMetric(
			"root.computer.fan",
			[]telegraf.Tag{
				{Key: "price", Value: "cheap"},
				{Key: "owner", Value: "gpu"},
			},
			[]telegraf.Field{
				{Key: "temperature", Value: float64(56.24)},
				{Key: "counter", Value: int64(123456789)},
			},
			const_TestTimestamp,
		),
		generateTestMetric(
			"root.computer.keyboard",
			[]telegraf.Tag{},
			[]telegraf.Field{
				{Key: "temperature", Value: float64(30.33)},
				{Key: "counter", Value: int64(123456789)},
				{Key: "unsigned", Value: uint64(math.MaxInt64 + 1000)},
				{Key: "string", Value: "Made in China."},
				{Key: "bool", Value: bool(false)},
			},
			const_TestTimestamp,
		),
	}
)

// compare two RecordsWithTags, returns True if and only if they are the same.
func compareRecords(rwt1 *RecordsWithTags, rwt2 *RecordsWithTags, Log telegraf.Logger) bool {
	if len(rwt1.DeviceId_list) == len(rwt2.DeviceId_list) &&
		len(rwt1.Measurements_list) == len(rwt2.Measurements_list) &&
		len(rwt1.Values_list) == len(rwt2.Values_list) &&
		len(rwt1.DataTypes_list) == len(rwt2.DataTypes_list) &&
		len(rwt1.Timestamp_list) == len(rwt2.Timestamp_list) {
		// ok
	} else {
		Log.Errorf("compareRecords Cechk failed. Two RecordsWithTags has different shape.")
		return false
	}
	for index, deviceID := range rwt1.DeviceId_list {
		if !(deviceID == rwt2.DeviceId_list[index]) {
			Log.Errorf("compareRecords Cechk failed. rwt1.DeviceId_list[%d]=%v, rwt2.DeviceId_list[%d]=%v.",
				index, deviceID, index, rwt2.DeviceId_list[index])
			return false
		}
	}
	for index, m_list := range rwt1.Measurements_list {
		if !(len(m_list) == len(rwt2.Measurements_list[index])) {
			Log.Errorf("compareRecords Cechk failed. Two Measurements_list has different shape. %d : %d",
				len(m_list), len(rwt2.Measurements_list[index]))
			return false
		}
		for index_2, m := range rwt1.Measurements_list[index] {
			if !(m == rwt2.Measurements_list[index][index_2]) {
				Log.Errorf("compareRecords Cechk failed. rwt1.Measurements_list[%d][%d]=%v, rwt2.Measurements_list[%d][%d]=%v.",
					index, index_2, m, index, index_2, rwt2.Measurements_list[index][index_2])
				return false
			}
		}
	}
	for index, m_list := range rwt1.Values_list {
		if !(len(m_list) == len(rwt2.Values_list[index])) {
			Log.Errorf("compareRecords Cechk failed. Two Values_list has different shape. %d : %d",
				len(m_list), len(rwt2.Values_list[index]))
			return false
		}
		for index_2, m := range rwt1.Values_list[index] {
			if !(m == rwt2.Values_list[index][index_2]) {
				Log.Errorf("compareRecords Cechk failed. rwt1.Values_list[%d][%d]=%v, rwt2.Values_list[%d][%d]=%v.",
					index, index_2, m, index, index_2, rwt2.Values_list[index][index_2])
				return false
			}
		}
	}
	for index, m_list := range rwt1.DataTypes_list {
		if !(len(m_list) == len(rwt2.DataTypes_list[index])) {
			Log.Errorf("compareRecords Cechk failed. Two DataTypes_list has different shape. %d : %d",
				len(m_list), len(rwt2.DataTypes_list[index]))
			return false
		}
		for index_2, m := range rwt1.DataTypes_list[index] {
			if !(m == rwt2.DataTypes_list[index][index_2]) {
				Log.Errorf("compareRecords Cechk failed. rwt1.DataTypes_list[%d][%d]=%v, rwt2.DataTypes_list[%d][%d]=%v.",
					index, index_2, m, index, index_2, rwt2.DataTypes_list[index][index_2])
				return false
			}
		}
	}
	for index, timestamp := range rwt1.Timestamp_list {
		if !(timestamp == rwt2.Timestamp_list[index]) {
			Log.Errorf("compareRecords Cechk failed. rwt1.DeviceId_list[%d]=%v, rwt2.DeviceId_list[%d]=%v.",
				index, timestamp, index, rwt2.Timestamp_list[index])
			return false
		}
	}
	return true
}

func TestMetricConvertion_01(t *testing.T) {
	var test_client = &IoTDB{
		Host:            target_host,
		Port:            target_port,
		User:            target_user,
		Password:        target_password,
		ConvertUint64To: "ToInt64",
		TimeStampUnit:   "nanosecond",
		TreateTagsAs:    "Measurements",
	}
	test_client.Log = testutil.Logger{}

	result, err := test_client.ConvertMetricsToRecordsWithTags(testMetrics)
	require.NoError(t, err)
	var testRecordsWithTags_01 = RecordsWithTags{
		DeviceId_list: []string{"root.computer.fan", "root.computer.fan", "root.computer.keyboard"},
		Measurements_list: [][]string{
			{"temperature", "counter"}, {"temperature", "counter"},
			{"temperature", "counter", "unsigned", "string", "bool"},
		},
		Values_list: [][]interface{}{
			{float64(42.55), int64(987654321)},
			{float64(56.24), int64(123456789)},
			{float64(30.33), int64(123456789), int64(math.MaxInt64), "Made in China.", bool(false)},
		},
		DataTypes_list: [][]client.TSDataType{
			{client.DOUBLE, client.INT64},
			{client.DOUBLE, client.INT64},
			{client.DOUBLE, client.INT64, client.INT64, client.TEXT, client.BOOLEAN},
		},
		Timestamp_list: []int64{
			const_TestTimestamp.UnixNano(), const_TestTimestamp.UnixNano(), const_TestTimestamp.UnixNano(),
		},
	}
	require.True(t, compareRecords(result, &testRecordsWithTags_01, test_client.Log))
}

func TestMetricConvertion_02(t *testing.T) {
	var test_client = &IoTDB{
		Host:            target_host,
		Port:            target_port,
		User:            target_user,
		Password:        target_password,
		ConvertUint64To: "Text",
		TimeStampUnit:   "nanosecond",
		TreateTagsAs:    "Measurements",
	}
	test_client.Log = testutil.Logger{}

	result, err := test_client.ConvertMetricsToRecordsWithTags(testMetrics)
	require.NoError(t, err)
	var testRecordsWithTags_02 = RecordsWithTags{
		DeviceId_list: []string{"root.computer.fan", "root.computer.fan", "root.computer.keyboard"},
		Measurements_list: [][]string{
			{"temperature", "counter"}, {"temperature", "counter"},
			{"temperature", "counter", "unsigned", "string", "bool"},
		},
		Values_list: [][]interface{}{
			{float64(42.55), int64(987654321)},
			{float64(56.24), int64(123456789)},
			{float64(30.33), int64(123456789), fmt.Sprintf("%d", uint64(math.MaxInt64+1000)), "Made in China.", bool(false)},
		},
		DataTypes_list: [][]client.TSDataType{
			{client.DOUBLE, client.INT64},
			{client.DOUBLE, client.INT64},
			{client.DOUBLE, client.INT64, client.TEXT, client.TEXT, client.BOOLEAN},
		},
		Timestamp_list: []int64{
			const_TestTimestamp.UnixNano(), const_TestTimestamp.UnixNano(), const_TestTimestamp.UnixNano(),
		},
	}
	require.True(t, compareRecords(result, &testRecordsWithTags_02, test_client.Log))
}

func TestMetricConvertion_03(t *testing.T) {
	var test_client = &IoTDB{
		Host:            target_host,
		Port:            target_port,
		User:            target_user,
		Password:        target_password,
		ConvertUint64To: "ToInt64",
		TimeStampUnit:   "second",
		TreateTagsAs:    "Measurements",
	}
	test_client.Log = testutil.Logger{}

	result, err := test_client.ConvertMetricsToRecordsWithTags(testMetrics)
	require.NoError(t, err)
	var testRecordsWithTags_03 = RecordsWithTags{
		DeviceId_list: []string{"root.computer.fan", "root.computer.fan", "root.computer.keyboard"},
		Measurements_list: [][]string{
			{"temperature", "counter"}, {"temperature", "counter"},
			{"temperature", "counter", "unsigned", "string", "bool"},
		},
		Values_list: [][]interface{}{
			{float64(42.55), int64(987654321)},
			{float64(56.24), int64(123456789)},
			{float64(30.33), int64(123456789), int64(math.MaxInt64), "Made in China.", bool(false)},
		},
		DataTypes_list: [][]client.TSDataType{
			{client.DOUBLE, client.INT64},
			{client.DOUBLE, client.INT64},
			{client.DOUBLE, client.INT64, client.INT64, client.TEXT, client.BOOLEAN},
		},
		Timestamp_list: []int64{
			const_TestTimestamp.Unix(), const_TestTimestamp.Unix(), const_TestTimestamp.Unix(),
		},
	}
	require.True(t, compareRecords(result, &testRecordsWithTags_03, test_client.Log))
}

func TestTagsConvertion_05(t *testing.T) {
	var test_client = &IoTDB{
		Host:            target_host,
		Port:            target_port,
		User:            target_user,
		Password:        target_password,
		ConvertUint64To: "ToInt64",
		TimeStampUnit:   "nanosecond",
		TreateTagsAs:    "Measurements",
	}
	test_client.Log = testutil.Logger{}

	result, err := test_client.ConvertMetricsToRecordsWithTags(testMetrics)
	require.NoError(t, err)
	err = test_client.ModifiyRecordsWithTags(result)
	require.NoError(t, err)
	testRecordsWithTags_05 := RecordsWithTags{
		DeviceId_list: []string{"root.computer.fan", "root.computer.fan", "root.computer.keyboard"},
		Measurements_list: [][]string{
			{"temperature", "counter", "owner", "price"}, {"temperature", "counter", "owner", "price"},
			{"temperature", "counter", "unsigned", "string", "bool"},
		},
		Values_list: [][]interface{}{
			{float64(42.55), int64(987654321), "cpu", "expensive"},
			{float64(56.24), int64(123456789), "gpu", "cheap"},
			{float64(30.33), int64(123456789), int64(math.MaxInt64), "Made in China.", bool(false)},
		},
		DataTypes_list: [][]client.TSDataType{
			{client.DOUBLE, client.INT64, client.TEXT, client.TEXT},
			{client.DOUBLE, client.INT64, client.TEXT, client.TEXT},
			{client.DOUBLE, client.INT64, client.INT64, client.TEXT, client.BOOLEAN},
		},
		Timestamp_list: []int64{
			const_TestTimestamp.UnixNano(), const_TestTimestamp.UnixNano(), const_TestTimestamp.UnixNano(),
		},
	}
	require.True(t, compareRecords(result, &testRecordsWithTags_05, test_client.Log))
}

func TestTagsConvertion_06(t *testing.T) {
	var test_client = &IoTDB{
		Host:            target_host,
		Port:            target_port,
		User:            target_user,
		Password:        target_password,
		ConvertUint64To: "ToInt64",
		TimeStampUnit:   "nanosecond",
		TreateTagsAs:    "DeviceID_subtree",
	}
	test_client.Log = testutil.Logger{}

	result, err := test_client.ConvertMetricsToRecordsWithTags(testMetrics)
	require.NoError(t, err)
	err = test_client.ModifiyRecordsWithTags(result)
	require.NoError(t, err)
	testRecordsWithTags_06 := RecordsWithTags{
		DeviceId_list: []string{"root.computer.fan.cpu.expensive", "root.computer.fan.gpu.cheap", "root.computer.keyboard"},
		Measurements_list: [][]string{
			{"temperature", "counter"}, {"temperature", "counter"},
			{"temperature", "counter", "unsigned", "string", "bool"},
		},
		Values_list: [][]interface{}{
			{float64(42.55), int64(987654321)},
			{float64(56.24), int64(123456789)},
			{float64(30.33), int64(123456789), int64(math.MaxInt64), "Made in China.", bool(false)},
		},
		DataTypes_list: [][]client.TSDataType{
			{client.DOUBLE, client.INT64},
			{client.DOUBLE, client.INT64},
			{client.DOUBLE, client.INT64, client.INT64, client.TEXT, client.BOOLEAN},
		},
		Timestamp_list: []int64{
			const_TestTimestamp.UnixNano(), const_TestTimestamp.UnixNano(), const_TestTimestamp.UnixNano(),
		},
	}
	require.True(t, compareRecords(result, &testRecordsWithTags_06, test_client.Log))
}
