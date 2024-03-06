using Api.Models;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;

namespace Api.Endpoints
{
  public class MarketDataEndpoints
  {
    const int MAX_TYPES = 10;
    const int MAX_REGIONS = 10;
    const int MAX_DAYS = 90;
    /// <summary>
    /// Maps the endpoints for the MarketData controller
    /// </summary>
    /// <param name="app"></param>
    public static void MapEndpoints(WebApplication app)
    {
      app.MapPost("/api/market-data/prices/daily", GetDailyPrices).WithName("GetDailyPrices").WithOpenApi();
      app.MapPost("/api/market-data/prices/monthly", GetMonthlyPrices).WithName("GetMonthlyPrices").WithOpenApi();
      app.MapPost("/api/market-data/prices/annual", GetAnnualPrices).WithName("GetAnnualPrices").WithOpenApi();
      app.MapPost("/api/market-data/volume/daily", GetDailyVolumeData).WithName("GetDailyVolumeData").WithOpenApi();
      app.MapPost("/api/market-data/volume/monthly", GetMonthlyVolumeData).WithName("GetMonthlyVolumeData").WithOpenApi();
      app.MapPost("/api/market-data/volume/annual", GetAnnualVolumeData).WithName("GetAnnualVolumeData").WithOpenApi();
    }

    private static async Task<IResult> GetDailyPrices(MarketDb db, [FromBody] MarketDataFilter filter)
    {
      if (!IsValidRegionsFilter(filter.Regions))
      {
        return Results.BadRequest($"Too many regions selected; must be {MAX_REGIONS} or fewer.");
      }

      if (!IsValidTypesFilter(filter.Types))
      {
        return Results.BadRequest($"Too many types selected; must be {MAX_TYPES} or fewer.");
      }

      if (!IsValidDateFilter(filter.DateRange))
      {
        return Results.BadRequest($"Date range is too long for daily data; must be {MAX_DAYS} days or fewer.");
      }

      var results =
        await
        db
        .MarketData
        .Where(m => m.Date.Date >= filter.DateRange.StartDate && m.Date.Date <= filter.DateRange.EndDate)
        .Where(m => filter.Regions.Contains(m.RegionID))
        .Where(m => filter.Types.Contains(m.TypeID))
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
        .GroupBy(m => new { m.Date.Date.Month, m.Date.Date.Year, m.RegionID, m.TypeID })
        .Select(m => new MonthlyPriceDTO
        {
          Month = new MonthDate
          {
            Month = m.Key.Month,
            Year = m.Key.Year
          },
          Region = new RegionDTO
          {
            Id = m.Key.RegionID,
            Name = m.First().Region.Value
          },
          Type = new TypeDTO
          {
            Id = m.Key.TypeID,
            Name = m.First().Type.Value
          },
          Average = m.Average(m => m.Average),
          Highest = m.Max(m => m.Highest),
          Lowest = m.Min(m => m.Lowest)

        })
        .ToListAsync();

      return Results.Ok(results);
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
        .GroupBy(m => new { m.Date.Date.Year, m.RegionID, m.TypeID })
        .Select(m => new AnnualPriceDTO
        {
          Year = m.Key.Year,
          Region = new RegionDTO
          {
            Id = m.Key.RegionID,
            Name = m.First().Region.Value
          },
          Type = new TypeDTO
          {
            Id = m.Key.TypeID,
            Name = m.First().Type.Value
          },
          Average = m.Average(m => m.Average),
          Highest = m.Max(m => m.Highest),
          Lowest = m.Min(m => m.Lowest)

        })
        .ToListAsync();

      return Results.Ok(results);
    }
    private static async Task<IResult> GetDailyVolumeData(MarketDb db, [FromBody] MarketDataFilter filter)
    {
      if (!IsValidDateFilter(filter.DateRange))
      {
        return Results.BadRequest($"Date range is too long for daily data; must be {MAX_DAYS} days or fewer.");
      }
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
        }).ToListAsync();

      return Results.Ok(results);
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
        .GroupBy(m => new { m.Date.Date.Month, m.Date.Date.Year, m.RegionID, m.TypeID })
        .Select(m => new MonthlyVolumeDTO
        {
          Month = new MonthDate
          {
            Month = m.Key.Month,
            Year = m.Key.Year
          },
          Region = new RegionDTO
          {
            Id = m.Key.RegionID,
            Name = m.First().Region.Value,
          },
          Type = new TypeDTO
          {
            Id = m.Key.TypeID,
            Name = m.First().Type.Value
          },
          Volume = m.First().Volume,
          OrderCount = m.First().OrderCount
        }).ToListAsync();

      return Results.Ok(results);
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
        .GroupBy(m => new { m.Date.Date.Year, m.RegionID, m.TypeID })
        .Select(m => new AnnualVolumeDTO
        {
          Year = m.Key.Year,
          Region = new RegionDTO
          {
            Id = m.Key.RegionID,
            Name = m.First().Region.Value,
          },
          Type = new TypeDTO
          {
            Id = m.Key.TypeID,
            Name = m.First().Type.Value
          },
          Volume = m.First().Volume,
          OrderCount = m.First().OrderCount
        }).ToListAsync();

      return Results.Ok(results);
    }



    private static bool IsValidTypesFilter(ICollection<int> types)
    {
      return types.Count <= MAX_TYPES;
    }

    private static bool IsValidRegionsFilter(ICollection<int> regions)
    {
      return regions.Count <= MAX_REGIONS;
    }

    private static bool IsValidDateFilter(DateRange dateRange)
    {
      return dateRange.EndDate.DayNumber - dateRange.StartDate.DayNumber <= MAX_DAYS;
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

  public class DailyPriceDTO
  {
    public required DateOnly Date { get; set; }
    public required RegionDTO Region { get; set; }
    public required TypeDTO Type { get; set; }
    public required double Average { get; set; }
    public required double Highest { get; set; }
    public required double Lowest { get; set; }
  }

  public class MonthlyPriceDTO
  {
    public required MonthDate Month { get; set; }
    public required RegionDTO Region { get; set; }
    public required TypeDTO Type { get; set; }
    public required double Average { get; set; }
    public required double Highest { get; set; }
    public required double Lowest { get; set; }
  }

  public class AnnualPriceDTO
  {
    public required int Year { get; set; }
    public required RegionDTO Region { get; set; }
    public required TypeDTO Type { get; set; }
    public required double Average { get; set; }
    public required double Highest { get; set; }
    public required double Lowest { get; set; }

  }

  public class DailyVolumeDTO
  {
    public required DateOnly Date { get; set; }
    public required RegionDTO Region { get; set; }
    public required TypeDTO Type { get; set; }
    public required long Volume { get; set; }
    public required long OrderCount { get; set; }
  }

  public class MonthlyVolumeDTO
  {
    public required MonthDate Month { get; set; }
    public required RegionDTO Region { get; set; }
    public required TypeDTO Type { get; set; }
    public required long Volume { get; set; }
    public required long OrderCount { get; set; }
  }

  public class AnnualVolumeDTO
  {
    public required int Year { get; set; }
    public required RegionDTO Region { get; set; }
    public required TypeDTO Type { get; set; }
    public required long Volume { get; set; }
    public required long OrderCount { get; set; }
  }

  public class MonthDate
  {
    public required int Month { get; set; }
    public required int Year { get; set; }
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