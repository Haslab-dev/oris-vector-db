import { createHashHistory, createRouter, createRoute, createRootRoute } from "@tanstack/react-router"
import Layout from "./components/layout"
import DashboardPage from "./pages/dashboard"
import CollectionsPage from "./pages/collections"
import CollectionDetailPage from "./pages/collection-detail"
import SearchPage from "./pages/search"
import PerformancePage from "./pages/performance"
import SettingsPage from "./pages/settings"

const rootRoute = createRootRoute({
  component: Layout,
})

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: DashboardPage,
})

const collectionsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/collections",
  component: CollectionsPage,
})

const collectionDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/collections/$name",
  component: CollectionDetailPage,
})

const searchRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/search",
  component: SearchPage,
})

const performanceRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/performance",
  component: PerformancePage,
})

const settingsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/settings",
  component: SettingsPage,
})

const routeTree = rootRoute.addChildren([
  indexRoute,
  collectionsRoute,
  collectionDetailRoute,
  searchRoute,
  performanceRoute,
  settingsRoute,
])

const history = createHashHistory()

export const router = createRouter({ routeTree, history })

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router
  }
}
