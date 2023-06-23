package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func Test_run(t *testing.T) {
	t.Run("fails when no filename is specified", func(t *testing.T) {
		err := run([]string{"-f"}, io.Discard)
		assert.ErrorContains(t, err, "flag needs an argument: -f")
	})

	t.Run("fails when reading from non-existing file", func(t *testing.T) {
		err := run([]string{"-f", "nowhere"}, io.Discard)
		assert.ErrorContains(t, err, "no such file")
	})

	t.Run("fails when reading invalid cloudevent from file", func(t *testing.T) {
		file, err := os.CreateTemp("", "invalid-ce*.json")
		assert.NilError(t, err)

		_, err = file.WriteString(invalidCloudEvent)
		assert.NilError(t, err)

		t.Cleanup(func() {
			err = os.Remove(file.Name())
			assert.NilError(t, err)
		})

		err = run([]string{"-f", file.Name()}, io.Discard)
		assert.ErrorContains(t, err, "specversion: no specversion")
	})

	t.Run("succeeds when reading valid cloudevent from file", func(t *testing.T) {
		file, err := os.CreateTemp("", "valid-ce*.json")
		assert.NilError(t, err)

		_, err = file.WriteString(validCloudEventV1)
		assert.NilError(t, err)

		t.Cleanup(func() {
			err = os.Remove(file.Name())
			assert.NilError(t, err)
		})

		var stdout bytes.Buffer
		err = run([]string{"-f", file.Name()}, &stdout)
		assert.NilError(t, err)
		assert.Assert(t,
			strings.Contains(stdout.String(), "Event is a valid CloudEvent (spec version: 1.0)"),
			"got string: %s", stdout.String())
	})
}

const (
	validCloudEventV1 = `{
  "specversion": "1.0",
  "id": "8c00dd69-650c-1789-ccbb-7c3bc1ba76e0",
  "source": "some.source",
  "type": "test.event",
  "time": "2023-06-23T09:05:31Z",
  "awsregion": "eu-central-1",
  "account": "1234567890",
  "data": {
    "topic": "json-values-topic",
    "partition": 0,
    "offset": 0,
    "timestamp": 1684841916831,
    "timestampType": "CreateTime",
    "headers": [],
    "key": null,
    "value": {
      "orderItems": [
        "item-1",
        "item-2"
      ],
      "orderCreatedTime": "Tue May 23 13:38:46 CEST 2023"
    }
  }
}`

	invalidCloudEvent = `{
  "id": "8c00dd69-650c-1789-ccbb-7c3bc1ba76e0",
  "source": "some.source",
  "time": "2023-06-23T09:05:31Z",
  "data": {}
  }`
)
