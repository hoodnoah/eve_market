using Api.Models;
using Microsoft.EntityFrameworkCore;

namespace Api.Endpoints
{
  public class CompletedDates
  {
    public static void MapEndpoints(WebApplication app)
    {
      app.MapGet("/api/dates", GetCompletedDatesAsync).WithName("GetCompletedDates").WithOpenApi();
    }

    private static async Task<IResult> GetCompletedDatesAsync(MarketDb db)
    {
      var completedDates = await db
        .CompletedDates
        .Select(cd => new CompletedDatesDTO
        {
          Id = cd.Id,
          Date = cd.Date
        }).ToListAsync();

      return Results.Ok(completedDates);
    }
  }

  public class CompletedDatesDTO
  {
    public required int Id { get; set; }
    public required DateOnly Date { get; set; }
  }
}