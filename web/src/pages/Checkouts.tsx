import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { useCheckouts } from "@/query/checkouts";
import { Copy } from "lucide-react";

export default function Checkouts() {
  const checkouts = useCheckouts().data
  console.log(checkouts)

  if (!checkouts) {
    return <div>Loading...</div>
  }

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
                <p>Notes: {checkout.notes}</p>
              </CardContent>
            </Card>
          )
        }
      </div>
    </div>
  )
}
