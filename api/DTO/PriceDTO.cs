using System.Diagnostics.CodeAnalysis;
using Api.Models;
using Util;

namespace Api.DTO
{
  public class PriceDTO : ItemDTO
  {
    public required double Average { get; set; }
    public required double Highest { get; set; }
    public required double Lowest { get; set; }
  }

  public class DailyPriceDTO : PriceDTO
  {
    public required DateOnly Date { get; set; }
  }

  public class MonthlyPriceDTO : PriceDTO
  {
    public required MonthDate Month { get; set; }
  }

  public class AnnualPriceDTO : PriceDTO
  {
    public required int Year { get; set; }
  }

  public class MonthDate
  {
    public required int Month { get; set; }
    public required int Year { get; set; }
  }
}