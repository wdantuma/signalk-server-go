using Microsoft.Extensions.Logging;
using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using Google.Protobuf;
using System.IO.Compression;

namespace ShapeTask
{
    public class VectorTileOptimizer
    {
        public ILogger Logger { get; }
        public VectorTileOptimizer(ILoggerFactory loggerFactory)
        {
            Logger = loggerFactory.CreateLogger(typeof(VectorTileOptimizer));
        }

        public void Optimize(string path)
        {
            var tileFiles = Directory.GetFiles(path, "*.pbf", SearchOption.AllDirectories);
            foreach(var tileFile in tileFiles)
            {
                var newTileFile = Path.ChangeExtension(tileFile, "pbf_new");
                VectorTile.Tile vt;
                int inCount = 0;
                int outCount = 0;
                List<ulong> pointCache = new List<ulong>();
                using (var zinput = File.OpenRead(tileFile))
                using (var input = new GZipStream(zinput,CompressionMode.Decompress))
                {
                    using (var zoutput = File.Create(newTileFile))
                    using (var output = new GZipStream(zoutput, CompressionMode.Compress))
                    {
                        vt = VectorTile.Tile.Parser.ParseFrom(input);
                        foreach (var layer in vt.Layers)
                        {
                            int index = 0;
                            int count = inCount = layer.Features.Count;

                            while(index<count)
                            {
                                var feature = layer.Features[index];
                                //only optimize point features
                                if(feature.Type == VectorTile.Tile.Types.GeomType.Point)
                                {
                                    if(feature.Geometry.Count == 3 && feature.Geometry[0] == 9) //  one moveto command
                                    {
                                        var key = feature.Geometry[1] << 32 | feature.Geometry[2];
                                        if(pointCache.Contains(key))
                                        {
                                            layer.Features.RemoveAt(index);
                                            count--;
                                        } else
                                        {
                                            pointCache.Add(key);
                                            index++;
                                            outCount++;
                                        }
                                    }
                                }
                            }
                        }

                        vt.WriteTo(output);
                    }
                }
                //swap
                File.Delete(tileFile);
                File.Move(newTileFile, tileFile);
            }            
        }
    }
}
