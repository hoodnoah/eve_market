using Api.DTO;
using Api.Models;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Util;

namespace Api.Endpoints
{
  public class MarketDataEndpoints
  {
    const int MAX_TYPES = 5;
    const int MAX_REGIONS = 5;
    const int PAGESIZE = 30;

    /// <summary>
    /// Maps the endpoints for the MarketData controller
    /// </summary>
    /// <param name="app"></param>
    public static void MapEndpoints(WebApplication app)
    {
      // register the whole /api/market-data endpoint for prices, volume
      var prices = app.MapGroup("/api/market-data/prices");
      var volume = app.MapGroup("/api/market-data/volume");

      prices.MapPost("/daily", GetDailyPrices).WithName("GetDailyPrices").WithOpenApi();
      prices.MapPost("/monthly", GetMonthlyPrices).WithName("GetMonthlyPrices").WithOpenApi();
      prices.MapPost("/annual", GetAnnualPrices).WithName("GetAnnualPrices").WithOpenApi();
      volume.MapPost("/daily", GetDailyVolumeData).WithName("GetDailyVolumeData").WithOpenApi();
      volume.MapPost("/monthly", GetMonthlyVolumeData).WithName("GetMonthlyVolumeData").WithOpenApi();
      volume.MapPost("/annual", GetAnnualVolumeData).WithName("GetAnnualVolumeData").WithOpenApi();
    }

    private static async Task<IResult> GetDailyPrices(MarketDb db, [FromBody] PaginatedMarketDataFilter filter)
    {
      if (!IsValidRegionsFilter(filter.Regions))
      {
        return Results.BadRequest($"Too many regions selected; must be {MAX_REGIONS} or fewer.");
      }

      if (!IsValidTypesFilter(filter.Types))
      {
        return Results.BadRequest($"Too many types selected; must be {MAX_TYPES} or fewer.");
      }

      var results =
        await
        db
        .MarketData
        .Where(m => m.Date.Date >= filter.DateRange.StartDate && m.Date.Date <= filter.DateRange.EndDate)
        .Where(m => filter.Regions.Contains(m.RegionID))
        .Where(m => filter.Types.Contains(m.TypeID))
        .Where(m => m.Date.Date < filter.BeforeDate)
        .OrderByDescending(m => m.Date.Date)
        .Take(PAGESIZE)
        .Select(m => new DailyPriceDTO
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
          Lowest = m.Lowest
        })
        .ToListAsync();

      return TypedResults.Ok(results);
    }

    private static async Task<IResult> GetMonthlyPrices(MarketDb db, [FromBody] MarketDataFilter filter)
    {
      if (!IsValidRegionsFilter(filter.Regions))
      {
        return Results.BadRequest($"Too many regions selected; must be {MAX_REGIONS} or fewer.");
      }

      if (!IsValidTypesFilter(filter.Types))
      {
        return Results.BadRequest($"Too many types selected; must be {MAX_TYPES} or fewer.");
      }

      var results =
        await
        db
        .MarketData
        .Where(m => m.Date.Date >= filter.DateRange.StartDate && m.Date.Date <= filter.DateRange.EndDate)
        .Where(m => filter.Regions.Contains(m.RegionID))
        .Where(m => filter.Types.Contains(m.TypeID))
        .GroupBy(m => new MonthlyGroupKey { Month = m.Date.Date.Month, Year = m.Date.Date.Year, RegionID = m.RegionID, TypeID = m.TypeID })
        .Select(g => new MonthlyPriceDTO
        {
          Month = new MonthDate
          {
            Month = g.First().Date.Date.Month,
            Year = g.First().Date.Date.Year
          },
          Region = new RegionDTO
          {
            Id = g.First().RegionID,
            Name = g.First().Region.Value
          },
          Type = new TypeDTO
          {
            Id = g.First().TypeID,
            Name = g.First().Type.Value
          },
          Average = g.Average(g => g.Average),
          Highest = g.Max(g => g.Highest),
          Lowest = g.Min(g => g.Lowest)
        })
        .ToListAsync();

      return TypedResults.Ok(results);
    }

    private static async Task<IResult> GetAnnualPrices(MarketDb db, [FromBody] MarketDataFilter filter)
    {
      if (!IsValidRegionsFilter(filter.Regions))
      {
        return Results.BadRequest($"Too many regions selected; must be {MAX_REGIONS} or fewer.");
      }

      if (!IsValidTypesFilter(filter.Types))
      {
        return Results.BadRequest($"Too many types selected; must be {MAX_TYPES} or fewer.");
      }

      var results =
        await
        db
        .MarketData
        .Where(m => m.Date.Date >= filter.DateRange.StartDate && m.Date.Date <= filter.DateRange.EndDate)
        .Where(m => filter.Regions.Contains(m.RegionID))
        .Where(m => filter.Types.Contains(m.TypeID))
        .GroupBy(g => new AnnualGroupKey { Year = g.Date.Date.Year, RegionID = g.RegionID, TypeID = g.TypeID })
        .Select(g => new AnnualPriceDTO
        {
          Year = g.First().Date.Date.Year,
          Region = new RegionDTO
          {
            Id = g.First().RegionID,
            Name = g.First().Region.Value
          },
          Type = new TypeDTO
          {
            Id = g.First().TypeID,
            Name = g.First().Type.Value
          },
          Average = g.Average(g => g.Average),
          Highest = g.Max(g => g.Highest),
          Lowest = g.Min(g => g.Lowest)
        })
        .ToListAsync();

      return TypedResults.Ok(results);
    }

    private static async Task<IResult> GetDailyVolumeData(MarketDb db, [FromBody] PaginatedMarketDataFilter filter)
    {
      if (!IsValidRegionsFilter(filter.Regions))
      {
        return Results.BadRequest($"Too many regions selected; must be {MAX_REGIONS} or fewer.");
      }

      if (!IsValidTypesFilter(filter.Types))
      {
        return Results.BadRequest($"Too many types selected; must be {MAX_TYPES} or fewer.");
      }

      var results =
        await
        db
        .MarketData
        .Where(m => m.Date.Date >= filter.DateRange.StartDate && m.Date.Date <= filter.DateRange.EndDate)
        .Where(m => filter.Regions.Contains(m.RegionID))
        .Where(m => filter.Types.Contains(m.TypeID))
        .Where(m => m.Date.Date < filter.BeforeDate)
        .Take(PAGESIZE)
        .Select(m => new DailyVolumeDTO
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
          Volume = m.Volume,
          OrderCount = m.OrderCount
        })
        .ToListAsync();

      return TypedResults.Ok(results);
    }

    private static async Task<IResult> GetMonthlyVolumeData(MarketDb db, [FromBody] MarketDataFilter filter)
    {
      if (!IsValidRegionsFilter(filter.Regions))
      {
        return Results.BadRequest($"Too many regions selected; must be {MAX_REGIONS} or fewer.");
      }
      if (!IsValidTypesFilter(filter.Types))
      {
        return Results.BadRequest($"Too many types selected; must be {MAX_TYPES} or fewer.");
      }

      var results =
        await
        db
        .MarketData
        .Where(m => m.Date.Date >= filter.DateRange.StartDate && m.Date.Date <= filter.DateRange.EndDate)
        .Where(m => filter.Regions.Contains(m.RegionID))
        .Where(m => filter.Types.Contains(m.TypeID))
        .GroupBy(m => new MonthlyGroupKey { Month = m.Date.Date.Month, Year = m.Date.Date.Year, RegionID = m.RegionID, TypeID = m.TypeID })
        .Select(g => new MonthlyVolumeDTO
        {
          Month = new MonthDate
          {
            Month = g.First().Date.Date.Month,
            Year = g.First().Date.Date.Year
          },
          Region = new RegionDTO
          {
            Id = g.First().RegionID,
            Name = g.First().Region.Value
          },
          Type = new TypeDTO
          {
            Id = g.First().TypeID,
            Name = g.First().Type.Value
          },
          Volume = g.Sum(g => g.Volume),
          OrderCount = g.Sum(g => g.OrderCount)
        })
        .ToListAsync();

      return TypedResults.Ok(results);
    }

    private static async Task<IResult> GetAnnualVolumeData(MarketDb db, [FromBody] MarketDataFilter filter)
    {
      if (!IsValidRegionsFilter(filter.Regions))
      {
        return Results.BadRequest($"Too many regions selected; must be {MAX_REGIONS} or fewer.");
      }
      if (!IsValidTypesFilter(filter.Types))
      {
        return Results.BadRequest($"Too many types selected; must be {MAX_TYPES} or fewer.");
      }

      var results =
        await
        db
        .MarketData
        .Where(m => m.Date.Date >= filter.DateRange.StartDate && m.Date.Date <= filter.DateRange.EndDate)
        .Where(m => filter.Regions.Contains(m.RegionID))
        .Where(m => filter.Types.Contains(m.TypeID))
        .GroupBy(m => new AnnualGroupKey { Year = m.Date.Date.Year, RegionID = m.RegionID, TypeID = m.TypeID })
        .Select(g => new AnnualVolumeDTO
        {
          Year = g.First().Date.Date.Year,
          Region = new RegionDTO
          {
            Id = g.First().RegionID,
            Name = g.First().Region.Value
          },
          Type = new TypeDTO
          {
            Id = g.First().TypeID,
            Name = g.First().Type.Value
          },
          Volume = g.Sum(g => g.Volume),
          OrderCount = g.Sum(g => g.OrderCount)
        })
        .ToListAsync();

      return TypedResults.Ok(results);
    }

    private static bool IsValidTypesFilter(ICollection<int> types)
    {
      return types.Count <= MAX_TYPES;
    }

    private static bool IsValidRegionsFilter(ICollection<int> regions)
    {
      return regions.Count <= MAX_REGIONS;
    }
  }

  public class MarketDataFilter
  {
    public required DateRange DateRange { get; set; }
    public required List<int> Regions { get; set; } = new List<int>();
    public required List<int> Types { get; set; } = new List<int>();
  }

  public class PaginatedMarketDataFilter : MarketDataFilter
  {
    public required DateOnly BeforeDate { get; set; }
  }

  public class DateRange
  {
    public required DateOnly StartDate { get; set; }
    public required DateOnly EndDate { get; set; }
  }
}

