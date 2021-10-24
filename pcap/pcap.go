package pcap

import (
   "encoding/hex"
   "github.com/google/gopacket"
   "github.com/google/gopacket/layers"
   "github.com/google/gopacket/pcapgo"
   "io"
)

type Handshake []byte

func Handshakes(r io.Reader) ([]Handshake, error) {
   read, err := pcapgo.NewReader(r)
   if err != nil {
      return nil, err
   }
   var hands []Handshake
   for {
      data, _, err := read.ReadPacketData()
      if err == io.EOF {
         return hands, nil
      } else if err != nil {
         return nil, err
      }
      pack := gopacket.NewPacket(
         data, read.LinkType(), gopacket.DecodeStreamsAsDatagrams,
      )
      tls, ok := pack.Layer(layers.LayerTypeTLS).(*layers.TLS)
      if ok && tls.Handshake != nil {
         hands = append(hands, tls.BaseLayer.Contents)
      }
   }
}

func (h Handshake) String() string {
   return hex.EncodeToString(h)
}
