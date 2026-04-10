import { ReactQueryDevtools } from "@tanstack/react-query-devtools"
import { Link, Outlet } from "@tanstack/react-router"
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools"
import { ThemeProvider } from "@/components/theme-provider"
import { ModeToggle } from "@/components/mode-toggle"
import { TooltipProvider } from "@/components/ui/tooltip"
import UserInfo from "@/components/user"

export const Root = () => {
  return (
    <TooltipProvider>
      <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
        <header>
          <Link to="/"><h1>RITOC Inventory</h1></Link>
          <div className="absolute left-2" >
            <ModeToggle />
          </div>
          <div className="absolute right-2" >
            <UserInfo />
          </div>
        </header>
        <Outlet />
        <ReactQueryDevtools buttonPosition="bottom-right" />
        <TanStackRouterDevtools position="bottom-right" />
      </ThemeProvider>
    </TooltipProvider>
  )
}
