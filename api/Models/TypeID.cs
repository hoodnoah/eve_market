using System.ComponentModel.DataAnnotations.Schema;
using Microsoft.EntityFrameworkCore;

namespace Api.Models
{
  [Table("type_id")]
  public class TypeId
  {
    public required int Id { get; set; }
    public required string Value { get; set; }
    public ICollection<MarketData> MarketData { get; set; } = new List<MarketData>();
  }
}