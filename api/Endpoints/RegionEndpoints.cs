using Api.Models;
using Microsoft.EntityFrameworkCore;

namespace Api.Endpoints
{
  public class RegionEndpoints
  {
    public static void MapEndpoints(WebApplication app)
    {
      app.MapGet("/api/regions", GetRegionsAsync).WithName("GetRegions").WithOpenApi();
    }

    private static async Task<IResult> GetRegionsAsync(MarketDb db)
    {
      var regions = await db
        .RegionId
        .Select(r => new RegionDTO
        {
          Id = r.Id,
          Name = r.Value
        }).ToListAsync();
      return Results.Ok(regions);
    }
  }

  public class RegionDTO
  {
    public required int Id { get; set; }
    public required string Name { get; set; }
  }
}