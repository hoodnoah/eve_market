using System.ComponentModel.DataAnnotations.Schema;
using Microsoft.EntityFrameworkCore;

namespace Api.Models
{
  [Table("region_id")]
  public class RegionId
  {
    public required int Id { get; set; }
    public required string Value { get; set; }
    public ICollection<MarketData> MarketData { get; set; } = new List<MarketData>();
  }
}