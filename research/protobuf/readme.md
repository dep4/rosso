# ProtoBuf

## Why not use generic add?

If you use generic add with method interface, it would allow recurive slices. We
could solve this by also adding type interface, but its probably simpler to just
not use generics for this.
