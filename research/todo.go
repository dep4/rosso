package strconv

var Percent = Unit{100, "%"}

var Cardinals = []Unit{
   {1e-3, " thousand"},
   {1e-6, " million"},
   {1e-9, " billion"},
   {1e-12, " trillion"},
}

var Rates = []Unit{
   {1e-3, " kilobyte/s"},
   {1e-6, " megabyte/s"},
   {1e-9, " gigabyte/s"},
   {1e-12, " terabyte/s"},
}

var Sizes = []Unit{
   {1e-3, " kilobyte"},
   {1e-6, " megabyte"},
   {1e-9, " gigabyte"},
   {1e-12, " terabyte"},
}
