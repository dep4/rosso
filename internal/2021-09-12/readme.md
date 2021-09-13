# September 12 2021

After more looking, I found this entry with an identical User-Agent:

~~~json
{
   "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36",
   "Count": 21350,
   "md5": "46346146986ace52994f40a34f06b1ce",
   "Last_seen": "2021-02-26 16:43:43"
}
~~~

and the cross reference:

~~~json
{
   "md5": "46346146986ace52994f40a34f06b1ce",
   "ja3": "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53-255,0-11-10-35-5-16-23-13-43-45-51,23-24,0-1-2",
   "First_reported": "2020-05-25 15:46:11"
}
~~~

in this case, the cross reference has all fields listed, as expected. I am not
sure what is causing the missing fields in some cases.
