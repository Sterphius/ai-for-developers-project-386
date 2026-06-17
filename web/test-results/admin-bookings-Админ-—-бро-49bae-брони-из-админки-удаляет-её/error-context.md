# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: admin-bookings.spec.ts >> Админ — бронирования >> 14: отмена брони из админки удаляет её
- Location: e2e/admin-bookings.spec.ts:42:3

# Error details

```
Error: locator.click: Error: strict mode violation: getByLabel('Отменить бронирование') resolved to 2 elements:
    1) <button aria-label="Отменить бронирование" class="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:size-4 [&_svg]:shrink-0 border border-input bg-background hover:bg-accent hover:text-accent-foreground h-9 w-9 shrink-0">…</button> aka getByRole('button', { name: 'Отменить бронирование' }).first()
    2) <button aria-label="Отменить бронирование" class="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:size-4 [&_svg]:shrink-0 border border-input bg-background hover:bg-accent hover:text-accent-foreground h-9 w-9 shrink-0">…</button> aka getByRole('button', { name: 'Отменить бронирование' }).nth(1)

Call log:
  - waiting for getByLabel('Отменить бронирование')

```

# Page snapshot

```yaml
- generic [ref=e3]:
  - banner [ref=e4]:
    - generic [ref=e5]:
      - generic [ref=e6]:
        - img [ref=e7]
        - text: Кабинет владельца
      - generic [ref=e9]:
        - generic [ref=e10]:
          - generic [ref=e11]: Алексей Малышев
          - generic [ref=e12]: owner@example.com · профиль по умолчанию
        - link "К гостю" [ref=e13] [cursor=pointer]:
          - /url: /
          - img [ref=e14]
          - text: К гостю
  - generic [ref=e16]:
    - navigation [ref=e17]:
      - link "Типы событий" [ref=e18] [cursor=pointer]:
        - /url: /admin
        - img [ref=e19]
        - text: Типы событий
      - link "Предстоящие встречи" [ref=e24] [cursor=pointer]:
        - /url: /admin/bookings
        - img [ref=e25]
        - text: Предстоящие встречи
    - main [ref=e28]:
      - generic [ref=e29]:
        - generic [ref=e30]:
          - heading "Предстоящие встречи" [level=1] [ref=e31]
          - paragraph [ref=e32]: "Все бронирования всех типов, по времени. Зона: Asia/Yekaterinburg"
        - generic [ref=e34]:
          - heading "Среда, 17 июня" [level=2] [ref=e35]
          - generic [ref=e36]:
            - generic [ref=e37]:
              - generic [ref=e38]:
                - generic [ref=e39]:
                  - generic [ref=e40]: 16:08 – 16:23
                  - generic [ref=e41]: Test e2e-ab-1781684533255-h7fj
                - generic [ref=e42]:
                  - generic [ref=e43]:
                    - img [ref=e44]
                    - text: Admin Delete Guest
                  - generic [ref=e47]:
                    - img [ref=e48]
                    - text: admin-delete@example.com
              - button "Отменить бронирование" [ref=e51] [cursor=pointer]:
                - img [ref=e52]
            - generic [ref=e55]:
              - generic [ref=e56]:
                - generic [ref=e57]:
                  - generic [ref=e58]: 20:08 – 20:23
                  - generic [ref=e59]: e2e-cr-1781684533240-uxdf
                - generic [ref=e60]:
                  - generic [ref=e61]:
                    - img [ref=e62]
                    - text: Orphan Guest
                  - generic [ref=e65]:
                    - img [ref=e66]
                    - text: orphan@test.com
              - button "Отменить бронирование" [ref=e69] [cursor=pointer]:
                - img [ref=e70]
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
> 57 |     await page.getByLabel("Отменить бронирование").click();
     |                                                    ^ Error: locator.click: Error: strict mode violation: getByLabel('Отменить бронирование') resolved to 2 elements:
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
  74 |     await expect(page.getByText("Бронирований нет")).toBeVisible();
  75 |   });
  76 | });
  77 | 
```