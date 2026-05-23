import { createRootRouteWithContext, createRoute, Link } from "@tanstack/react-router";
import { QueryClient } from '@tanstack/react-query'
import { House } from 'lucide-react'
import dragonbeaver from './assets/dragon-beaver.png'
import { Root } from "./pages/root";
import Items from "./pages/Items";
import Login from "./pages/Login";
import Index from "./pages/Index";
import Checkouts from "./pages/Checkouts";


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
  component: Index
})

const itemsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/items',
  component: Items
})

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  component: Login
})

const checkoutsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/checkouts',
  component: Checkouts
})

export const routeTree = rootRoute.addChildren([indexRoute, itemsRoute, loginRoute, checkoutsRoute])
