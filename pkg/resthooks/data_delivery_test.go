// SPDX-License-Identifier: Apache-2.0
package resthooks

import (
	"encoding/json"
	"testing"

	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/rmr"
	"github.com/stretchr/testify/assert"
)

type dataDeliverySender struct {
	message string
	mtype   int
	subid   int
	calls   int
}

func (s *dataDeliverySender) RmrSendToXapp(message string, mtype, subid int) bool {
	s.message, s.mtype, s.subid = message, mtype, subid
	s.calls++
	return true
}

func TestDataDeliveryMalformedBodies(t *testing.T) {
	tests := []struct {
		name string
		body interface{}
	}{
		{"nil", nil},
		{"typed nil object", map[string]interface{}(nil)},
		{"array", []interface{}{}},
		{"string", "data"},
		{"number", 1.0},
		{"empty object", map[string]interface{}{}},
		{"missing job", map[string]interface{}{"payload": "data"}},
		{"null job", map[string]interface{}{"job": nil, "payload": "data"}},
		{"numeric job", map[string]interface{}{"job": 1, "payload": "data"}},
		{"empty job", map[string]interface{}{"job": "", "payload": "data"}},
		{"missing payload", map[string]interface{}{"job": "job-1"}},
		{"null payload", map[string]interface{}{"job": "job-1", "payload": nil}},
		{"object payload", map[string]interface{}{"job": "job-1", "payload": map[string]interface{}{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender := new(dataDeliverySender)
			hook := &Resthook{iRmrSenderInst: sender}
			var err error
			assert.NotPanics(t, func() { err = hook.DataDelivery(tt.body) })
			assert.Error(t, err)
			assert.Zero(t, sender.calls)
		})
	}
}

func TestDataDeliveryValidBodies(t *testing.T) {
	for _, payload := range []string{"", "measurements"} {
		t.Run("payload="+payload, func(t *testing.T) {
			sender := new(dataDeliverySender)
			hook := &Resthook{iRmrSenderInst: sender}
			assert.NoError(t, hook.DataDelivery(map[string]interface{}{"job": "job-1", "payload": payload}))
			assert.Equal(t, 1, sender.calls)
			assert.Equal(t, a1EIDataDelivery, sender.mtype)
			assert.Equal(t, rmr.DefaultSubId, sender.subid)
			var body map[string]string
			assert.NoError(t, json.Unmarshal([]byte(sender.message), &body))
			assert.Equal(t, map[string]string{"ei_job_id": "job-1", "payload": payload}, body)
		})
	}
}
