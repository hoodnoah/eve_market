using System.Diagnostics.CodeAnalysis;
using Api.Models;

namespace Api.DTO
{
  public class CompletedDatesDTO
  {
    public required int Id { get; set; }
    public required DateOnly Date { get; set; }

    [SetsRequiredMembers]
    public CompletedDatesDTO(CompletedDates cd)
    {
      Id = cd.Id;
      Date = cd.Date;
    }
  }
}