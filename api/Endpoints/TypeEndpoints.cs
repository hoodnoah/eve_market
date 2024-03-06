using Api.Models;
using Microsoft.EntityFrameworkCore;

namespace Api.Endpoints
{
  public class TypeEndpoints
  {
    public static void MapEndpoints(WebApplication app)
    {
      app.MapGet("/api/types", GetTypesAsync).WithName("GetTypes").WithOpenApi();
    }

    private static async Task<IResult> GetTypesAsync(MarketDb db)
    {
      var types = await db
        .TypeId
        .Select(t => new TypeDTO
        {
          Id = t.Id,
          Name = t.Value
        }).ToListAsync();
      return Results.Ok(types);
    }
  }

  public class TypeDTO
  {
    public required int Id { get; set; }
    public required string Name { get; set; }
  }
}