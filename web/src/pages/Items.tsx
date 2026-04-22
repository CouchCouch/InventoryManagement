import NewItem from "@/components/new-item";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useDeleteItem, useItems, useTypes } from "@/query/items";
import { meQueryOptions } from "@/query/user";
import { useSuspenseQuery } from "@tanstack/react-query";
import { ShoppingCart, Trash2 } from "lucide-react";
import { useState } from "react";

export default function Items() {
  const [typeFilter, setTypeFilter] = useState<string>('none')

  const items = useItems(typeFilter).data

  const user = useSuspenseQuery(meQueryOptions).data.user
  const types = useTypes().data || []

  const {mutate: deleteItem, isPending} = useDeleteItem()

  if (!items) {
    return <div>Loading...</div>
  }

  return (
    <div className="m-2">
      <div className="m-2">
        <Select onValueChange={(value: string | null) => setTypeFilter(value || 'none')}>
          <SelectTrigger className="w-full max-w-48">
            <SelectValue />
          </SelectTrigger>
          <SelectContent alignItemWithTrigger={false} >
            <SelectItem key="none" value="Select Type">Select Type</SelectItem>
            {types?.map((type) =>
              <SelectItem key={type} value={type}>
                {type}
              </SelectItem>
            )}
          </SelectContent>
        </Select>
      </div>
      <div className="grid sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-2">
        {
          items.map((item) =>
            <Card key={item.id}>
              <CardHeader>
                <CardTitle>
                  {item.name}
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p>Type: {item.type}</p>
                <p>Notes: {item.notes}</p>
              </CardContent>
              <CardFooter className={`justify-end grid ${user ? "grid-cols-2" : "grid-cols-1"} gap-2`}>
                {
                  user &&
                    <Button variant="destructive" onClick={() => deleteItem(item.id)} disabled={isPending}>
                      <Trash2 />
                    </Button>
                }
                <Button variant="secondary" >
                  <ShoppingCart />
                </Button>
              </CardFooter>
            </Card>
          )
        }
      </div>
      {user && <NewItem />}
    </div>
  )
}
