package hunting

import (
	"fmt"
	"log"
	"math"

	"github.com/golang/glog"
	deviant "github.com/recluse-games/deviant-protobuf/genproto/go"
)

type gridLocation struct {
	X float64
	Y float64
}

// Vertex a vertex represents a point in our grid.
type Vertex struct {
	X int
	Y int
}

// EntityVertexPair A mapping of entities and vertexs that they exist at.
type EntityVertexPair struct {
	entity *deviant.Entity
	vertex *Vertex
}

// CardVertexPair A mapping of cards and the vertex they would be played at.
type CardVertexPair struct {
	card   *deviant.Card
	vertex *Vertex
}

// GenerateEntityLocationPairs Returns the entities of a certain alignment and their locations in points.
func GenerateEntityLocationPairs(alignment deviant.Alignment, entities []*deviant.EntitiesRow) []EntityVertexPair {
	entityVertexPairs := []EntityVertexPair{}

	for y, entityRow := range entities {
		for x, entity := range entityRow.Entities {
			entityVertexPair := EntityVertexPair{
				entity: entity,
				vertex: &Vertex{
					X: x,
					Y: y,
				},
			}

			entityVertexPairs = append(entityVertexPairs, entityVertexPair)
		}
	}

	return entityVertexPairs
}

// GetEntityVertex Get the location of yourself
func GetEntityVertex(desiredEntity *deviant.Entity, entities []*deviant.EntitiesRow) *Vertex {
	for y, entityRow := range entities {
		for x, entity := range entityRow.Entities {
			if desiredEntity.Id == entity.Id {
				entityVertexPair := &Vertex{
					X: x,
					Y: y,
				}

				return entityVertexPair
			}
		}
	}

	return nil
}

// Generate a list of all avaliable locations based on AP cost - FloodFill
func floodFill(startx int32, starty int32, x int32, y int32, filledID string, blockedID string, limit int32, tiles []*[]*deviant.Tile) {
	if (*tiles[x])[y].Id != blockedID && (*tiles[x])[y].Id != filledID {
		var apCostX int32
		var apCostY int32

		if startx > x {
			apCostX = startx - x
		} else if startx < x {
			apCostX = x - startx
		} else {
			apCostX = 0
		}

		if starty > y {
			apCostY = starty - y
		} else if starty < y {
			apCostY = y - starty
		} else {
			apCostY = 0
		}

		newTile := &deviant.Tile{}
		newTile.X = int32(x)
		newTile.Y = int32(y)
		newTile.Id = filledID
		(*tiles[x])[y] = newTile

		message := fmt.Sprintf("%v", (*tiles[x])[y].Id)
		glog.Info(message)

		if limit-apCostX-apCostY >= 0 {
			if x+1 <= 8 {
				floodFill(startx, starty, x+1, y, filledID, blockedID, limit, tiles)
			}

			if y+1 <= 7 {
				floodFill(startx, starty, x, y+1, filledID, blockedID, limit, tiles)
			}

			if x-1 >= 0 {
				floodFill(startx, starty, x-1, y, filledID, blockedID, limit, tiles)
			}

			if y-1 >= 0 {
				floodFill(startx, starty, x, y-1, filledID, blockedID, limit, tiles)
			}
		}
	}
}

// GeneratePermissableMoves Generate a list of permissable moves.
func GeneratePermissableMoves(requestedMoveAction *deviant.EntityMoveAction, avaliableAp int32, entities *deviant.Entities) []*deviant.Tile {
	finalTiles := []*deviant.Tile{}
	moveTargetTiles := []*[]*deviant.Tile{}

	for y := 0; y < len(entities.Entities); y++ {
		newRow := []*deviant.Tile{}

		for x := 0; x < len(entities.Entities[y].Entities); x++ {
			newTile := deviant.Tile{}
			newTile.X = int32(y)
			newTile.Y = int32(x)

			if entities.Entities[y].Entities[x].Id != "" {
				newTile.Id = "select_0002"
			}

			if int32(y) == requestedMoveAction.StartXPosition && int32(x) == requestedMoveAction.StartYPosition {
				newTile.Id = "select_0001"
			}

			newRow = append(newRow, &newTile)
		}

		moveTargetTiles = append(moveTargetTiles, &newRow)
	}

	floodFill(requestedMoveAction.StartXPosition, requestedMoveAction.StartYPosition, requestedMoveAction.StartXPosition, requestedMoveAction.StartYPosition, "select_0000", "select_0002", avaliableAp, moveTargetTiles)

	for _, row := range moveTargetTiles {
		for _, tile := range *row {
			if (*tile).Id == "select_0000" {
				finalTiles = append(finalTiles, tile)
			}
		}
	}

	return finalTiles
}

// IsMovePermissiable Determines if the move is permissable using a flood fill algorithm and ap cost.
func IsMovePermissiable(activeEntity *deviant.Entity, requestedMoveAction *deviant.EntityMoveAction, encounter *deviant.Encounter) bool {
	isRequestedMoveValid := false

	validTiles := GeneratePermissableMoves(requestedMoveAction, activeEntity.Ap, encounter.Board.Entities)

	for _, tile := range validTiles {
		if tile.X == requestedMoveAction.FinalXPosition && tile.Y == requestedMoveAction.FinalYPosition {
			isRequestedMoveValid = true
		}
	}

	return isRequestedMoveValid
}

// GenerateValidMoveVertexes Generate list of vertexs
func GenerateValidMoveVertexes(entity *deviant.Entity, entities []*deviant.EntitiesRow, encounter *deviant.Encounter) *Vertex {
	validMoveVertexs := []*Vertex{}
	startingLocation := GetEntityVertex(entity, entities)

	for y, entityRow := range entities {
		for x, entity := range entityRow.Entities {
			entityVertexPair := &Vertex{
				X: x,
				Y: y,
			}

			desiredMove := &deviant.EntityMoveAction{
				StartXPosition: int32(startingLocation.X),
				StartYPosition: int32(startingLocation.Y),
				FinalXPosition: int32(entityVertexPair.X),
				FinalYPosition: int32(entityVertexPair.Y),
			}

			isVertexValid := IsMovePermissiable(entity, desiredMove, encounter)

			if isVertexValid {
				validMoveVertexs = append(validMoveVertexs, entityVertexPair)
			}
		}
	}

	return nil
}

// LOOP

// Generate list of all avaliable cards in hand in each rotation at current location based on avaliable AP
func rotateTilePatterns(ocx float64, ocy float64, px float64, py float64, rotationAngle float64) *gridLocation {
	var radians = (math.Pi / 180) * rotationAngle
	var s = math.Sin(radians)
	var c = math.Cos(radians)

	// translate point back to origin:
	px -= ocx
	py -= ocy

	// rotate point
	var xnew = px*c - py*s
	var ynew = px*s + py*c

	// translate point back:
	px = xnew + ocx
	py = ynew + ocy

	return &gridLocation{px, py}
}

func convertDirectionToDegree(characterRotation deviant.EntityRotationNames) float64 {
	switch characterRotation {
	case deviant.EntityRotationNames_NORTH:
		return 180.00
	case deviant.EntityRotationNames_SOUTH:
		return 0.00
	case deviant.EntityRotationNames_EAST:
		return 270.00
	case deviant.EntityRotationNames_WEST:
		return 90.00
	}

	return 0.00
}

func generateOffSetVertex(startingLocation *Vertex, offSet *deviant.Offset) *Vertex {
	switch offSet.Direction {
	case deviant.Direction_UP:
		newVertex := &Vertex{
			X: startingLocation.X,
			Y: startingLocation.Y,
		}

		for i := 0; int32(i) < offSet.Distance; i++ {
			newVertex.X = newVertex.X + 1
		}

		return newVertex
	case deviant.Direction_DOWN:
		newVertex := &Vertex{
			X: startingLocation.X,
			Y: startingLocation.Y,
		}

		for i := 0; int32(i) < offSet.Distance; i++ {
			newVertex.X = newVertex.X - 1
		}

		return newVertex
	case deviant.Direction_LEFT:
		newVertex := &Vertex{
			X: startingLocation.X,
			Y: startingLocation.Y,
		}

		for i := 0; int32(i) < offSet.Distance; i++ {
			newVertex.Y = newVertex.Y + 1
		}

		return newVertex
	case deviant.Direction_RIGHT:
		newVertex := &Vertex{
			X: startingLocation.X,
			Y: startingLocation.Y,
		}

		for i := 0; int32(i) < offSet.Distance; i++ {
			newVertex.Y = newVertex.Y - 1
		}

		return newVertex
	}

	return nil
}

func generatePatternVertex(startingLocation *Vertex, offSet *Vertex, direction deviant.Direction, distance int32) []*Vertex {
	patternVertexes := []*Vertex{}

	switch direction {
	case deviant.Direction_UP:
		newVertex := &Vertex{}
		for i := 0; int32(i) < distance; i++ {
			newVertex = &Vertex{
				X: startingLocation.X + offSet.X,
				Y: startingLocation.Y + offSet.Y,
			}

			patternVertexes = append(patternVertexes, newVertex)
			offSet.X = offSet.X + 1
		}

		return patternVertexes
	case deviant.Direction_DOWN:
		newVertex := &Vertex{}
		for i := 0; int32(i) < distance; i++ {
			newVertex = &Vertex{
				X: startingLocation.X + offSet.X,
				Y: startingLocation.Y + offSet.Y,
			}

			patternVertexes = append(patternVertexes, newVertex)
			offSet.X = offSet.X + -1
		}

		return patternVertexes
	case deviant.Direction_LEFT:
		newVertex := &Vertex{}
		for i := 0; int32(i) < distance; i++ {
			newVertex = &Vertex{
				X: startingLocation.X + offSet.X,
				Y: startingLocation.Y + offSet.Y,
			}

			patternVertexes = append(patternVertexes, newVertex)
			offSet.Y = offSet.Y + 1
		}

		return patternVertexes
	case deviant.Direction_RIGHT:
		newVertex := &Vertex{}
		for i := 0; int32(i) < distance; i++ {
			newVertex = &Vertex{
				X: startingLocation.X + offSet.X,
				Y: startingLocation.Y + offSet.Y,
			}

			patternVertexes = append(patternVertexes, newVertex)
			offSet.Y = offSet.Y + -1
		}

		return patternVertexes
	}

	return nil
}

func GenerateCardVertexPairs(entity *deviant.Entity, entities []*deviant.EntitiesRow, rotation deviant.EntityRotationNames) []*CardVertexPair {
	nonRotatedCardVertexPairs := []*CardVertexPair{}
	rotatedCardVertexPairs := []*CardVertexPair{}
	entityLocation := GetEntityVertex(entity, entities)

	for _, card := range entity.Hand.Cards {
		for _, play := range card.Action.Pattern {
			offsetVertex := &Vertex{
				X: 0,
				Y: 0,
			}

			for _, offSet := range play.Offset {
				offsetVertex = generateOffSetVertex(offsetVertex, offSet)
			}

			playPatternVertexes := generatePatternVertex(entityLocation, offsetVertex, play.Direction, play.Distance)

			for _, vertex := range playPatternVertexes {
				cardVertexPair := &CardVertexPair{
					card:   card,
					vertex: vertex,
				}

				log.Output(0, fmt.Sprintf("%v", cardVertexPair.vertex))

				nonRotatedCardVertexPairs = append(nonRotatedCardVertexPairs, cardVertexPair)
			}
		}
	}

	for _, cardVertexPair := range nonRotatedCardVertexPairs {

		// CAUTION: HACK - This logic should be moved somewhere else to apply rotations directly to the cards themselves maybe?
		var rotationDegree = convertDirectionToDegree(rotation)
		var rotatedPlayPair = rotateTilePatterns(float64(entityLocation.X), float64(entityLocation.Y), float64(cardVertexPair.vertex.X), float64(cardVertexPair.vertex.Y), rotationDegree)

		var x = int(math.RoundToEven(rotatedPlayPair.X))
		var y = int(math.RoundToEven(rotatedPlayPair.Y))

		rotatedCardVertexPair := &CardVertexPair{
			vertex: &Vertex{
				X: x,
				Y: y,
			},
			card: cardVertexPair.card,
		}

		rotatedCardVertexPairs = append(rotatedCardVertexPairs, rotatedCardVertexPair)
	}

	return rotatedCardVertexPairs
}

// Move to next avaliable location

// Decrement AP

// END

// Filter list of all plays down to plays which hit an enemy

// Determine HP of each enemy after each hit

// Sort list by lowest enemy HP

// Determine move that preserves the most ap but is highest on the list

// Return stream of actions to process.
