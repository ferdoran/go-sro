select r.wRegionID, r.ContinentName, n.fLocalPosX, n.fLocalPosY, n.fLocalPosZ, n.wInitialDir, n.nRadius, n.nGenerateRadius, t.ObjID, c.CodeName
from SRO_SHARD.REGION_REFERENCE r, SRO_SHARD.SPAWN_REF_NESTS n, SRO_SHARD.SPAWN_REF_TACTICS t, SRO_SHARD.CHAR_REF_DATA c
WHERE r.ContinentName = 'Jangan' 
AND r.wRegionID = n.nRegionDBID
AND n.dwTacticsID = t.TacticsID
AND t.ObjID = c.RefObjIDSPAWN_REF_NESTS