using Api.Models;
using Microsoft.EntityFrameworkCore;
using Microsoft.OpenApi.Any;
using Microsoft.OpenApi.Models;
using Util;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddDbContext<MarketDb>(options =>
{
    var connectionString =
        builder
        .Configuration
        .GetConnectionString("MySqlLocal")
        ?? throw new ArgumentNullException("MySqlLocal connection string resolved to null.");
    options.UseMySQL(connectionString);
});
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(options =>
{
    options.MapType<DateOnly>(() => new OpenApiSchema
    {
        Type = "string",
        Format = "date",
        Example = new OpenApiString("2021-01-01")
    });
});
builder.Services.AddControllers()
    .AddJsonOptions(options =>
    {
        options.JsonSerializerOptions.Converters.Add(new DateOnlyJSONConverter());
    });

// build app
var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();

Api.Endpoints.CompletedDatesEndpoints.MapEndpoints(app);
Api.Endpoints.RegionEndpoints.MapEndpoints(app);
Api.Endpoints.TypeEndpoints.MapEndpoints(app);
Api.Endpoints.MarketDataEndpoints.MapEndpoints(app);

app.Run();