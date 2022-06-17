# ProtoBuf

Where does type get checked? Start with:

~~~
DecodeTag
~~~

Then:

~~~
ConsumeTag
~~~

Then:

~~~
ConsumeField
~~~

Then:

~~~
ConsumeFieldValue
~~~

Then:

~~~
consumeFieldValueD
~~~

Then:

~~~
errCodeReserved
~~~

Then:

~~~
errReserved
errReserved    = errors.New("cannot parse reserved wire type")
~~~
