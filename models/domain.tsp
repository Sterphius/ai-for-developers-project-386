namespace Calendar;

@doc("Owner of the calendar who manages event types and sees all bookings.")
model Owner {
  id: Identifier;
  name: string;
}

@doc("Type of event defined by the calendar owner.")
model EventType {
  id: Identifier;
  name: string;
  description: string;
  durationMinutes: Minutes;
  ownerId: Identifier;
}

@doc("Free or occupied time interval in the booking calendar.")
model Slot {
  startAt: DateTime;
  endAt: DateTime;
  isFree: boolean;
}

@doc("Booking created by a guest for a specific event type and slot.")
model Booking {
  id: Identifier;
  eventTypeId: Identifier;
  slot: Slot;
  createdAt: DateTime;
}

@doc("Booking with event type details for owner views.")
model BookingWithEventType {
  booking: Booking;
  eventType: PublicEventType;
}

@doc("View of an event type that is visible to guests.")
model PublicEventType {
  id: Identifier;
  name: string;
  description: string;
  durationMinutes: Minutes;
}

@doc("Request used by a guest to create a booking.")
model CreateBookingRequest {
  eventTypeId: Identifier;
  slotStartAt: DateTime;
}

@doc("Request used by the owner to create an event type.")
model CreateEventTypeRequest {
  id: Identifier;
  name: string;
  description: string;
  durationMinutes: Minutes;
}

@doc("Standard error envelope for domain and validation failures.")
model ErrorResponse {
  code: string;
  message: string;
}
