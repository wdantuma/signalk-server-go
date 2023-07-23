using Microsoft.Extensions.Logging;
using System;
using System.Collections.Generic;
using System.Drawing;
using System.IO;
using System.IO.Compression;
using System.Text;
using Base.Core.Common.Render;
using Color = System.Drawing.Color;
using Base.Core.Common.Render.Model;
using TileLayer = VectorTile.Tile.Types.Layer;
using Google.Protobuf.Collections;
using System.Linq;
using ICSharpCode.SharpZipLib.GZip;

namespace ShapeTask
{
    public class VectorTileRenderer
    {
        public ILogger Logger { get; }
        public VectorTileRenderer(ILoggerFactory loggerFactory)
        {
            Logger = loggerFactory.CreateLogger(typeof(VectorTileRenderer));
        }

        private void DrawPoint(Graphics graphics, uint extent, Pen pen, RepeatedField<uint> geometry)
        {
            float cX, cY;
            cX = cY = 0F;
            float f = 256F / extent;
            int p = 0;
            if(geometry.Count>0)
            {
                var id = geometry[p] & 0x7;
                if (id != 1) throw new Exception("Point feature only expects a moveto command");
                var count = geometry[p] >> 3;
                p++;
                while (p < geometry.Count)
                {
                    var dx = ((geometry[p] >> 1) ^ (-(geometry[p] & 1)));
                    var dy = ((geometry[p + 1] >> 1) ^ (-(geometry[p + 1] & 1)));
                    var pX = (cX + dx);
                    var pY = (cY + dy);
                    graphics.FillEllipse(new SolidBrush(pen.Color), pX*f, pY*f, 5, 5);
                    cX = pX;
                    cY = pY;
                    p += 2;
                }
            }
        }

        private void DrawPolygon(Graphics graphics, uint extent, Pen pen, RepeatedField<uint> geometry)
        {
            float cX, cY;
            cX = cY = 0F;
            float f = 256F / extent;
            int p = 0;
            List<PointF> ring = new List<PointF>();
            if (geometry.Count > 0)
            {
                while(p<geometry.Count)
                {
                    var id = geometry[p] & 0x7;
                    var count = geometry[p] >> 3;
                    p++;
                    for(int n=0;n<count;n++)
                    {
                        if(id==1 || id == 2)
                        {
                            var dx = ((geometry[p] >> 1) ^ (-(geometry[p] & 1)));
                            var dy = ((geometry[p + 1] >> 1) ^ (-(geometry[p + 1] & 1)));
                            var pX = (cX + dx);
                            var pY = (cY + dy);
                            p += 2;
                            cX = pX;
                            cY = pY;
                            if (id == 1)
                            {
                                ring.Clear();
                                ring.Add(new PointF(pX*f, pY*f));
                            }
                            else if (id == 2)
                            {
                                ring.Add(new PointF(pX*f, pY*f));
                            }
                        } else
                        {
                            ring.Add(ring[0]);
                            graphics.FillPolygon(new SolidBrush(pen.Color), ring.ToArray());
                            ring.Clear();
                        }                       
                    }                  
                }              
            }
        }

        private void DrawLinestring(Graphics graphics, uint extent, Pen pen, RepeatedField<uint> geometry)
        {
            float cX, cY;
            cX = cY = 0F;
            float f = 256F / extent;
            int p = 0;
            List<PointF> ring = new List<PointF>();
            if (geometry.Count > 0)
            {
                while (p < geometry.Count)
                {
                    var id = geometry[p] & 0x7;
                    var count = geometry[p] >> 3;
                    p++;
                    for (int n = 0; n < count; n++)
                    {
                        if (id == 1 || id == 2)
                        {
                            var dx = ((geometry[p] >> 1) ^ (-(geometry[p] & 1)));
                            var dy = ((geometry[p + 1] >> 1) ^ (-(geometry[p + 1] & 1)));
                            var pX = (cX + dx);
                            var pY = (cY + dy);
                            p += 2;
                            cX = pX;
                            cY = pY;
                            if (id == 1)
                            {
                                ring.Clear();
                                ring.Add(new PointF(pX * f, pY * f));
                            }
                            else if (id == 2)
                            {
                                ring.Add(new PointF(pX * f, pY * f));
                            }
                        } else
                        {
                            throw new Exception("Close path not expected in linestring");
                        }                       
                    }
                }               
            }
        }

        public void Render(TileLayer tileLayer,Layer layer,VectorTile.Tile.Types.Feature feature, Graphics graphics, uint extent)
        {
            float x, y;
            x = y = 0F;
            float f = 256F / extent;
            var pen = new Pen(Color.Black,1);
            if( layer is object && layer.Renderer is Base.Core.Common.Render.Model.IColorMapProvider && layer.Renderer is Base.Core.Common.Render.Model.IBandProvider)
            {
                var gradient = (layer.Renderer as IColorMapProvider).ColorMap.Gradient;
                var band = (layer.Renderer as IBandProvider).Band;
                var bandIndex = 0;
                if(!string.IsNullOrEmpty(layer.IndexKey))
                {
                    bandIndex = tileLayer.Keys.ToList().IndexOf(layer.IndexKey);
                } else
                {
                    bandIndex = band.Index;
                }

                var valueIndex = (int)feature.Tags[(bandIndex * 2) + 1];

                var allValues = new[] { tileLayer.Values[valueIndex].DoubleValue, (double)tileLayer.Values[valueIndex].UintValue, (double)tileLayer.Values[valueIndex].IntValue, (double)tileLayer.Values[valueIndex].FloatValue, (double)tileLayer.Values[valueIndex].SintValue };
                var value = allValues.Where(v => v != 0.0).SingleOrDefault();
               
                var color = gradient.GetColorAt(band.Histogram.ValueToFactor(value));
                pen = new Pen(Color.FromArgb(color.Alpha, color.Red, color.Green, color.Blue));
            }
            switch(feature.Type)
            {
                case VectorTile.Tile.Types.GeomType.Point: DrawPoint(graphics, extent, pen, feature.Geometry);break;
                case VectorTile.Tile.Types.GeomType.Polygon: DrawPolygon(graphics, extent, pen, feature.Geometry); break;
                case VectorTile.Tile.Types.GeomType.Linestring: DrawLinestring(graphics, extent, pen, feature.Geometry); break;
                default:break;
            }           
        }

        public void Render(Layer layer,VectorTile.Tile tile, Graphics graphics)
        {
            foreach(var l in tile.Layers)
            {
                foreach(var feature in l.Features)
                {                  
                    Render(l,layer,feature, graphics, l.Extent);
                }
            }
        }

        public void Render(Layer layer, string basePath)
        {
            var tileFiles = Directory.GetFiles(basePath, "*.pbf", SearchOption.AllDirectories);
            List<ulong> pointCache = new List<ulong>();
            foreach (var tileFile in tileFiles)
            {
                var pathSuffix = Base.Core.Common.Helper.GetPlatformIndependentDirectoryPath(tileFile).Replace(basePath, "");
                var pngTileFile =  Path.ChangeExtension($"{basePath}/image_tiles_{layer.Index}{pathSuffix}", ".png");
                Directory.CreateDirectory(Path.GetDirectoryName(pngTileFile));
                using (var zinput = File.OpenRead(tileFile))                    
                using (var input = new GZipInputStream(zinput))
                {
                    using (var output = File.Create(pngTileFile))
                    {
                        var vt = VectorTile.Tile.Parser.ParseFrom(input);
                        var bitmap = new Bitmap(256, 256);
                        bitmap.MakeTransparent();
                        using (var graphics = Graphics.FromImage(bitmap))
                        {
                            // graphics.SmoothingMode = System.Drawing.Drawing2D.SmoothingMode.AntiAlias;
                            Render(layer, vt, graphics);
                        }
                        bitmap.Save(output, System.Drawing.Imaging.ImageFormat.Png);
                    }
                }           
            }
        }
    }
}
