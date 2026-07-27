import { Link, Outlet, useLocation } from "@tanstack/react-router"
import {
  LayoutDashboard,
  Database,
  Search,
  BarChart3,
  Settings,
  BookOpen,
  Play,
} from "lucide-react"
import { cn } from "../lib/utils"

const navItems = [
  { icon: LayoutDashboard, label: "Dashboard", path: "/" },
  { icon: Play, label: "Playground", path: "/playground" },
  { icon: Database, label: "Collections", path: "/collections" },
  { icon: Search, label: "Search", path: "/search" },
  { icon: BarChart3, label: "Performance", path: "/performance" },
  { icon: BookOpen, label: "Docs", path: "/docs" },
  { icon: Settings, label: "Settings", path: "/settings" },
]

export default function Layout() {
  const location = useLocation()

  return (
    <div className="flex h-screen bg-zinc-50 dark:bg-zinc-950">
      {/* Sidebar */}
      <aside className="w-56 border-r border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 flex flex-col">
        <div className="p-5 border-b border-zinc-200 dark:border-zinc-800">
          <h1 className="text-lg font-bold tracking-tight text-zinc-900 dark:text-zinc-50">
            Oris
          </h1>
          <p className="text-xs text-zinc-500 mt-0.5">
            Vector Retrieval Engine
          </p>
        </div>
        <nav className="flex-1 p-3 space-y-1">
          {navItems.map((item) => {
            const active = location.pathname === item.path
            return (
              <Link
                key={item.path}
                to={item.path}
                className={cn(
                  "flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors",
                  active
                    ? "bg-zinc-100 text-zinc-900 dark:bg-zinc-800 dark:text-zinc-50"
                    : "text-zinc-600 hover:bg-zinc-100 dark:text-zinc-400 dark:hover:bg-zinc-800 dark:hover:text-zinc-50",
                )}
              >
                <item.icon className="h-4 w-4" />
                {item.label}
              </Link>
            )
          })}
        </nav>
        <div className="p-4 border-t border-zinc-200 dark:border-zinc-800">
          <p className="text-xs text-zinc-400">v0.1.0</p>
        </div>
      </aside>

      {/* Main content */}
      <main className="flex-1 overflow-auto">
        <div className="p-6 max-w-6xl mx-auto">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
