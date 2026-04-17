import { User } from "lucide-react"
import { RouterButton } from "./ui/router-button"
import { meQueryOptions, type UserT } from "@/query/user"
import { useState } from "react"
import { useSuspenseQuery } from "@tanstack/react-query"

export default function UserInfo() {
  const meQuery = useSuspenseQuery(meQueryOptions)
  const user = meQuery.data.user as UserT

  return (
    <RouterButton to="/login" variant="ghost" className="h-auto p-3 flex items-center">
      <User className="w-1/2! h-1/2!" />
      <span className="leading-none mt-1">{user?.name || "Login"}</span>
    </RouterButton>
  )
}
