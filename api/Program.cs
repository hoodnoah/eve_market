using Api.Models;
using Microsoft.EntityFrameworkCore;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddDbContext<MarketDb>(options =>
{
    var connectionString = builder.Configuration.GetConnectionString("MySqlLocal");
    var serverVersion = new MySqlServerVersion(new Version(8, 3, 0));
    options.UseMySql(connectionString, serverVersion).EnableSensitiveDataLogging().EnableDetailedErrors();
});
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

// build app
var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();

app.MapGet("/hello", () => "Hello World!");

app.MapGet("/regions", async (MarketDb dbContext) =>
{
    return await dbContext.RegionId.ToListAsync();
}).WithName("GetRegions").WithOpenApi();

app.Run();