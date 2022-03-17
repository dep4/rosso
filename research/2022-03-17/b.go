package main

func main() {
   obj := object{
      "one": object{
         "two": object{"three": 3},
      },
   }
   three := newToken[int](obj).get("one").get("two").get("three")
   println(three.value == 3)
}
