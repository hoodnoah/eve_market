using System.Diagnostics.CodeAnalysis;
using Api.Models;

namespace Api.DTO
{
  public class RegionDTO
  {
    public required int Id { get; set; }
    public required string Name { get; set; }
  }
}