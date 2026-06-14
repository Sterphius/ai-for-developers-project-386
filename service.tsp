import "@typespec/http";

using TypeSpec.Http;

@service({
  title: "Calendar Booking API",
})
@server("http://localhost:8080", "Local development server")
namespace Calendar;
