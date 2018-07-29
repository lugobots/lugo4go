package tatics

// PlayerRegion define a virtual rectangle  which starts at point A and is finished at point B
// These points are based in the home team field, so they must be converted to the away team if
// the team is player as the visitor. Don't convert this coords as in a mirror. They must
// be rotated 180 degress.
//type PlayerRegion struct {
//	CornerA Physics.Point
//	CornerB Physics.Point
//}

//func (p *PlayerRegion) InitialPosition() Physics.Point {
//	return p.CentralDefense()
//}
//
//func (p *PlayerRegion) CentralDefense() Physics.Point {
//	return Physics.Point{
//		PosX: p.CornerA.PosX + ((p.CornerB.PosX - p.CornerA.PosX) / 5),
//		PosY: p.CornerA.PosY + (p.CornerB.PosY - p.CornerA.PosY) / 2,
//	}
//}

//var regionWidth = int(math.Round(Units.CourtWidth * 0.6))
//var regionLength = int(math.Round(Units.CourtHeight * 0.3))
//var regionOverlap = ((regionLength * 4) - Units.CourtHeight) / 3

//var HomePlayersRegions = map[BasicTypes.PlayerNumber]PlayerRegion{
//	"2": {
//		CornerA: Physics.Point{0, 0},
//		CornerB: Physics.Point{regionWidth, regionLength},
//	},
//	"3": {
//		CornerA: Physics.Point{0, regionLength - regionOverlap},
//		CornerB: Physics.Point{regionWidth, (2 * regionLength) - regionOverlap},
//	},
//	"4": {
//		CornerA: Physics.Point{0, 2 * (regionLength - regionOverlap)},
//		CornerB: Physics.Point{regionWidth, 3*regionLength - 2*regionOverlap},
//	},
//	"5": {
//		CornerA: Physics.Point{0, 3 * (regionLength - regionOverlap)},
//		CornerB: Physics.Point{regionWidth, Units.CourtHeight},
//	},
//
//
//	"6": {
//		CornerA: Physics.Point{Units.CourtHeight / 3, 0},
//		CornerB: Physics.Point{regionWidth + (Units.CourtHeight / 3), regionLength},
//	},
//	"7": {
//		CornerA: Physics.Point{Units.CourtHeight / 3, regionLength - regionOverlap},
//		CornerB: Physics.Point{regionWidth + (Units.CourtHeight / 3), (2 * regionLength) - regionOverlap},
//	},
//	"8": {
//		CornerA: Physics.Point{Units.CourtHeight / 3, 2 * (regionLength - regionOverlap)},
//		CornerB: Physics.Point{regionWidth + (Units.CourtHeight / 3), 3*regionLength - 2*regionOverlap},
//	},
//	"9": {
//		CornerA: Physics.Point{Units.CourtHeight / 3, 3 * (regionLength - regionOverlap)},
//		CornerB: Physics.Point{regionWidth + (Units.CourtHeight / 3), Units.CourtHeight},
//	},
//
//
//	"10": {
//		CornerA: Physics.Point{Units.CourtWidth / 3, regionLength - regionOverlap},
//		CornerB: Physics.Point{regionWidth + (Units.CourtWidth / 3), (2 * regionLength) - regionOverlap},
//	},
//	"11": {
//		CornerA: Physics.Point{Units.CourtWidth / 3, 2 * (regionLength) - regionOverlap},
//		CornerB: Physics.Point{regionWidth + (Units.CourtWidth / 3), 3*regionLength - 2*regionOverlap},
//	},
//}

//type PlayerPosition struct {
//	InitialPosition Physics.Point
//Region PlayerRegion
//Number BasicTypes.PlayerNumber
//}
//
