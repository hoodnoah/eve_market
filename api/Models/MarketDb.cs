using Microsoft.EntityFrameworkCore;

namespace Api.Models
{
  public class MarketDb : DbContext
  {
    public MarketDb(DbContextOptions<MarketDb> options) : base(options) { }

    public DbSet<CompletedDates> CompletedDates { get; set; }
    public DbSet<RegionId> RegionId { get; set; }
    public DbSet<TypeId> TypeId { get; set; }
    public DbSet<MarketData> MarketData { get; set; }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
      modelBuilder.Entity<MarketData>()
        .HasOne(m => m.Date)
        .WithMany(d => d.MarketData)
        .HasForeignKey(m => m.DateID);

      modelBuilder.Entity<MarketData>()
        .HasOne(m => m.Region)
        .WithMany(r => r.MarketData)
        .HasForeignKey(m => m.RegionID);

      modelBuilder.Entity<MarketData>()
        .HasOne(m => m.Type)
        .WithMany(t => t.MarketData)
        .HasForeignKey(m => m.TypeID);

      // register cast from DateOnly -> DateTime, vice-versa
      modelBuilder.Entity<CompletedDates>(builder =>
      {
        builder.Property(x => x.Date)
          .HasConversion<Util.DateOnlyConverter>();
      });
    }
  }
}