# About

Validates a given *structured and JSON-encoded CloudEvent* against the latest [CloudEvent SPEC](https://cloudevents.io/) (based on the Go SDK `v2.14.0`).

# Installation

Download a released binary from the Github [releases](https://github.com/embano1/ce-validator/releases) or build via `go` (see [below](#build)).

# Usage

```bash
./ce-validator -h
Usage of ce-validator:
  -f string
        Read event from file
  -i    Read event from standard input
flag: help requested
```

## Example

Given this input event:

```json
{
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
}
```

When reading from Standard Input (`stdin`) e.g., using [`pbaste`](https://ss64.com/osx/pbpaste.html):

```bash
pbpaste | ./ce-validator -i
Event is a valid CloudEvent (spec version: 1.0)
```

# Build

> **Note**  
> Requires Go >=1.20

```bash
go install github.com/embano1/ce-validator@latest
```