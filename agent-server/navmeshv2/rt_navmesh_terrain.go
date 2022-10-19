package navmeshv2

import (
	"fmt"
	"github.com/g3n/engine/math32"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	BlocksX     = 6
	BlocksY     = 6
	BlocksTotal = BlocksX * BlocksY

	TilesX     = 96
	TilesY     = 96
	TilesTotal = TilesX * TilesY

	VerticesX     = TilesX + 1
	VerticesY     = TilesY + 1
	VerticesTotal = VerticesX * VerticesY

	TerrainWidth     = TilesX * TileWidth
	TerrainHeight    = TilesY * TileHeight
	TerrainWidthInt  = 1920
	TerrainHeightInt = 1920
)

type RtNavmeshTerrain struct {
	RtNavmeshBase
	Region Region

	tileMap   [TilesTotal]RtNavmeshTile
	planeMap  [BlocksTotal]RtNavmeshPlane
	heightMap [VerticesTotal]float32

	Objects []RtNavmeshInstObj `json:"-"`
	Cells   []RtNavmeshCellQuad

	GlobalEdges   []RtNavmeshEdgeGlobal `json:"-"`
	InternalEdges []RtNavmeshEdgeInternal
}

func (t RtNavmeshTerrain) GetNavmeshType() RtNavmeshType {
	return RtNavmeshTypeTerrain
}

func NewRtNavmeshTerrain(filename string, region Region) RtNavmeshTerrain {
	return RtNavmeshTerrain{
		RtNavmeshBase: RtNavmeshBase{Filename: filename},
		Region:        region,
		Objects:       make([]RtNavmeshInstObj, 0),
		Cells:         make([]RtNavmeshCellQuad, 0),
		GlobalEdges:   make([]RtNavmeshEdgeGlobal, 0),
		InternalEdges: make([]RtNavmeshEdgeInternal, 0),
	}
}

func (t RtNavmeshTerrain) GetCell(index int) RtNavmeshCell {
	return &t.Cells[index]
}

func (t RtNavmeshTerrain) GetTile(x, y int) RtNavmeshTile {
	return t.tileMap[y*TilesY+x]
}

func (t RtNavmeshTerrain) GetHeight(x, y int) float32 {
	return t.heightMap[y*VerticesY+x]
}

func (t RtNavmeshTerrain) GetPlane(xBlock, zBlock int) RtNavmeshPlane {
	return t.planeMap[zBlock*BlocksY+xBlock]
}

func (t RtNavmeshTerrain) ResolveCell(pos *math32.Vector3) (RtNavmeshCellQuad, error) {
	if pos.X < 0 || pos.X > TerrainWidth || pos.Z < 0 || pos.Z > TerrainHeight {
		return RtNavmeshCellQuad{}, fmt.Errorf("position %v is not in terrain", pos)
	}
	tile := t.GetTile(int(pos.X/TileWidth), int(pos.Z/TileHeight))
	return t.Cells[tile.CellIndex], nil
}

func (t RtNavmeshTerrain) ResolveHeight(pos *math32.Vector3) float32 {
	tileX := int(pos.X / TileWidth)
	tileZ := int(pos.Z / TileHeight)

	if tileX < 0 {
		tileX = 0
	}
	if tileZ < 0 {
		tileZ = 0
	}

	tileX1 := tileX + 1
	tileZ1 := tileZ + 1
	if tileX1 >= tileX {
		tileX1 = tileX
	}
	if tileZ1 >= tileZ {
		tileZ1 = tileZ
	}

	h1 := t.GetHeight(tileX, tileZ)
	h2 := t.GetHeight(tileX, tileZ1)
	h3 := t.GetHeight(tileX1, tileZ)
	h4 := t.GetHeight(tileX1, tileZ1)

	// h1--------h3
	// |   |      |
	// |   |      |
	// h5--+------h6
	// |   |      |
	// h2--------h4

	tileOffsetX := pos.X - (TileWidth * float32(tileX))
	tileOffsetXLength := tileOffsetX / TileWidth
	tileOffsetZ := pos.Z - (TileHeight * float32(tileZ))
	tileOffsetZLength := tileOffsetZ / TileHeight

	h5 := h1 + (h2-h1)*tileOffsetZLength
	h6 := h3 + (h4-h3)*tileOffsetZLength
	yHeight := h5 + (h6-h5)*tileOffsetXLength

	return yHeight
}

func (t RtNavmeshTerrain) ToBlenderObj() []byte {
	var sb strings.Builder

	sb.WriteString("# SRO Navmesh Terrain\n")
	sb.WriteString("o " + t.Filename + "\n")
	sb.WriteString("g vertices\n")

	//for x := 0; x < VerticesX; x++ {
	//	for y := 0; y < VerticesY; y++ {
	//		x1 := float32(x) * TileWidth
	//		z1 := float32(y) * TileHeight
	//		sb.WriteString(fmt.Sprintf("v %.6f %.6f %.6f\n", x1, t.GetHeight(x, y), z1))
	//	}
	//}
	//
	//for x := 0; x < TilesX; x++ {
	//	for y := 0; y < TilesY; y++ {
	//		a := y*VerticesY + x
	//		b := a + 1
	//		c := a + VerticesX
	//		sb.WriteString(fmt.Sprintf("f %d %d %d\n", a+1, b+1, c+1))
	//		sb.WriteString(fmt.Sprintf("f %d %d %d\n", b+1, c+1, c+2))
	//	}
	//}

	vertexCounter := 1
	for _, c := range t.Cells {
		y1 := t.ResolveHeight(math32.NewVector3(c.Rect.Min.X, 0, c.Rect.Min.Y))
		y2 := t.ResolveHeight(math32.NewVector3(c.Rect.Min.X, 0, c.Rect.Max.Y))
		y3 := t.ResolveHeight(math32.NewVector3(c.Rect.Max.X, 0, c.Rect.Min.Y))
		y4 := t.ResolveHeight(math32.NewVector3(c.Rect.Max.X, 0, c.Rect.Max.Y))

		sb.WriteString(fmt.Sprintf("v %.6f %.6f %.6f\n", c.Rect.Min.X, y1, c.Rect.Min.Y))
		sb.WriteString(fmt.Sprintf("v %.6f %.6f %.6f\n", c.Rect.Min.X, y2, c.Rect.Max.Y))
		sb.WriteString(fmt.Sprintf("v %.6f %.6f %.6f\n", c.Rect.Max.X, y3, c.Rect.Min.Y))
		sb.WriteString(fmt.Sprintf("v %.6f %.6f %.6f\n", c.Rect.Max.X, y4, c.Rect.Max.Y))

		sb.WriteString(fmt.Sprintf("f %d %d %d\n", vertexCounter, vertexCounter+1, vertexCounter+2))
		sb.WriteString(fmt.Sprintf("f %d %d %d\n", vertexCounter+1, vertexCounter+2, vertexCounter+3))
		vertexCounter += 4
	}

	for _, oi := range t.Objects {
		for _, oc := range oi.Object.Cells {
			a := oc.Triangle.A.Clone().ApplyMatrix4(oi.LocalToWorld)
			b := oc.Triangle.B.Clone().ApplyMatrix4(oi.LocalToWorld)
			c := oc.Triangle.C.Clone().ApplyMatrix4(oi.LocalToWorld)

			sb.WriteString(fmt.Sprintf("v %.6f %.6f %.6f\n", a.X, a.Y, a.Z))
			sb.WriteString(fmt.Sprintf("v %.6f %.6f %.6f\n", b.X, b.Y, b.Z))
			sb.WriteString(fmt.Sprintf("v %.6f %.6f %.6f\n", c.X, c.Y, c.Z))
			sb.WriteString(fmt.Sprintf("f %d %d %d\n", vertexCounter, vertexCounter+1, vertexCounter+2))
			vertexCounter += 3
		}
	}

	return []byte(sb.String())
}

func (t *RtNavmeshTerrain) DetectUnmappedObjectEdges() {
	for _, o := range t.Objects {
		for _, ge := range o.Object.GlobalEdges {
			if !ge.Flag.IsBlocked() && !ge.Flag.IsBridge() {
				// potential linking candidate
				if ge.DstCellIndex == -1 || ge.DstCellIndex == 65535 {
					a := ge.Line.A.Clone().ApplyMatrix4(o.LocalToWorld)
					b := ge.Line.B.Clone().ApplyMatrix4(o.LocalToWorld)

					ray := math32.NewRay(a, b.Clone().SubVectors(b, a).Normalize())
					distance := a.DistanceTo(b)
					steps := int(distance / 5)
					cells := make(map[int]RtNavmeshCellQuad)
					for i := 0; i < steps; i++ {
						p := ray.At(float32(i*5), nil)
						cell, err := t.ResolveCell(p)
						if err != nil {
							logrus.Error(err)
							continue
						}

						cells[cell.Index] = cell
						cellContainsObject := false
						for _, obj := range cell.Objects {
							if obj.ID == o.ID {
								cellContainsObject = true
								break
							}
						}

						if !cellContainsObject {
							logrus.Infof("Cell %d does not contain object %d yet", cell.Index, o.ID)
							cell.Objects = append(cell.Objects, o)
						}
					}

					for _, cell := range cells {
						e := RtNavmeshEdgeGlobal{RtNavmeshEdgeBase{
							RtNavmeshEdgeMeshType: RtNavmeshEdgeMeshTypeTerrain,
							Mesh:                  t,
							Index:                 len(t.GlobalEdges) + 1,
							Line: LineSegment{
								A: a,
								B: b,
							},
							Flag:         ge.Flag,
							SrcDirection: 0,
							DstDirection: 0,
							SrcCellIndex: cell.Index,
							DstCellIndex: ge.SrcCellIndex,
							SrcCell:      &cell,
							DstCell:      ge.SrcCell,
							EventData:    ge.EventData,
						},
							int(t.Region.ID),
							o.WorldID,
						}

						t.GlobalEdges = append(t.GlobalEdges, e)
					}
				}
			}
		}
	}
}

func (t *RtNavmeshTerrain) ToJson() []byte {
	t.DetectUnmappedObjectEdges()
	var sb strings.Builder
	sb.WriteString("{") // root open
	sb.WriteString(fmt.Sprintf(`"file":"%s",`, t.Filename))
	sb.WriteString(fmt.Sprintf(`"region":{"id":%d,"x":%d,"z":%d},`, t.Region.ID, t.Region.X, t.Region.Y))
	sb.WriteString(`"heights": [`) // tiles open
	for i, h := range t.heightMap {
		sb.WriteString(fmt.Sprintf("%f", h))
		if i < len(t.heightMap)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("],")
	sb.WriteString(`"tiles": [`) // tiles open
	for i, tile := range t.tileMap {
		sb.WriteString("{") // tile open
		sb.WriteString(fmt.Sprintf(`"cellId":%d,`, tile.CellIndex))
		sb.WriteString(fmt.Sprintf(`"flag":%d`, tile.Flag))
		sb.WriteString("}") // tile close
		if i < len(t.tileMap)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("],")         // tiles close
	sb.WriteString(`"cells": [`) // cells open
	for i, c := range t.Cells {
		sb.WriteString("{") // cell open
		sb.WriteString(fmt.Sprintf(`"id":%d,`, c.Index))
		y1 := t.ResolveHeight(math32.NewVector3(c.Rect.Min.X, 0, c.Rect.Min.Y))
		y2 := t.ResolveHeight(math32.NewVector3(c.Rect.Max.X, 0, c.Rect.Max.Y))
		//y2 := t.ResolveHeight(math32.NewVector3(c.Rect.Min.X, 0, c.Rect.Max.Y))
		//y3 := t.ResolveHeight(math32.NewVector3(c.Rect.Max.X, 0, c.Rect.Min.Y))

		sb.WriteString(fmt.Sprintf(`"min":{"x":%.6f,"y":%.6f,"z":%.6f},`, c.Rect.Min.X, y1, c.Rect.Min.Y))
		sb.WriteString(fmt.Sprintf(`"max":{"x":%.6f,"y":%.6f,"z":%.6f},`, c.Rect.Max.X, y2, c.Rect.Max.Y))
		sb.WriteString(fmt.Sprintf(`"objects": [`)) // objects open
		for j, o := range c.Objects {
			sb.WriteString(fmt.Sprintf("%d", o.ID))
			if j < len(c.Objects)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString("]") // objects close
		sb.WriteString("}") // cell close
		if i < len(t.Cells)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("],") // cells close

	sb.WriteString(`"objects":[`) // objects open
	for i, oi := range t.Objects {

		sb.WriteString("{") // object open
		sb.WriteString(fmt.Sprintf(`"id":%d,`, oi.ID))

		sb.WriteString(fmt.Sprintf(`"position":{"x":%.6f,"y":%.6f,"z":%.6f},`, oi.Position.X, oi.Position.Y, oi.Position.Z))
		sb.WriteString(fmt.Sprintf(`"rotation":{"x":%.6f,"y":%.6f,"z":%.6f,"w":%.6f},`, oi.Rotation.X, oi.Rotation.Y, oi.Rotation.Z, oi.Rotation.W))

		sb.WriteString(`"cells":[`) // cells open
		for j, oc := range oi.Object.Cells {
			sb.WriteString("{") // cell open
			sb.WriteString(fmt.Sprintf(`"id":%d,`, oc.Index))
			sb.WriteString(fmt.Sprintf(`"worldId":%d,`, oi.WorldID))
			sb.WriteString(fmt.Sprintf(`"flag":%d,`, oc.Flag))

			a := oc.Triangle.A.Clone().ApplyMatrix4(oi.LocalToWorld)
			b := oc.Triangle.B.Clone().ApplyMatrix4(oi.LocalToWorld)
			c := oc.Triangle.C.Clone().ApplyMatrix4(oi.LocalToWorld)

			sb.WriteString(fmt.Sprintf(`"a":{"x":%.6f,"y":%.6f,"z":%.6f},`, a.X, a.Y, a.Z))
			sb.WriteString(fmt.Sprintf(`"b":{"x":%.6f,"y":%.6f,"z":%.6f},`, b.X, b.Y, b.Z))
			sb.WriteString(fmt.Sprintf(`"c":{"x":%.6f,"y":%.6f,"z":%.6f}`, c.X, c.Y, c.Z))

			sb.WriteString("}") // cell close
			if j < len(oi.Object.Cells)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString("],")                // cells close
		sb.WriteString(`"internalEdges":[`) // internal edges open
		for j, e := range oi.Object.InternalEdges {
			a := e.Line.A.Clone().ApplyMatrix4(oi.LocalToWorld)
			b := e.Line.B.Clone().ApplyMatrix4(oi.LocalToWorld)
			sb.WriteString("{") // edge open
			sb.WriteString(fmt.Sprintf(`"id":%d,`, e.Index))
			sb.WriteString(fmt.Sprintf(`"flag":%d,`, e.Flag))
			sb.WriteString(fmt.Sprintf(`"eventData":%d,`, e.EventData))
			sb.WriteString(fmt.Sprintf(`"srcCellId":%d,`, e.SrcCellIndex))
			sb.WriteString(fmt.Sprintf(`"dstCellId":%d,`, e.DstCellIndex))
			sb.WriteString(fmt.Sprintf(`"srcDir":%d,`, e.SrcDirection))
			sb.WriteString(fmt.Sprintf(`"dstDir":%d,`, e.DstDirection))
			sb.WriteString(fmt.Sprintf(`"a":{"x":%.6f,"y":%.6f,"z":%.6f},`, a.X, a.Y, a.Z))
			sb.WriteString(fmt.Sprintf(`"b":{"x":%.6f,"y":%.6f,"z":%.6f}`, b.X, b.Y, b.Z))
			sb.WriteString("}")
			if j < len(oi.Object.InternalEdges)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString("],")              // internal edges close
		sb.WriteString(`"globalEdges":[`) // global edges open
		for j, e := range oi.Object.GlobalEdges {
			a := e.Line.A.Clone().ApplyMatrix4(oi.LocalToWorld)
			b := e.Line.B.Clone().ApplyMatrix4(oi.LocalToWorld)
			sb.WriteString("{") // edge open
			sb.WriteString(fmt.Sprintf(`"id":%d,`, e.Index))
			sb.WriteString(fmt.Sprintf(`"flag":%d,`, e.Flag))
			sb.WriteString(fmt.Sprintf(`"eventData":%d,`, e.EventData))
			sb.WriteString(fmt.Sprintf(`"srcCellId":%d,`, e.SrcCellIndex))
			sb.WriteString(fmt.Sprintf(`"dstCellId":%d,`, e.DstCellIndex))
			sb.WriteString(fmt.Sprintf(`"srcDir":%d,`, e.SrcDirection))
			sb.WriteString(fmt.Sprintf(`"dstDir":%d,`, e.DstDirection))
			sb.WriteString(fmt.Sprintf(`"srcMeshId":%d,`, e.SrcMeshIndex))
			sb.WriteString(fmt.Sprintf(`"dstMeshId":%d,`, e.DstMeshIndex))
			sb.WriteString(fmt.Sprintf(`"a":{"x":%.6f,"y":%.6f,"z":%.6f},`, a.X, a.Y, a.Z))
			sb.WriteString(fmt.Sprintf(`"b":{"x":%.6f,"y":%.6f,"z":%.6f}`, b.X, b.Y, b.Z))
			sb.WriteString("}")
			if j < len(oi.Object.GlobalEdges)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString("],")            // global edges close
		sb.WriteString(`"edgeLinks":[`) // edge links open
		for j, e := range oi.EdgeLinks {
			sb.WriteString("{") // edge link open
			sb.WriteString(fmt.Sprintf(`"objId":%d,`, e.LinkedObjID))
			sb.WriteString(fmt.Sprintf(`"objEdgeId":%d,`, e.LinkedObjEdgeID))
			sb.WriteString(fmt.Sprintf(`"edgeId":%d`, e.EdgeID))
			sb.WriteString("}") // edge link close
			if j < len(oi.EdgeLinks)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString("]") // edge links close
		sb.WriteString("}") // object close
		if i < len(t.Objects)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("],")                // objects close
	sb.WriteString(`"internalEdges":[`) // internal edges open
	for i, e := range t.InternalEdges {
		a := e.Line.A
		b := e.Line.B
		sb.WriteString("{") // edge open
		sb.WriteString(fmt.Sprintf(`"id":%d,`, e.Index))
		sb.WriteString(fmt.Sprintf(`"flag":%d,`, e.Flag))
		sb.WriteString(fmt.Sprintf(`"srcCellId":%d,`, e.SrcCellIndex))
		sb.WriteString(fmt.Sprintf(`"dstCellId":%d,`, e.DstCellIndex))
		sb.WriteString(fmt.Sprintf(`"srcDir":%d,`, e.SrcDirection))
		sb.WriteString(fmt.Sprintf(`"dstDir":%d,`, e.DstDirection))
		sb.WriteString(fmt.Sprintf(`"a":{"x":%.6f,"y":%.6f,"z":%.6f},`, a.X, a.Y, a.Z))
		sb.WriteString(fmt.Sprintf(`"b":{"x":%.6f,"y":%.6f,"z":%.6f}`, b.X, b.Y, b.Z))
		sb.WriteString("}")
		if i < len(t.InternalEdges)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("],")              // internal edges close
	sb.WriteString(`"globalEdges":[`) // global edges open
	for i, e := range t.GlobalEdges {
		a := e.Line.A
		b := e.Line.B
		sb.WriteString("{") // edge open
		sb.WriteString(fmt.Sprintf(`"id":%d,`, e.Index))
		sb.WriteString(fmt.Sprintf(`"flag":%d,`, e.Flag))
		sb.WriteString(fmt.Sprintf(`"srcCellId":%d,`, e.SrcCellIndex))
		sb.WriteString(fmt.Sprintf(`"dstCellId":%d,`, e.DstCellIndex))
		sb.WriteString(fmt.Sprintf(`"srcDir":%d,`, e.SrcDirection))
		sb.WriteString(fmt.Sprintf(`"dstDir":%d,`, e.DstDirection))
		sb.WriteString(fmt.Sprintf(`"srcMeshId":%d,`, e.SrcMeshIndex))
		sb.WriteString(fmt.Sprintf(`"dstMeshId":%d,`, e.DstMeshIndex))
		sb.WriteString(fmt.Sprintf(`"a":{"x":%.6f,"y":%.6f,"z":%.6f},`, a.X, a.Y, a.Z))
		sb.WriteString(fmt.Sprintf(`"b":{"x":%.6f,"y":%.6f,"z":%.6f}`, b.X, b.Y, b.Z))
		sb.WriteString("}")
		if i < len(t.GlobalEdges)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]") // global edges close
	sb.WriteString("}") // root close
	return []byte(sb.String())
}
