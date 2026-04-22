import { User } from "lucide-react"
import { RouterButton } from "./ui/router-button"
import { meQueryOptions } from "@/query/user"
import { useSuspenseQuery } from "@tanstack/react-query"
import { Popover, PopoverContent, PopoverTrigger } from "./ui/popover"
import { Button } from "./ui/button"
import { useLogout } from "@/query/login"

export default function UserInfo() {
  const meQuery = useSuspenseQuery(meQueryOptions)
  const user = meQuery.data.user
  const { mutate: logout } = useLogout()

  if (user) {
    return (
      <Popover>
        <PopoverTrigger render={<Button variant="ghost" className="h-auto p-3 flex items-center" />}>
          <User className="w-1/2! h-1/2!" />
          <span className="leading-none mt-1">{user.name}</span>
        </PopoverTrigger>
        <PopoverContent align="center">
          <Button onClick={() => logout()} className="w-full" variant="destructive">
            Logout
          </Button>
        </PopoverContent>
      </Popover>
    )
  }

  return (
    <RouterButton to="/login" variant="ghost" className="h-auto p-3 flex items-center">
      <User className="w-1/2! h-1/2!" />
      <span className="leading-none mt-1">Login</span>
    </RouterButton>
  )
}
