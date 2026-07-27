import { createHashHistory, createRouter, createRoute, createRootRoute } from "@tanstack/react-router"
import Layout from "./components/layout"
import DashboardPage from "./pages/dashboard"
import CollectionsPage from "./pages/collections"
import CollectionDetailPage from "./pages/collection-detail"
import SearchPage from "./pages/search"
import PerformancePage from "./pages/performance"
import SettingsPage from "./pages/settings"
import PlaygroundPage from "./pages/playground"
import DocsOverviewPage from "./pages/docs/overview"
import GettingStartedPage from "./pages/docs/getting-started"
import ArchitecturePage from "./pages/docs/architecture"
import DenseEnginePage from "./pages/docs/dense-engine"
import SparseEnginePage from "./pages/docs/sparse-engine"
import MetadataEnginePage from "./pages/docs/metadata-engine"
import StorageEnginePage from "./pages/docs/storage-engine"
import SegmentManagerPage from "./pages/docs/segment-manager"
import QueryPlannerPage from "./pages/docs/query-planner"
import ApiReferencePage from "./pages/docs/api-reference"

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

const docsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/docs",
  component: DocsOverviewPage,
})

const playgroundRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/playground",
  component: PlaygroundPage,
})

const docsGettingStartedRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/docs/getting-started",
  component: GettingStartedPage,
})

const docsArchitectureRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/docs/architecture",
  component: ArchitecturePage,
})

const docsDenseEngineRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/docs/dense-engine",
  component: DenseEnginePage,
})

const docsSparseEngineRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/docs/sparse-engine",
  component: SparseEnginePage,
})

const docsMetadataEngineRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/docs/metadata-engine",
  component: MetadataEnginePage,
})

const docsStorageEngineRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/docs/storage-engine",
  component: StorageEnginePage,
})

const docsSegmentManagerRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/docs/segment-manager",
  component: SegmentManagerPage,
})

const docsQueryPlannerRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/docs/query-planner",
  component: QueryPlannerPage,
})

const docsApiReferenceRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/docs/api-reference",
  component: ApiReferencePage,
})

const routeTree = rootRoute.addChildren([
  indexRoute,
  collectionsRoute,
  collectionDetailRoute,
  searchRoute,
  performanceRoute,
  settingsRoute,
  docsRoute,
  playgroundRoute,
  docsGettingStartedRoute,
  docsArchitectureRoute,
  docsDenseEngineRoute,
  docsSparseEngineRoute,
  docsMetadataEngineRoute,
  docsStorageEngineRoute,
  docsSegmentManagerRoute,
  docsQueryPlannerRoute,
  docsApiReferenceRoute,
])

const history = createHashHistory()

export const router = createRouter({ routeTree, history })

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router
  }
}
