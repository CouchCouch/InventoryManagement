import { useEffect, useState } from "react"
import { CheckoutT, fetchCheckouts } from "../types/checkout"

export const CheckouotDashboard = () => {
  const [checkouts, setCheckouts] = useState<CheckoutT[]>()
  const [error, setError] = useState("")

  useEffect(() => {
    fetchCheckouts()
      .then(response => {
        setCheckouts(response.checkouts)
      })
      .catch(e => {
        setError(e)
        console.log(e)
      })
  }, [])

  if (error) {
    return <div>Error: {error}</div>
  }

  return (
    <>

    </>
  )
}
