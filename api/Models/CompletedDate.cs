using Microsoft.EntityFrameworkCore;
using System.ComponentModel.DataAnnotations.Schema;

namespace Api.Models
{
  [Table("completed_dates")]
  public class CompletedDates
  {
    public required int Id { get; set; }
    public required DateOnly Date { get; set; }
    public ICollection<MarketData> MarketData { get; set; } = new List<MarketData>();

  }
}