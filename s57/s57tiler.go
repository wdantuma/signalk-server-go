package s57

// Convert S-57 format to MVT vector tiles
// see MVT spec at https://github.com/mapbox/vector-tile-spec/tree/master/2.1

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lukeroth/gdal"
	"github.com/wdantuma/signalk-server-go/resources/charts"
	"github.com/wdantuma/signalk-server-go/s57/dataset"
	m "github.com/wdantuma/signalk-server-go/s57/mercantile"
	"github.com/wdantuma/signalk-server-go/s57/vectortile"
	"google.golang.org/protobuf/proto"
)

const (
	TILE_EXTENT = 4096
)

type s57Tiler struct {
	transform gdal.CoordinateTransform
	datasets  []dataset.Dataset
	valuesMap map[string]uint32
	values    []string
	keysMap   map[string]uint32
	keys      []string
}

func NewS57Tiler(datasets []dataset.Dataset) *s57Tiler {
	src := gdal.CreateSpatialReference("")
	src.FromEPSG(4326)
	dst := gdal.CreateSpatialReference("")
	dst.FromEPSG(3857)

	driver, err := gdal.GetDriverByName("S57")
	if err != nil {
		fmt.Println(err.Error())
	}
	driver.Register()

	return &s57Tiler{transform: gdal.CreateCoordinateTransform(src, dst), datasets: datasets}
}

func (s *s57Tiler) startLayer() {
	s.valuesMap = make(map[string]uint32)
	s.values = make([]string, 0)
	s.keysMap = make(map[string]uint32)
	s.keys = make([]string, 0)
}

func (s *s57Tiler) to3857(x float64, y float64) (float64, float64) {
	xs := make([]float64, 1)
	xs[0] = y
	ys := make([]float64, 1)
	ys[0] = x
	zs := make([]float64, 1)
	zs[0] = 0

	s.transform.Transform(1, xs, ys, zs)

	return xs[0], ys[0]
}

func (s *s57Tiler) toTileCoordinate(tileBounds m.Extrema, x float64, y float64, z float64) (int32, int32, int32) {
	tx, ty := s.to3857(x, y)

	ulx, uly := s.to3857(tileBounds.W, tileBounds.N)
	lrx, lry := s.to3857(tileBounds.E, tileBounds.S)

	xf := TILE_EXTENT / (lrx - ulx)
	yf := TILE_EXTENT / (uly - lry)
	xx := (tx - ulx) * xf
	yy := (uly - ty) * yf
	return int32(xx), int32(yy), 0
}

func getCommand(command int, count int) uint32 {
	cmd := (command & 0x7) | (count << 3)
	return uint32(cmd)
}

func getCoordinate(coordinate int32) uint32 {
	return uint32((coordinate << 1) ^ (coordinate >> 31))
}

func (s *s57Tiler) toMvtLinestringGeometry(geometry *gdal.Geometry, tileBounds m.Extrema) []uint32 {
	mvtGeometry := make([]uint32, 0)
	count := geometry.PointCount()
	if count > 1 {
		// moveto
		mvtGeometry = append(mvtGeometry, getCommand(1, 1))
		x, y, _ := geometry.Point(0)
		lastx, lasty, _ := s.toTileCoordinate(tileBounds, x, y, 0)
		mvtGeometry = append(mvtGeometry, getCoordinate(lastx))
		mvtGeometry = append(mvtGeometry, getCoordinate(lasty))
		// lineto
		mvtGeometry = append(mvtGeometry, getCommand(2, geometry.PointCount()-1))
		for i := 1; i < count; i++ {

			x, y, _ := geometry.Point(i)
			xx, yy, _ := s.toTileCoordinate(tileBounds, x, y, 0)
			dx := xx - lastx
			dy := yy - lasty
			mvtGeometry = append(mvtGeometry, getCoordinate(dx))
			mvtGeometry = append(mvtGeometry, getCoordinate(dy))
			lastx = xx
			lasty = yy
		}
	}

	return mvtGeometry
}

func (s *s57Tiler) toMvtPolygonGeometry(geometry *gdal.Geometry, tileBounds m.Extrema) []uint32 {
	mvtGeometry := s.toMvtLinestringGeometry(geometry, tileBounds)
	// close path
	mvtGeometry = append(mvtGeometry, getCommand(7, 1))
	return mvtGeometry
}

func (s *s57Tiler) toMvtPointGeometry(geometry *gdal.Geometry, tileBounds m.Extrema) []uint32 {
	mvtGeometry := make([]uint32, 0)
	count := geometry.PointCount()
	mvtGeometry = append(mvtGeometry, getCommand(1, count))
	for i := 0; i < count; i++ {
		x, y, _ := geometry.Point(i)
		xx, yy, _ := s.toTileCoordinate(tileBounds, x, y, 0)
		mvtGeometry = append(mvtGeometry, getCoordinate(xx))
		mvtGeometry = append(mvtGeometry, getCoordinate(yy))
	}

	return mvtGeometry
}

func (s *s57Tiler) toMvtGeometry(featureType vectortile.Tile_GeomType, geometry *gdal.Geometry, tileBounds m.Extrema) []uint32 {
	mvtGeometry := make([]uint32, 0)
	geomcount := geometry.GeometryCount()
	pointCount := geometry.PointCount()
	if geomcount > 0 {
		for i := 0; i < geomcount; i++ {
			geom := geometry.Geometry(i)
			switch featureType {
			case vectortile.Tile_POINT:
				mvtGeometry = append(mvtGeometry, s.toMvtPointGeometry(&geom, tileBounds)...)
			case vectortile.Tile_LINESTRING:
				mvtGeometry = append(mvtGeometry, s.toMvtLinestringGeometry(&geom, tileBounds)...)
			case vectortile.Tile_POLYGON:
				mvtGeometry = append(mvtGeometry, s.toMvtPolygonGeometry(&geom, tileBounds)...)
			}
		}
	} else if pointCount > 0 {
		switch featureType {
		case vectortile.Tile_POINT:
			mvtGeometry = append(mvtGeometry, s.toMvtPointGeometry(geometry, tileBounds)...)
		case vectortile.Tile_LINESTRING:
			mvtGeometry = append(mvtGeometry, s.toMvtLinestringGeometry(geometry, tileBounds)...)
		case vectortile.Tile_POLYGON:
			mvtGeometry = append(mvtGeometry, s.toMvtPolygonGeometry(geometry, tileBounds)...)
		}
	}

	return mvtGeometry
}

func (s *s57Tiler) getMvtFeatureType(geometry *gdal.Geometry) *vectortile.Tile_GeomType {
	geomType := geometry.Type()
	var mvtGeomType vectortile.Tile_GeomType
	switch geomType {
	case gdal.GT_LineString, gdal.GT_MultiLineString25D:
		mvtGeomType = vectortile.Tile_LINESTRING
	case gdal.GT_Polygon, gdal.GT_MultiPolygon25D:
		mvtGeomType = vectortile.Tile_POLYGON
	case gdal.GT_Point, gdal.GT_Point25D:
		mvtGeomType = vectortile.Tile_POINT
	default:
		mvtGeomType = vectortile.Tile_UNKNOWN
	}
	return &mvtGeomType
}

func (s *s57Tiler) toMvtFeature(feature *gdal.Feature, tileBounds m.Extrema) *vectortile.Tile_Feature {
	geom := feature.Geometry()
	mvtFeature := vectortile.Tile_Feature{}
	mvtFeature.Type = s.getMvtFeatureType(&geom)

	if *mvtFeature.Type != vectortile.Tile_UNKNOWN {
		mvtFeature.Geometry = s.toMvtGeometry(*mvtFeature.Type, &geom, tileBounds)
		// write tags
		for i := 0; i < feature.FieldCount(); i++ {
			key := feature.FieldDefinition(i).Name()
			value := feature.FieldAsString(i)
			if value != "" {
				fieldType := feature.FieldDefinition(i).Type()
				switch fieldType {
				case gdal.FT_StringList:
					parts := feature.FieldAsStringList(i)
					if parts != nil {
						value = strings.Join(parts, ",")
					}
					break
				default:
					value = feature.FieldAsString(i)
				}

				if _, ok := s.keysMap[key]; !ok {
					s.keysMap[key] = uint32(len(s.keys))
					s.keys = append(s.keys, key)
				}
				if _, ok := s.valuesMap[value]; !ok {
					s.valuesMap[value] = uint32(len(s.values))
					s.values = append(s.values, value)
				}
				mvtFeature.Tags = append(mvtFeature.Tags, s.keysMap[key])
				mvtFeature.Tags = append(mvtFeature.Tags, s.valuesMap[value])
			}
		}

		return &mvtFeature
	}
	return nil
}

func (s *s57Tiler) GetFeatures(layer gdal.Layer, tile m.TileID, tileBounds m.Extrema) []*vectortile.Tile_Feature {

	features := make([]*vectortile.Tile_Feature, 0)

	layer.SetSpatialFilterRect(tileBounds.W, tileBounds.S, tileBounds.E, tileBounds.N)
	for feature := layer.NextFeature(); feature != nil; feature = layer.NextFeature() {
		mvtFeature := s.toMvtFeature(feature, tileBounds)
		if mvtFeature != nil {
			features = append(features, mvtFeature)
		}
	}

	return features
}

func (s *s57Tiler) GetTilesForBounds(tiles map[string]m.TileID, bounds m.Extrema, zoomLevel int) map[string]m.TileID {
	if tiles == nil {
		tiles = make(map[string]m.TileID)
	}
	ulTile := m.Tile(bounds.W, bounds.N, zoomLevel)
	lrTile := m.Tile(bounds.E, bounds.S, zoomLevel)
	for col := ulTile.X; col <= lrTile.X; col++ {
		for row := ulTile.Y; row <= lrTile.Y; row++ {
			key := fmt.Sprintf("%d,%d,%d", col, row, zoomLevel)
			tile := m.TileID{X: col, Y: row, Z: uint64(zoomLevel)}
			tiles[key] = tile
		}
	}
	return tiles
}

func (s *s57Tiler) GetTiles(dataset dataset.Dataset, zoomLevel int) map[string]m.TileID {
	tiles := make(map[string]m.TileID)
	for _, file := range dataset.Files {
		datasource := gdal.OpenDataSource(file.Path, 0)
		for i := 0; i < datasource.LayerCount(); i++ {
			l := datasource.LayerByIndex(i)
			ext, err := l.Extent(true)
			if err == nil {
				tiles = s.GetTilesForBounds(tiles, m.Extrema{W: ext.MinX(), N: ext.MaxY(), E: ext.MaxX(), S: ext.MinY()}, zoomLevel)
			}
		}
		datasource.Release()
	}
	return tiles
}

func (s *s57Tiler) GenerateMetaData(outPath string, dataset dataset.Dataset) {
	path := filepath.Join(outPath, dataset.Id, "sk-chart-meta.json")
	metaData := charts.ChartMetaData{Id: dataset.Id, Name: dataset.Id, Created: time.Now().UTC(), ChartType: "S-57", ChartFormat: "pbf"}

	out, _ := json.Marshal(metaData)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), 0700) // Create your file
	}
	err := os.WriteFile(path, out, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *s57Tiler) GenerateTile(outPath string, dataset dataset.Dataset, tile m.TileID) {
	mvtTile := vectortile.Tile{}

	// for test only selected layers

	layers := []string{"BOYLAT", "BOYCAR", "BOYINB", "BOYISD", "BOYSAW", "BOYSPP", "BCNLAT", "BCNCAR", "BCNISN", "BCNSAW", "BCNSPP", "LIGHTS", "DEPARE", "SEAARE", "COALNE", "RESARE", "UNSARE", "LNDARE", "BUAARE", "NAVLNE", "RECTRC", "CANALS"}

	for _, layerName := range layers {
		ln := layerName
		var version uint32 = 2
		var extent uint32 = TILE_EXTENT
		s.startLayer()
		mvtLayer := vectortile.Tile_Layer{Name: &ln, Version: &version, Extent: &extent}
		for _, file := range dataset.Files {
			datasource := gdal.OpenDataSource(file.Path, 0)

			bounds := m.Bounds(tile)
			tileEnvelope := gdal.Envelope{}
			tileEnvelope.SetMaxX(bounds.E)
			tileEnvelope.SetMaxY(bounds.N)
			tileEnvelope.SetMinX(bounds.W)
			tileEnvelope.SetMinY(bounds.S)

			l := datasource.LayerByName(layerName)

			c, ok := l.FeatureCount(false)
			if ok && c > 0 {
				ext, err := l.Extent(true)
				if err == nil && ext.Intersects(tileEnvelope) {
					features := s.GetFeatures(l, tile, bounds)
					mvtLayer.Features = append(mvtLayer.Features, features...)
				}
			}
			datasource.Release()
		}
		if len(mvtLayer.Features) > 0 {
			// keys
			for _, k := range s.keys {
				mvtLayer.Keys = append(mvtLayer.Keys, k)
			}
			// values
			for _, v := range s.values {
				vv := v
				value := vectortile.Tile_Value{StringValue: &vv}
				mvtLayer.Values = append(mvtLayer.Values, &value)
			}

			mvtTile.Layers = append(mvtTile.Layers, &mvtLayer)
		}
	}

	path := filepath.Join(outPath, dataset.Id, strconv.Itoa(int(tile.X)), strconv.Itoa(int(tile.Y)), strconv.Itoa(int(tile.Z)))
	if len(mvtTile.Layers) > 0 {
		out, _ := proto.Marshal(&mvtTile)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(path), 0700) // Create your file
		}
		err := os.WriteFile(path, out, 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		os.Remove(path)
	}
}
