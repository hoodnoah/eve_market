using System.ComponentModel.DataAnnotations.Schema;
using Microsoft.EntityFrameworkCore;


namespace Api.Models
{
  [Table("market_data")]
  [PrimaryKey(nameof(DateID), nameof(RegionID), nameof(TypeID))]
  public class MarketData
  {
    public int DateID { get; set; }
    public required CompletedDates Date { get; set; }
    public int RegionID { get; set; }
    public required RegionId Region { get; set; }
    public int TypeID { get; set; }
    public required TypeId Type { get; set; }
    public double Average { get; set; }
    public double Highest { get; set; }
    public double Lowest { get; set; }
    public long Volume { get; set; }
    public long OrderCount { get; set; }
  }
}