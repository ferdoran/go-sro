# Diskussion 1

**Datum: 22.09.2019**

## Begriffe

| Begriff              | Beschreibung | Synonyme |
| -------------------- | ------------ |--------- |
| Division             | Beschreibt eine Erdregion, wie z.B. Japan, Vietnam, Korea. Eine Division kann dabei mehrere Farms besitzen, bspw. Nord Vietnam, S�den Vietnam | - |
| Shard                | Eine Spielwelt oder ein Server in der Serverliste                                                                    |                          - |
| Zone                 | Beschreibt hier eine selbst verwaltete Region auf der Map. Z.B. Jangan, Hotan (+ Takla, Karakoram), etc.             |                          - |
| Game Server          | Hosted eine _Zone_. Pro _Shard_ existieren mehrere _Game Server_                                                     | Zone Server, Region Server |
| Shard Server         | Verwaltet den Zustand einer _Shard_ / _World_                                                                        | World Server               |
| Global Manager       | Verwaltet eine Division. Der SMC verbindet sich hiermit                                                              | Division Server            |
| Farm Manager         | Gruppiert alle _Shards_ die in seinem Server Zentrum stehen.                                                         |                          - |
| Certification Server | Verifiziert Lizenzen fuer Divisions                                                                                   |                          - |

## Beschluesse

- Farm Manager wird nicht benoetigt
- Global Manager wird nicht benoetigt
- Certification Server wird nicht benoetigt
- Billing Server wird nicht benoetigt
- Beim Transfer von Gateway Server auf Agent Server wird ein Token ausgestellt, mit dem der Client sich beim Agent Server authentifizieren kann

## Ideen

Moegliche Sachen die man aus dem GameServer rausnehmen k�nnte:

- Community (Friendlist, Memo)
- Akademie groesstenteils (nicht betroffen Map Position von Akademie Membern)
- Inventory? (Exchange mit anderen Spielern).
- Chat?
- Consignment (bis auf Interaction mit NPC)
- CTF? -> Laeuft auf einem Zone Server
- Guide?
- Guild + Union
- Quest?
- Job?
- Party? (Party Matching, verschiedene Zones, ...)
- FW? -> Laeuft auf einem Zone Server
- Silk / Billing / Gatcha (Magic Pop)
- Skill (Withdraw, Learn, + Masteries)
- Stall (Inventory, Chat, ...)?
- TAP?
- Teleport (Teleport Req NPC, ...)
- Timed Job?

## Anderes

- NavMeshes aus dem Client laden. Diese muessen noch berechnet werden.
- Farm Manager hat eine DB Verbindung. KA warum
- Invites werden auf dem jeweiligen Server gemanaged?
