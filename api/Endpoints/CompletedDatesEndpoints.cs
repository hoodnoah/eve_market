using Api.Models;
using Api.DTO;
using Microsoft.EntityFrameworkCore;

namespace Api.Endpoints
{
  public class CompletedDatesEndpoints
  {
    public static void MapEndpoints(WebApplication app)
    {
      app.MapGet("/api/dates", GetCompletedDatesAsync).WithName("GetCompletedDates").WithOpenApi();
    }

    private static async Task<IResult> GetCompletedDatesAsync(MarketDb db)
    {
      var completedDates = await db
        .CompletedDates
        .Select(cd => new CompletedDatesDTO(cd)
        ).ToListAsync();

      return TypedResults.Ok(completedDates);
    }
  }


}