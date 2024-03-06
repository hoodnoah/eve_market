using System.ComponentModel.DataAnnotations.Schema;

namespace Api.Models
{
  [Table("completed_dates")]
  public class CompletedDates
  {
    [Column("id")]
    public required int Id { get; set; }
    [Column("date")]
    public required DateOnly Date { get; set; }
    public ICollection<MarketData> MarketData { get; set; } = new List<MarketData>();

  }
}