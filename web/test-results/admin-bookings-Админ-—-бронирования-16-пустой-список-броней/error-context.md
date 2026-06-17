# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: admin-bookings.spec.ts >> Админ — бронирования >> 16: пустой список броней
- Location: e2e/admin-bookings.spec.ts:71:3

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByText('Бронирований нет')
Expected: visible
Timeout: 5000ms
Error: element(s) not found

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for getByText('Бронирований нет')

```

```yaml
- banner:
  - img
  - text: Кабинет владельца Алексей Малышев owner@example.com · профиль по умолчанию
  - link "К гостю":
    - /url: /
    - img
    - text: К гостю
- navigation:
  - link "Типы событий":
    - /url: /admin
    - img
    - text: Типы событий
  - link "Предстоящие встречи":
    - /url: /admin/bookings
    - img
    - text: Предстоящие встречи
- main:
  - heading "Предстоящие встречи" [level=1]
  - paragraph: "Все бронирования всех типов, по времени. Зона: Asia/Yekaterinburg"
  - heading "Среда, 17 июня" [level=2]
  - text: 13:23 – 13:53 Вводный звонок
  - img
  - text: Тест Гость
  - img
  - text: test@example.com
  - button "Отменить бронирование":
    - img
  - text: 16:08 – 16:23 Test e2e-ab-1781684533255-h7fj
  - img
  - text: Admin Delete Guest
  - img
  - text: admin-delete@example.com
  - button "Отменить бронирование":
    - img
```

# Test source

```ts
  1  | import { test, expect } from "@playwright/test";
  2  | import {
  3  |   createEventType,
  4  |   createBooking,
  5  |   deleteBooking,
  6  |   deleteEventType,
  7  |   futureSlotMinutes,
  8  | } from "./helpers";
  9  | 
  10 | const uid = () =>
  11 |   `e2e-ab-${Date.now()}-${Math.random().toString(36).slice(2, 6)}`;
  12 | 
  13 | let offset = 120;
  14 | 
  15 | test.describe("Админ — бронирования", () => {
  16 |   // Сценарий 13: бронь отображается в админке
  17 |   test("13: созданная бронь видна на странице предстоящих встреч", async ({
  18 |     page,
  19 |   }) => {
  20 |     offset += 20;
  21 |     const eventTypeId = uid();
  22 |     await createEventType(eventTypeId, { durationMinutes: 15 });
  23 |     const start = futureSlotMinutes(offset, 15);
  24 |     const booking = await createBooking(
  25 |       eventTypeId,
  26 |       start,
  27 |       "Admin Test Guest",
  28 |       "admin-test@example.com",
  29 |     );
  30 | 
  31 |     await page.goto("/admin/bookings");
  32 |     await expect(page.getByText("Admin Test Guest")).toBeVisible();
  33 |     await expect(
  34 |       page.getByText("admin-test@example.com"),
  35 |     ).toBeVisible();
  36 | 
  37 |     await deleteBooking(booking.id);
  38 |     await deleteEventType(eventTypeId);
  39 |   });
  40 | 
  41 |   // Сценарий 14: удаление брони из админки
  42 |   test("14: отмена брони из админки удаляет её", async ({ page }) => {
  43 |     offset += 20;
  44 |     const eventTypeId = uid();
  45 |     await createEventType(eventTypeId, { durationMinutes: 15 });
  46 |     const start = futureSlotMinutes(offset, 15);
  47 |     const booking = await createBooking(
  48 |       eventTypeId,
  49 |       start,
  50 |     "Admin Delete Guest",
  51 |       "admin-delete@example.com",
  52 |     );
  53 | 
  54 |     await page.goto("/admin/bookings");
  55 |     await expect(page.getByText("Admin Delete Guest")).toBeVisible();
  56 | 
  57 |     await page.getByLabel("Отменить бронирование").click();
  58 |     await page.getByRole("button", { name: "Отменить" }).click();
  59 | 
  60 |     await expect(
  61 |       page.getByText("Бронирование отменено."),
  62 |     ).toBeVisible();
  63 |     await expect(
  64 |       page.getByText("Admin Delete Guest"),
  65 |     ).not.toBeVisible();
  66 | 
  67 |     await deleteEventType(eventTypeId);
  68 |   });
  69 | 
  70 |   // Сценарий 16: список броней пуст
  71 |   test("16: пустой список броней", async ({ page }) => {
  72 |     // Не создаём бронь — идём сразу в пустой список
  73 |     await page.goto("/admin/bookings");
> 74 |     await expect(page.getByText("Бронирований нет")).toBeVisible();
     |                                                      ^ Error: expect(locator).toBeVisible() failed
  75 |   });
  76 | });
  77 | 
```