using System.ComponentModel.DataAnnotations.Schema;
using Microsoft.EntityFrameworkCore;


namespace Api.Models
{
  [Table("market_data")]
  [PrimaryKey(nameof(DateID), nameof(RegionID), nameof(TypeID))]
  public class MarketData
  {
    [Column("date_id")]
    public required int DateID { get; set; }
    public required CompletedDates Date { get; set; }
    [Column("region_id")]
    public required int RegionID { get; set; }
    public required RegionId Region { get; set; }
    [Column("type_id")]
    public required int TypeID { get; set; }
    public required TypeId Type { get; set; }
    [Column("average")]
    public required double Average { get; set; }
    [Column("highest")]
    public required double Highest { get; set; }
    [Column("lowest")]
    public required double Lowest { get; set; }
    [Column("volume")]
    public required long Volume { get; set; }
    [Column("order_count")]
    public required long OrderCount { get; set; }
  }
}