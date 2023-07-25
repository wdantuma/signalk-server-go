package s57

// Convert S-57 format to MVT vector tiles
// see MVT spec at https://github.com/mapbox/vector-tile-spec/tree/master/2.1

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/lukeroth/gdal"
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

func (s *s57Tiler) toMvtFeature(feature *gdal.Feature, tileBounds m.Extrema) *vectortile.Tile_Feature {
	geom := feature.Geometry()
	mvtFeature := vectortile.Tile_Feature{}
	featureType := vectortile.Tile_POINT
	mvtFeature.Type = &featureType
	x, y, _ := geom.Point(0)
	tx, ty := s.to3857(x, y)

	ulx, uly := s.to3857(tileBounds.W, tileBounds.N)
	lrx, lry := s.to3857(tileBounds.E, tileBounds.S)

	xf := TILE_EXTENT / (lrx - ulx)
	yf := TILE_EXTENT / (uly - lry)

	cmd := (1 & 0x7) | (1 << 3)

	mvtFeature.Geometry = append(mvtFeature.Geometry, uint32(cmd))
	xx := (tx - ulx) * xf
	vx := (uint32(xx) << 1) ^ (uint32(xx) >> 31)
	mvtFeature.Geometry = append(mvtFeature.Geometry, vx)
	yy := (uly - ty) * yf
	vy := (uint32(yy) << 1) ^ (uint32(yy) >> 31)
	mvtFeature.Geometry = append(mvtFeature.Geometry, vy)
	return &mvtFeature
}

func (s *s57Tiler) GetFeatures(layer gdal.Layer, tile m.TileID, tileBounds m.Extrema) []*vectortile.Tile_Feature {

	features := make([]*vectortile.Tile_Feature, 0)

	layer.SetSpatialFilterRect(tileBounds.W, tileBounds.S, tileBounds.E, tileBounds.N)
	for feature := layer.NextFeature(); feature != nil; feature = layer.NextFeature() {
		features = append(features, s.toMvtFeature(feature, tileBounds))
	}

	return features
}

func (s *s57Tiler) GetTiles(dataset dataset.Dataset, zoomLevel int) map[string]m.TileID {
	tiles := make(map[string]m.TileID)
	for _, file := range dataset.Files {
		datasource := gdal.OpenDataSource(file.Path, 0)
		for i := 0; i < datasource.LayerCount(); i++ {
			l := datasource.LayerByIndex(i)
			ext, err := l.Extent(true)
			if err == nil {
				ulTile := m.Tile(ext.MinX(), ext.MaxY(), zoomLevel)
				lrTile := m.Tile(ext.MaxX(), ext.MinY(), zoomLevel)
				for col := ulTile.X; col <= lrTile.X; col++ {
					for row := ulTile.Y; row <= lrTile.Y; row++ {
						key := fmt.Sprintf("%d,%d,%d", col, row, zoomLevel)
						tile := m.TileID{X: col, Y: row, Z: uint64(zoomLevel)}
						tiles[key] = tile
					}
				}

			}
		}
		datasource.Release()
	}
	return tiles
}

func (s *s57Tiler) GenerateTile(outPath string, dataset dataset.Dataset, tile m.TileID) {
	mvtTile := vectortile.Tile{}

	// for test now only buoys

	layers := []string{"BOYLAT"}

	for _, layerName := range layers {
		var version uint32 = 2
		var extent uint32 = TILE_EXTENT
		mvtLayer := vectortile.Tile_Layer{Name: &layerName, Version: &version, Extent: &extent}
		for _, file := range dataset.Files {
			datasource := gdal.OpenDataSource(file.Path, 0)
			defer datasource.Release()

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
		}
		mvtTile.Layers = append(mvtTile.Layers, &mvtLayer)
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
