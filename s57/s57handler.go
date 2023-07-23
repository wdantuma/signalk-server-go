package s57

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lukeroth/gdal"
	m "github.com/wdantuma/signalk-server-go/s57/mercantile"
	"github.com/wdantuma/signalk-server-go/s57/vectortile"
	"google.golang.org/protobuf/proto"
)

const (
	TILE_EXTENT = 4096
)

type s57MvtTileHandler struct {
	transform gdal.CoordinateTransform
}

func (s *s57MvtTileHandler) to3857(x float64, y float64) (float64, float64) {
	xs := make([]float64, 1)
	xs[0] = y
	ys := make([]float64, 1)
	ys[0] = x
	zs := make([]float64, 1)
	zs[0] = 0

	s.transform.Transform(1, xs, ys, zs)

	return xs[0], ys[0]
}

func (s *s57MvtTileHandler) toMvtFeature(feature *gdal.Feature, tile m.TileID) *vectortile.Tile_Feature {
	geom := feature.Geometry()
	mvtFeature := vectortile.Tile_Feature{}
	featureType := vectortile.Tile_POINT
	mvtFeature.Type = &featureType
	x, y, _ := geom.Point(0)
	tx, ty := s.to3857(x, y)

	bounds := m.Bounds(tile)
	ulx, uly := s.to3857(bounds.W, bounds.N)
	lrx, lry := s.to3857(bounds.E, bounds.S)

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

func (s *s57MvtTileHandler) GetLayer(layer gdal.Layer, tile m.TileID, tileEnvelope gdal.Envelope) []byte {

	layerName := layer.Name()
	var version uint32 = 2
	var extent uint32 = TILE_EXTENT
	var out = make([]byte, 0)
	mvtLayer := vectortile.Tile_Layer{Name: &layerName, Version: &version, Extent: &extent}
	n := 0
	layer.SetSpatialFilterRect(tileEnvelope.MinX(), tileEnvelope.MinY(), tileEnvelope.MaxX(), tileEnvelope.MaxY())
	for feature := layer.NextFeature(); feature != nil; feature = layer.NextFeature() {
		n++
		mvtLayer.Features = append(mvtLayer.Features, s.toMvtFeature(feature, tile))
	}
	if n > 0 {
		//fmt.Printf("POLYGON((%f %f,%f %f,%f %f,%f %f,%f %f))\nX:%d,Y:%d f: %d\n", tileEnvelope.MinX(), tileEnvelope.MinY(), tileEnvelope.MaxX(), tileEnvelope.MinY(), tileEnvelope.MaxX(), tileEnvelope.MaxY(), tileEnvelope.MinX(), tileEnvelope.MaxY(), tileEnvelope.MinX(), tileEnvelope.MinY(), tile.X, tile.Y, n)
		mvtTile := vectortile.Tile{}
		mvtTile.Layers = append(mvtTile.Layers, &mvtLayer)
		out, _ = proto.Marshal(&mvtTile)
	}

	return out
}

func NewS57MvtTileHandler() *s57MvtTileHandler {
	src := gdal.CreateSpatialReference("")
	src.FromEPSG(4326)
	dst := gdal.CreateSpatialReference("")
	dst.FromEPSG(3857)

	driver, err := gdal.GetDriverByName("S57")
	if err != nil {
		fmt.Println(err.Error())
	}
	driver.Register()

	return &s57MvtTileHandler{transform: gdal.CreateCoordinateTransform(src, dst)}
}

func (s *s57MvtTileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.mapbox-vector-tile")

	datasource := gdal.OpenDataSource("../../charts/20230721_U7Inland_Waddenzee_week 29_NL/ENC_ROOT/1R/7/1R7WAD01/1R7WAD01.000", 0)
	defer datasource.Release()

	vars := mux.Vars(r)
	x, _ := strconv.ParseInt(vars["x"], 10, 64)
	y, _ := strconv.ParseInt(vars["y"], 10, 64)
	z, _ := strconv.ParseUint(vars["z"], 10, 64)
	tile := m.TileID{X: x, Y: y, Z: z}

	bounds := m.Bounds(tile)
	tileEnvelope := gdal.Envelope{}
	tileEnvelope.SetMaxX(bounds.E)
	tileEnvelope.SetMaxY(bounds.N)
	tileEnvelope.SetMinX(bounds.W)
	tileEnvelope.SetMinY(bounds.S)

	// for test now only buoys

	l := datasource.LayerByName("BOYLAT")

	c, ok := l.FeatureCount(false)
	if ok && c > 0 {
		ext, err := l.Extent(true)
		if err == nil && ext.Intersects(tileEnvelope) {
			out := s.GetLayer(l, tile, tileEnvelope)
			if len(out) > 0 {
				w.Write(out)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
