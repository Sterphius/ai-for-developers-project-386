import "@typespec/http";

using TypeSpec.Http;

namespace Calendar;

@doc("Owner scenario: manage event types and inspect all upcoming bookings.")
namespace OwnerScenario {
  @doc("Create a new event type.")
  @route("/owner/event-types")
  @post
  op createEventType(@body request: CreateEventTypeRequest): EventType;

  @doc("List all bookings across all event types.")
  @route("/owner/bookings")
  @get
  op listBookings(): Booking[];
}
