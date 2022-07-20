//go:generate ../../../tools/readme_config_includer/generator
package iotdb

import (
	"testing"
	// "time"

	"github.com/stretchr/testify/require"

	// "github.com/influxdata/telegraf"
	// "github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/testutil"
)

func TestConnectAndClose(t *testing.T) {
	test_client := &IoTDB{
		Host:     "192.168.134.130",
		Port:     "6667",
		User:     "root",
		Password: "root",
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
		Host:     "192.168.134.130",
		Port:     "6667",
		User:     "root",
		Password: "root",
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
