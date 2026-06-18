import { Link, NavLink, Outlet } from "react-router-dom";
import { CalendarDays, ListChecks, LayoutGrid, ArrowLeft } from "lucide-react";
import { useOwner } from "@/api/hooks";
import { cn } from "@/lib/utils";

export function AdminLayout() {
  const { data: owner } = useOwner();

  return (
    <div className="min-h-screen bg-muted/30">
      <header className="border-b bg-background">
        <div className="container flex h-14 items-center justify-between">
          <div className="flex min-w-0 items-center gap-2 font-semibold">
            <CalendarDays className="h-5 w-5 shrink-0" />
            <span className="truncate">Кабинет владельца</span>
          </div>
          <div className="flex shrink-0 items-center gap-4">
            {owner && (
              <div className="hidden text-right text-sm sm:block">
                <div className="font-medium">{owner.name}</div>
                <div className="text-xs text-muted-foreground">
                  {owner.email} · профиль по умолчанию
                </div>
              </div>
            )}
            <Link
              to="/"
              className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground"
            >
              <ArrowLeft className="h-4 w-4" />
              <span className="hidden sm:inline">К гостю</span>
            </Link>
          </div>
        </div>
      </header>

      <nav className="flex gap-1 border-b bg-background px-4 md:hidden">
        <MobileNavItem to="/admin" end icon={<LayoutGrid className="h-4 w-4" />}>
          Типы событий
        </MobileNavItem>
        <MobileNavItem
          to="/admin/bookings"
          icon={<ListChecks className="h-4 w-4" />}
        >
          Предстоящие встречи
        </MobileNavItem>
      </nav>

      <div className="container flex gap-6 py-8">
        <aside className="hidden w-56 shrink-0 md:block">
          <nav className="space-y-1">
            <NavItem to="/admin" end icon={<LayoutGrid className="h-4 w-4" />}>
              Типы событий
            </NavItem>
            <NavItem
              to="/admin/bookings"
              icon={<ListChecks className="h-4 w-4" />}
            >
              Предстоящие встречи
            </NavItem>
          </nav>
        </aside>
        <main className="flex-1">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

function NavItem({
  to,
  end,
  icon,
  children,
}: {
  to: string;
  end?: boolean;
  icon: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <NavLink
      to={to}
      end={end}
      className={({ isActive }) =>
        cn(
          "flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium transition-colors",
          isActive
            ? "bg-primary text-primary-foreground"
            : "text-muted-foreground hover:bg-accent hover:text-foreground",
        )
      }
    >
      {icon}
      {children}
    </NavLink>
  );
}

function MobileNavItem({
  to,
  end,
  icon,
  children,
}: {
  to: string;
  end?: boolean;
  icon: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <NavLink
      to={to}
      end={end}
      className={({ isActive }) =>
        cn(
          "flex items-center gap-2 border-b-2 px-3 py-2 text-sm font-medium transition-colors",
          isActive
            ? "border-primary text-primary"
            : "border-transparent text-muted-foreground hover:text-foreground",
        )
      }
    >
      {icon}
      {children}
    </NavLink>
  );
}
