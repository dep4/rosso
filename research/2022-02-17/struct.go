package googleplay

func newRequest() request {
   var r request
   r.Version = 1
   r.Checkin.Build.SdkVersion = 1
   r.DeviceConfiguration.Keyboard = 1
   r.DeviceConfiguration.GlExtension = []string{"one", "two"}
   r.DeviceConfiguration.DeviceFeature = []deviceFeature{
      {"one"}, {"two"},
   }
   return r
}

type request struct {
   Version uint64 "14"
   Checkin struct {
      Build struct {
         SdkVersion uint64 "10"
      } "1"
   } "4"
   DeviceConfiguration struct {
      Keyboard uint64 "2"
      GlExtension []string "15"
      DeviceFeature []deviceFeature "26"
   } "18"
}

type deviceFeature struct {
   Name string "1"
}
