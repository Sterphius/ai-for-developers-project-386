import "@typespec/http";

using TypeSpec.Http;

namespace Calendar;

@doc("Public scenario for guests: view event types, inspect free slots, and create bookings.")
namespace GuestPublicScenario {
  @doc("List all event types available to guests.")
  @route("/public/event-types")
  @get
  op listEventTypes(): PublicEventType[];

  @doc("List free slots within the next 14 days for the selected event type.")
  @route("/public/event-types/{eventTypeId}/slots")
  @get
  op listAvailableSlots(@path eventTypeId: Identifier): Slot[];

  @doc("Create a booking for a free slot inside the 14-day window.")
  @route("/public/bookings")
  @post
  op createBooking(@body request: CreateBookingRequest): Booking;
}
