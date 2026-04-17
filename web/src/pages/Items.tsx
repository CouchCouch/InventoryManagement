import NewItem from "@/components/new-item";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { useDeleteItem, useItems } from "@/query/items";
import { ShoppingCart, Trash2 } from "lucide-react";

export default function Items() {
  const items = useItems().data || []

  const {mutate: deleteItem, isPending} = useDeleteItem()

  return (
    <div className="m-2">
      <div className="grid sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-2">
        {
          items.map(item => {
            return(
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
                <CardFooter className="justify-end grid grid-cols-2 gap-2">
                  <Button variant="destructive" onClick={() => deleteItem(item.id)} disabled={isPending}>
                    <Trash2 />
                  </Button>
                  <Button variant="secondary" >
                    <ShoppingCart />
                  </Button>
                </CardFooter>
              </Card>
            )
          })
        }
      </div>
      <NewItem />
    </div>
  )
}
