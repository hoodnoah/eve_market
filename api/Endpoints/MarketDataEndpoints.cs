using Api.Models;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;

namespace Api.Endpoints
{
  public class MarketDataEndpoints
  {
    public static void MapEndpoints(WebApplication app)
    {
      app.MapPost("/api/market-data/filter", GetFilteredMarketData).WithName("GetFilteredMarketData").WithOpenApi();
    }

    private static async Task<IResult> GetFilteredMarketData(MarketDb db, [FromBody] MarketDataFilter filter)
    {
      var query = db.MarketData.AsQueryable();

      // filter dates
      query = query.Where(m => m.Date.Date >= filter.DateRange.StartDate && m.Date.Date <= filter.DateRange.EndDate);

      // filter regions
      query = query.Where(m => filter.Regions.Contains(m.RegionID));

      // filter
      query = query.Where(m => filter.Types.Contains(m.TypeID));

      var results = await query
        .Select(m => new MarketDataDTO
        {
          Date = m.Date.Date,
          Region = new RegionDTO
          {
            Id = m.RegionID,
            Name = m.Region.Value
          },
          Type = new TypeDTO
          {
            Id = m.TypeID,
            Name = m.Type.Value
          },
          Average = m.Average,
          Highest = m.Highest,
          Lowest = m.Lowest,
          Volume = m.Volume,
          OrderCount = m.OrderCount
        }).ToListAsync();

      return Results.Ok(results);
    }
  }

  public class MarketDataDTO
  {
    public required DateOnly Date { get; set; }
    public required RegionDTO Region { get; set; }
    public required TypeDTO Type { get; set; }
    public required double Average { get; set; }
    public required double Highest { get; set; }
    public required double Lowest { get; set; }
    public required long Volume { get; set; }
    public required long OrderCount { get; set; }
  }

  public class MarketDataFilter
  {
    public required DateRange DateRange { get; set; }
    public required List<int> Regions { get; set; } = new List<int>();
    public required List<int> Types { get; set; } = new List<int>();

  }

  public class DateRange
  {
    public required DateOnly StartDate { get; set; }
    public required DateOnly EndDate { get; set; }
  }
}