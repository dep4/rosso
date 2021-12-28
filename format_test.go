package format

import (
   "fmt"
   "testing"
)

const trim = "https://east.manifest.na.theplatform.com/m/NnzsPC/T1Q_Z0YvTepY,thXzeZMxUtNl,E5dGFwKVZB7N,hB6R4we_olVJ,csINwUhvOCSf,RsLMUFAzCWSW.m3u8?sid=fce8c0b5-a196-4580-895b-d0675d794419&policy=189081367&date=1640742109851&ip=72.181.23.38&schema=1.1&manifest=M3U&tracking=true&switch=HLSServiceSecure&p2=null&am_sdkv=null&_fw_did=&nw=169843&f1=&am_extmp=default&uuid=optout-886598306373728&am_appv=null&mode=on-demand&uoo=1&sfid=9244572&s3=null&am_buildv=null&am_abvrtd=null&s4=poc&_fw_ae=&debug=false&csid=oneapp_phone_android_app_ondemand&am_cpsv=4.0.0-2&metr=1023&bundleId=&userAgent=Go-http-client%2F1.1&e1=default&prof=nbcu_android_cts_bl&afid=200265138&c3=null&a2=4.0.0-2&a3=null&a5=null&a6=0&rdid=android&am_crmid=null&am_playerv=null&did=optout&am_stitcherv=poc&am_abtestid=null&sig=cb8ba51c90f7981d8997c970c618ecdddc145072d5e23d36c8be6887ecaa36af"

func TestPercent(t *testing.T) {
   tots := []int{0, 3}
   for _, tot := range tots {
      val := PercentInt(2, tot)
      fmt.Println(val)
   }
}

func TestSymbol(t *testing.T) {
   nums := []int64{999, 1_234_567_890}
   for _, num := range nums {
      val := Number.LabelInt(num)
      fmt.Println(val)
   }
}

func TestTrim(t *testing.T) {
   val := Trim(trim)
   fmt.Println(val)
}
