package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/mobiledgex/edge-cloud-infra/shepherd/shepherd_common"
	"github.com/mobiledgex/edge-cloud-infra/shepherd/shepherd_platform/shepherd_unittest"
	"github.com/mobiledgex/edge-cloud/edgeproto"
	"github.com/mobiledgex/edge-cloud/log"
	"github.com/stretchr/testify/assert"
)

// Example output of resource-tracker
var testCloudletData = shepherd_common.CloudletMetrics{
	VCpuMax:  10,
	VCpuUsed: 5,
	MemMax:   500000,
	MemUsed:  1234,
	DiskMax:  400000,
	DiskUsed: 10000,
	NetRecv:  123456,
	NetSent:  654321,
}

func TestCloudletStats(t *testing.T) {
	var err error
	log.InitTracer("")
	defer log.FinishTracer()
	ctx := log.StartTestSpan(context.Background())

	testOperatorKey := edgeproto.OperatorKey{Name: "testoper"}
	cloudletKey = edgeproto.CloudletKey{
		OperatorKey: testOperatorKey,
		Name:        "testcloudlet",
	}

	testCloudletData.ComputeTS, err = types.TimestampProto(time.Now())
	assert.Nil(t, err, "Couldn't get current timestamp")
	testCloudletData.NetworkTS = testCloudletData.ComputeTS
	buf, err := json.Marshal(testCloudletData)
	assert.Nil(t, err, "marshal cloudlet metrics")
	pf = &shepherd_unittest.Platform{
		CloudletMetrics: string(buf),
	}

	cloudletStats, err := pf.GetPlatformStats(ctx)
	assert.Nil(t, err, "Get cloudlet stats")
	metrics := MarshalCloudletMetrics(&cloudletStats)
	// Should be two measurements
	assert.Equal(t, 2, len(metrics))
	// Verify the names
	assert.Equal(t, "cloudlet-utilization", metrics[0].Name)
	assert.Equal(t, "cloudlet-network", metrics[1].Name)
	// Verify metric tags
	for _, m := range metrics {
		for _, v := range m.Tags {
			if v.Name == "operator" {
				assert.Equal(t, cloudletKey.OperatorKey.Name, v.Val)
			}
			if v.Name == "cloudlet" {
				assert.Equal(t, cloudletKey.Name, v.Val)
			}
		}
	}
	// Verify metric values
	for _, v := range metrics[0].Vals {
		if v.Name == "vCpuUsed" {
			assert.Equal(t, testCloudletData.VCpuUsed, v.GetIval())
		} else if v.Name == "vCpuMax" {
			assert.Equal(t, testCloudletData.VCpuMax, v.GetIval())
		} else if v.Name == "memUsed" {
			assert.Equal(t, testCloudletData.MemUsed, v.GetIval())
		} else if v.Name == "memMax" {
			assert.Equal(t, testCloudletData.MemMax, v.GetIval())
		} else if v.Name == "diskUsed" {
			assert.Equal(t, testCloudletData.DiskUsed, v.GetIval())
		} else if v.Name == "diskMax" {
			assert.Equal(t, testCloudletData.DiskMax, v.GetIval())
		} else {
			errstr := fmt.Sprintf("Unexpected value in a metric(%v) - %s", v, v.Name)
			assert.FailNow(t, errstr)
		}
	}
	for _, v := range metrics[1].Vals {
		if v.Name == "netSent" {
			assert.Equal(t, testCloudletData.NetSent, v.GetIval())
		} else if v.Name == "netRecv" {
			assert.Equal(t, testCloudletData.NetRecv, v.GetIval())
		} else {
			errstr := fmt.Sprintf("Unexpected value in a metric(%v) - %s", v, v.Name)
			assert.FailNow(t, errstr)
		}
	}
}