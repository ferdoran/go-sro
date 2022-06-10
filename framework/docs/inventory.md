# Inventory & Items

Slots Ausrüstung (1-12) - gilt nur für Waffe, Schild, Klamotten und Accessories
Weitere Slots (Devil Spirit, Avatar) ???
Slots Inventar selbst (13 - X)

Packet um Items zu nutzen (0x704C - CLIENT_PLAYER_HANDLE)

Struktur
```
Opcode
Inventory slot  byte
ItemType        ushort
UniqueId        uint (nur für pet items wohl)
```

```
itemType = CashItem + (TypeID1 * 4) + (TypeID2 * 32) + (TypeID3 * 128) + (TypeID4 * 2048)
```

```
HP Recovery herb
[C -> S][704C]
0E               ................
EC 08            ................
```
```
Global chat
[C -> S][704C]
0D               ................
ED 29            .)..............
03 00            ................
68 68 68         hhh.............
```

[Ref 1](https://www.elitepvpers.com/forum/sro-coding-corner/3349888-help-send-packet.html)
[Ref 2](https://www.elitepvpers.com/forum/sro-pserver-guides-releases/4396959-0x704c-item-type-calculation-vsro.html)