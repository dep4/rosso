package googleplay

import (
   "strconv"
)

type message map[tag]interface{}

type tag struct {
   Number int64
   Type int64
}

func (t tag) MarshalText() ([]byte, error) {
   buf := []byte("Number:")
   buf = strconv.AppendInt(buf, t.Number, 10)
   buf = append(buf, " Type:"...)
   return strconv.AppendInt(buf, t.Type, 10), nil
}

func (m message) get(k int64) message {
   return m[tag{k, 0}].(message)
}

func (m message) set(k int64, v message) {
   m[tag{k, 0}] = v
}

func (m message) setUint64(k int64, v uint64) {
   m[tag{k, 1}] = v
}

// Normally we convert a ProtoBuf response into a struct.
// In this case, we convert a struct to a ProtoBuf request.
func (r request) message() message {
   m := make(message)
   m.set(4, make(message))
   m.get(4).set(1, make(message))
   m.set(18, make(message))
   m.setUint64(14, r.Version)
   m.get(4).get(1).setUint64(10, r.Checkin.Build.SdkVersion)
   m.get(18).setUint64(2, r.DeviceConfiguration.Keyboard)
   return m
}
