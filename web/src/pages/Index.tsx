import { RouterButton } from "@/components/ui/router-button";
import { meQueryOptions } from "@/query/user";
import { useSuspenseQuery } from "@tanstack/react-query";
import { ShelvingUnit, ShoppingBag } from "lucide-react";

export default function Index() {
  const user = useSuspenseQuery(meQueryOptions).data.user

  return (
    <div className="flex flex-col items-center justify-center pt-16">
      <h1 className="mb-4">Welcome{user ? ` ${user.name}` : ""}</h1>
      <div className={`grid ${user ? "grid-cols-2" : "grid-cols-1"} gap-2`}>
        <RouterButton to="/items" variant="ghost" className="h-auto p-3 flex flex-col items-center">
          <ShelvingUnit className="w-3/4! h-auto!" />
          <span className="leading-none mt-1">View Items</span>
        </RouterButton>
        {user &&
          <RouterButton to="/checkouts" variant="ghost" className="h-auto p-3 flex flex-col items-center">
            <ShoppingBag className="w-3/4! h-auto!" />
            <span className="leading-none mt-1">View Checkouts</span>
          </RouterButton>
        }
      </div>
    </div>
  )
}
