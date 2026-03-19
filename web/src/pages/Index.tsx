import { Button } from "@/components/ui/button";
import { ShelvingUnit, ShoppingBag, UserKey } from "lucide-react";

export default function Index() {
  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <h1 className="mb-4">Welcome!</h1>
      <div className="text-lg grid grid-cols-3 gap-2">
        <Button variant="ghost" className="h-auto p-3 flex flex-col items-center">
          <UserKey className="!w-3/4 !h-auto" />
          <span className="leading-none mt-1">Login</span>
        </Button>

        <Button variant="ghost" className="h-auto p-3 flex flex-col items-center">
          <ShelvingUnit className="!w-3/4 !h-auto" />
          <span className="leading-none mt-1">View Items</span>
        </Button>
        <Button variant="ghost" className="h-auto p-3 flex flex-col items-center">
          <ShoppingBag className="!w-3/4 !h-auto" />
          <span className="leading-none mt-1">View Checkouts</span>
        </Button>
      </div>
    </div>
  )
}
