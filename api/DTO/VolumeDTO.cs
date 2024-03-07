using System.Diagnostics.CodeAnalysis;
using Api.Models;
using Util;

namespace Api.DTO
{
  public class VolumeDTO : ItemDTO
  {
    public required long Volume { get; set; }
    public required long OrderCount { get; set; }
  }

  public class DailyVolumeDTO : VolumeDTO
  {
    public required DateOnly Date { get; set; }

  }

  public class MonthlyVolumeDTO : VolumeDTO
  {
    public required MonthDate Month { get; set; }
  }

  public class AnnualVolumeDTO : VolumeDTO
  {
    public required int Year { get; set; }
  }
}