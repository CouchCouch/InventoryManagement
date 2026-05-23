import NewItem from "@/components/new-item";
import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useDeleteItem, useItems, useTypes } from "@/query/items";
import { meQueryOptions } from "@/query/user";
import { useSuspenseQuery } from "@tanstack/react-query";
import { ShoppingCart, Trash2 } from "lucide-react";
import { useState } from "react";

type FilterProps = {
  value: string,
  onValueChange: (value: string | null) => void
  filterItems: string[],
  defaultValue: string
}

const Filter = ({ value, onValueChange, filterItems, defaultValue }: FilterProps ) => {
  return(
    <Select value={value} onValueChange={onValueChange}>
      <SelectTrigger className="w-full max-w-48">
        <SelectValue />
      </SelectTrigger>
      <SelectContent alignItemWithTrigger={false} >
        <SelectItem key="none" value={defaultValue}>{defaultValue}</SelectItem>
        {filterItems?.map((item) =>
          <SelectItem key={item} value={item}>
            {item}
          </SelectItem>
        )}
      </SelectContent>
    </Select>
  )
}

export default function Items() {
  const [typeFilter, setTypeFilter] = useState<string>('Select Type')

  const items = useItems(typeFilter).data

  const user = useSuspenseQuery(meQueryOptions).data.user
  const types = useTypes().data || []

  const {mutate: deleteItem, isPending} = useDeleteItem()

  if (!items) {
    return <div>Loading...</div>
  }

  return (
    <div className="m-2">
      <div className="my-4 grid-cols-5">
        <Filter
          value={typeFilter}
          defaultValue="Select Type"
          onValueChange={(value: string | null) => setTypeFilter(value || "Select Type")}
          filterItems={types}
        />
      </div>
      <div className="grid sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-2">
        {
          items.map((item) =>
            <Card key={item.id}>
              <CardHeader>
                <CardTitle>
                  {item.name}
                </CardTitle>
                <CardDescription>
                  <p>ID: {item.id.toUpperCase()}</p>
                  <p>Type: {item.type}</p>
                  <p>Notes: {item.notes}</p>
                </CardDescription>
              </CardHeader>
              {
                user &&
                  <CardFooter className={`justify-end grid ${user ? "grid-cols-2" : "grid-cols-1"} gap-2`}>
                    <Button variant="destructive" onClick={() => deleteItem(item.id)} disabled={isPending}>
                      <Trash2 />
                    </Button>
                    <Button variant="secondary" >
                      <ShoppingCart />
                    </Button>
                  </CardFooter>
              }
            </Card>
          )
        }
      </div>
      {user && <NewItem />}
    </div>
  )
}
