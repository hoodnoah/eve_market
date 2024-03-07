namespace Util
{
  public class MonthlyGroupKey
  {
    public required int Month { get; set; }
    public required int Year { get; set; }
    public required int RegionID { get; set; }
    public required int TypeID { get; set; }
  }

  public class AnnualGroupKey
  {
    public required int Year { get; set; }
    public required int RegionID { get; set; }
    public required int TypeID { get; set; }
  }
}