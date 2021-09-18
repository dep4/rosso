package ja3

import (
   "github.com/refraction-networking/utls"
   "strconv"
)

const (
   sepFieldByte byte = 44
   sepValueByte byte = 45
)

var greaseValues = map[uint16]bool{
   0x0a0a: true, 0x1a1a: true, 0x2a2a: true, 0x3a3a: true, 0x4a4a: true,
   0x5a5a: true, 0x6a6a: true, 0x7a7a: true, 0x8a8a: true, 0x9a9a: true,
   0xaaaa: true, 0xbaba: true, 0xcaca: true, 0xdada: true, 0xeaea: true,
   0xfafa: true,
}

func bare(hello *tls.ClientHelloSpec) []byte {
   // Version = uint16 => maximum = 65536 = 5chars + 1 field sep
   maxPossibleBufferLength := 5+1
   // CipherSuite = uint16 => maximum = 65536 = 5chars
   maxPossibleBufferLength += (5+1)*len(hello.CipherSuites)
   // uint16 = 2B => maximum = 65536 = 5chars
   maxPossibleBufferLength += (5+1)*len(hello.Extensions)
   // uint16 = 2B => maximum = 65536 = 5chars
   maxPossibleBufferLength += (5+1)*len(hello.SupportedGroups)
   // uint8 = 1B => maximum = 256 = 3chars
   maxPossibleBufferLength += (3+1)*len(hello.SupportedPoints)
   buffer := make([]byte, 0, maxPossibleBufferLength)
   buffer = strconv.AppendInt(buffer, int64(hello.TLSVersMin), 10)
   buffer = append(buffer, sepFieldByte)
   // collect cipher suites
   lastElem := len(hello.CipherSuites) - 1
   if len(hello.CipherSuites) > 1 {
      for _, e := range hello.CipherSuites[:lastElem] {
         // filter GREASE values
         if !greaseValues[uint16(e)] {
            buffer = strconv.AppendInt(buffer, int64(e), 10)
            buffer = append(buffer, sepValueByte)
         }
      }
   }
   // append last element if cipher suites are not empty
   if lastElem != -1 {
      // filter GREASE values
      if !greaseValues[uint16(hello.CipherSuites[lastElem])] {
         buffer = strconv.AppendInt(buffer, int64(hello.CipherSuites[lastElem]), 10)
      }
   }
   buffer = append(buffer, sepFieldByte)
   // collect extensions
   lastElem = len(hello.Extensions) - 1
   if len(hello.Extensions) > 1 {
      for _, e := range hello.Extensions[:lastElem] {
         // filter GREASE values
         if !greaseValues[uint16(e)] {
            buffer = strconv.AppendInt(buffer, int64(e), 10)
            buffer = append(buffer, sepValueByte)
         }
      }
   }
   // append last element if extensions are not empty
   if lastElem != -1 {
      // filter GREASE values
      if !greaseValues[uint16(hello.Extensions[lastElem])] {
         buffer = strconv.AppendInt(buffer, int64(hello.Extensions[lastElem]), 10)
      }
   }
   buffer = append(buffer, sepFieldByte)
   // collect supported groups
   lastElem = len(hello.SupportedGroups) - 1
   if len(hello.SupportedGroups) > 1 {
      for _, e := range hello.SupportedGroups[:lastElem] {
         // filter GREASE values
         if !greaseValues[uint16(e)] {
            buffer = strconv.AppendInt(buffer, int64(e), 10)
            buffer = append(buffer, sepValueByte)
         }
      }
   }
   // append last element if supported groups are not empty
   if lastElem != -1 {
      // filter GREASE values
      if !greaseValues[uint16(hello.SupportedGroups[lastElem])] {
         buffer = strconv.AppendInt(buffer, int64(hello.SupportedGroups[lastElem]), 10)
      }
   }
   buffer = append(buffer, sepFieldByte)
   // collect supported points
   lastElem = len(hello.SupportedPoints) - 1
   if len(hello.SupportedPoints) > 1 {
      for _, e := range hello.SupportedPoints[:lastElem] {
         buffer = strconv.AppendInt(buffer, int64(e), 10)
         buffer = append(buffer, sepValueByte)
      }
   }
   // append last element if supported points are not empty
   if lastElem != -1 {
      buffer = strconv.AppendInt(buffer, int64(hello.SupportedPoints[lastElem]), 10)
   }
   return buffer
}
