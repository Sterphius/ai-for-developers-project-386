import { Link, Outlet } from "react-router-dom";
import { CalendarDays, Settings } from "lucide-react";

export function Layout() {
  return (
    <div className="min-h-screen bg-muted/30">
      <header className="border-b bg-background">
        <div className="container flex h-14 items-center justify-between">
          <Link to="/" className="flex min-w-0 items-center gap-2 font-semibold">
            <CalendarDays className="h-5 w-5 shrink-0" />
            <span className="truncate">Календарь бронирования</span>
          </Link>
          <Link
            to="/admin"
            className="flex shrink-0 items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground"
          >
            <Settings className="h-4 w-4" />
            <span className="hidden sm:inline">Кабинет владельца</span>
          </Link>
        </div>
      </header>
      <main className="container py-8">
        <Outlet />
      </main>
    </div>
  );
}
