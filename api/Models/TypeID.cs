using System.ComponentModel.DataAnnotations.Schema;

namespace Api.Models
{
  [Table("type_id")]
  public class TypeId
  {
    [Column("id")]
    public required int Id { get; set; }
    [Column("value")]
    public required string Value { get; set; }
    public ICollection<MarketData> MarketData { get; set; } = new List<MarketData>();
  }
}