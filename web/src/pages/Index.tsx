import { RouterButton } from "@/components/ui/router-button";
import { meQueryOptions } from "@/query/user";
import { useQuery } from "@tanstack/react-query";
import { ShelvingUnit, ShoppingBag, UserCog2 } from "lucide-react";

export default function Index() {
  const { data: user } = useQuery(meQueryOptions)

  return (
    <div className="flex flex-col items-center justify-center pt-16">
      <h1 className="mb-4">Welcome{user?.user ? ` ${user.user.name}` : ""}</h1>
      <div className={`grid ${user?.user ? "grid-cols-3" : "grid-cols-1"} gap-2`}>
        <RouterButton to="/items" variant="ghost" className="h-auto p-3 flex flex-col items-center">
          <ShelvingUnit className="w-3/4! h-auto!" />
          <span className="leading-none mt-1">{user?.user ? "Manage Items" : "View Items"}</span>
        </RouterButton>
        {user?.user &&
          <>
            <RouterButton to="/checkouts" variant="ghost" className="h-auto p-3 flex flex-col items-center">
              <ShoppingBag className="w-3/4! h-auto!" />
              <span className="leading-none mt-1">Manage Checkouts</span>
            </RouterButton>
            <RouterButton to="/users" variant="ghost" className="h-auto p-3 flex flex-col items-center">
              <UserCog2 className="w-3/4! h-auto!" />
              <span className="leading-none mt-1">Manage Users</span>
            </RouterButton>
          </>
        }
      </div>
    </div>
  )
}
