import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { ErrorPage } from "@/components/error-page";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { LoadingPage } from "@/components/loading-page";
import { useCheckouts, useReturnCheckout, type CheckoutT } from "@/query/checkouts";
import { Copy } from "lucide-react";
import { useState } from "react";

const CheckoutReturnPopup = ({ checkout }: { checkout: CheckoutT }) => {
  const [selectedItems, setSelectedItems] = useState<string[]>([])
  const [open, setOpen] = useState(false)
  const { mutate: returnCheckout, isPending, error } = useReturnCheckout();

  const toggleItem = (id: string) => {
    setSelectedItems(prev =>
      prev.includes(id) ? prev.filter(i => i !== id) : [...prev, id]
    )
  }

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    returnCheckout(
      { id: checkout.id, items: selectedItems },
      { onSuccess: () => setOpen(false) }
    )
  }

  return (
    <>
      <Dialog open={open} onOpenChange={setOpen} >
        <form onSubmit={onSubmit}>
          <DialogTrigger render={<Button variant="secondary" disabled={isPending}>Return</Button>}/>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Return</DialogTitle>
              <DialogDescription>Select Items to Return</DialogDescription>
            </DialogHeader>
            <FieldGroup>
              {
                checkout.items.map((item) =>
                  <Field key={item.item.id} orientation="horizontal">
                    <Checkbox
                      id={item.item.id}
                      name={item.item.id}
                      checked={selectedItems.includes(item.item.id)}
                      onCheckedChange={() => toggleItem(item.item.id)}
                    />
                    <FieldLabel htmlFor={item.item.id}>{item.item.name}</FieldLabel>
                  </Field>
                )
              }
            </FieldGroup>
            {error && <p className="text-destructive px-6">{(error as Error).message}</p>}
            <DialogFooter>
              <Button type="submit" disabled={isPending}>Mark as Returned</Button>
            </DialogFooter>
          </DialogContent>
        </form>
      </Dialog>
    </>
  )
}

export default function Checkouts() {
  const { data: checkouts, error, refetch, isFetching } = useCheckouts()

  if (error) return <ErrorPage error={error} refetch={refetch} />

  if (!checkouts || isFetching) return <LoadingPage />

  return (
    <div className="m-2">
      <div className="grid sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-2">
        {
          checkouts.map((checkout) =>
            <Card key={checkout.id}>
              <CardHeader>
                <CardTitle>
                  Checkout {checkout.id}
                </CardTitle>
                <CardDescription>
                  <p>Date: {(new Date(checkout.checkout_date)).toLocaleDateString()}</p>
                  <p>Personal: {checkout.personal ? "Yes" : "No"}</p>
                </CardDescription>
              </CardHeader>
              <CardContent>
                <p>Name: {checkout.user.name}</p>
                <div className="inline-flex items-center">
                  <p>Email: {checkout.user.email}</p>
                  <Button className="h-auto p-0.5 leading-none" variant="ghost" onClick={() => navigator.clipboard.writeText(checkout.user.email)}><Copy /></Button>
                </div>
                <p>Eboard: {checkout.created_by.name}</p>
                {checkout.notes && <p>Notes: {checkout.notes}</p>}
                <p>Items:</p>
                <ul className="pl-1">
                  {checkout.items.map((item) =>
                    <li key={item.item.id}>{item.item.name}</li>
                  )}
                </ul>
              </CardContent>
              <CardFooter>
                <CheckoutReturnPopup checkout={checkout} />
              </CardFooter>
            </Card>
          )
        }
      </div>
    </div>
  )
}
