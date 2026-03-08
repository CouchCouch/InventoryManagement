import { createRootRouteWithContext, createRoute, Link } from "@tanstack/react-router";
import { QueryClient } from '@tanstack/react-query'
import { House } from 'lucide-react'
import dragonbeaver from './assets/dragon-beaver.png'
import { Root } from "./pages/root";
import Items from "./pages/Items";


const rootRoute = createRootRouteWithContext<{
  queryClient: QueryClient
}>()({
  component: Root,
  notFoundComponent: () => {
    return (
      <>
        <div className="justify-center items-center text-center">
          <h1>No page exists here</h1>
          <img src={dragonbeaver} alt="Dragon Beaver" className="w-48 m-auto" />
          <div className='flex justify-center'>
            <Link to="/" className="p-2 hover:bg-orange-500 rounded-md">
              <House color-black />
            </Link>
          </div>
        </div>
      </>
    )
  }
})

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: function Index() {
    return (
      <h1>Index</h1>
    )
  }
})

const itemsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/items',
  component: Items
})

export const routeTree = rootRoute.addChildren([indexRoute, itemsRoute])
