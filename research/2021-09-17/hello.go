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

func supportedGroups(hello *tls.ClientHelloSpec) []tls.CurveID {
   for _, ext := range hello.Extensions {
      sc, ok := ext.(*tls.SupportedCurvesExtension)
      if ok {
         return sc.Curves
      }
   }
   return nil
}

func supportedPoints(hello *tls.ClientHelloSpec) []uint8 {
   for _, ext := range hello.Extensions {
      sp, ok := ext.(*tls.SupportedPointsExtension)
      if ok {
         return sp.SupportedPoints
      }
   }
   return nil
}

func bare(hello *tls.ClientHelloSpec) []byte {
   groups := supportedGroups(hello)
   points := supportedPoints(hello)
   // Version = uint16 => maximum = 65536 = 5chars + 1 field sep
   maxPossibleBufferLength := 5+1
   // CipherSuite = uint16 => maximum = 65536 = 5chars
   maxPossibleBufferLength += (5+1)*len(hello.CipherSuites)
   // uint16 = 2B => maximum = 65536 = 5chars
   maxPossibleBufferLength += (5+1)*len(hello.Extensions)
   // uint16 = 2B => maximum = 65536 = 5chars
   maxPossibleBufferLength += (5+1)*len(groups)
   // uint8 = 1B => maximum = 256 = 3chars
   maxPossibleBufferLength += (3+1)*len(points)
   buffer := make([]byte, 0, maxPossibleBufferLength)
   buffer = strconv.AppendInt(buffer, int64(hello.TLSVersMin), 10)
   buffer = append(buffer, sepFieldByte)
   // collect cipher suites
   last := len(hello.CipherSuites) - 1
   if len(hello.CipherSuites) > 1 {
      for _, cipher := range hello.CipherSuites[:last] {
         // filter GREASE values
         if !greaseValues[uint16(cipher)] {
            buffer = strconv.AppendInt(buffer, int64(cipher), 10)
            buffer = append(buffer, sepValueByte)
         }
      }
   }
   // append last element if cipher suites are not empty
   if last != -1 {
      cipher := hello.CipherSuites[last]
      // filter GREASE values
      if !greaseValues[uint16(cipher)] {
         buffer = strconv.AppendInt(buffer, int64(cipher), 10)
      }
   }
   buffer = append(buffer, sepFieldByte)
   // collect extensions
   last = len(hello.Extensions) - 1
   if len(hello.Extensions) > 1 {
      for _, ext := range hello.Extensions[:last] {
         // filter GREASE values
         if !greaseValues[uint16(ext)] {
            buffer = strconv.AppendInt(buffer, int64(ext), 10)
            buffer = append(buffer, sepValueByte)
         }
      }
   }
   // append last element if extensions are not empty
   if last != -1 {
      ext := hello.Extensions[last]
      // filter GREASE values
      if !greaseValues[uint16(ext)] {
         buffer = strconv.AppendInt(buffer, int64(ext), 10)
      }
   }
   buffer = append(buffer, sepFieldByte)
   // collect supported groups
   last = len(groups) - 1
   if len(groups) > 1 {
      for _, group := range groups[:last] {
         // filter GREASE values
         if !greaseValues[uint16(group)] {
            buffer = strconv.AppendInt(buffer, int64(group), 10)
            buffer = append(buffer, sepValueByte)
         }
      }
   }
   // append last element if supported groups are not empty
   if last != -1 {
      // filter GREASE values
      if !greaseValues[uint16(groups[last])] {
         buffer = strconv.AppendInt(buffer, int64(groups[last]), 10)
      }
   }
   buffer = append(buffer, sepFieldByte)
   // collect supported points
   last = len(points) - 1
   if len(points) > 1 {
      for _, point := range points[:last] {
         buffer = strconv.AppendInt(buffer, int64(point), 10)
         buffer = append(buffer, sepValueByte)
      }
   }
   // append last element if supported points are not empty
   if last != -1 {
      buffer = strconv.AppendInt(buffer, int64(points[last]), 10)
   }
   return buffer
}
