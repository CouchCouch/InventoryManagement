import { useEffect, useState } from "react"
import { fetchHistory, CheckoutT } from "../types/checkout"

const HistoryItem = ({ checkout }: { checkout: CheckoutT }) => {
    return(
        <div className="pb-2">
            <p><span className="font-bold">Date:</span> {new Date(checkout.date).toLocaleDateString()}</p>
            <h4>Borower Info:</h4>
            <div className="px-2">
                <p>Name: {checkout.name}</p>
                <p>Name: {checkout.email}</p>
            </div>
            <p><span className="font-bold">Status:</span> {checkout.returned ? "Returned" : "Not Returned"}</p>
        </div>
    )
}

export const CheckoutHistory = ({ id }: { id: number }) => {
    const [history, setHistory] = useState<CheckoutT[]>();
    const [error, setError] = useState("")

    useEffect(() => {
        fetchHistory(id)
        .then(response => {
                setHistory(response.checkouts)
            })
            .catch(e => {
                setError(e)
                console.log(e)
            })
    }, [])

    if(error) {
        return (
            <div>
                <h2>History:</h2>
                <p>{error.toString()}</p>
            </div>
        )
    }

    if(!history) {
        return (
            <div>
                <h2>History:</h2>
                <p>No History to Display</p>
            </div>
        )
    }
    return(
        <div className="mt-2 mb-2 w-full">
            <h2>History:</h2>
            <div className="text-start p-2">
                {
                    history.map(checkout => {
                        return(
                            <HistoryItem
                                key={checkout.id}
                                checkout={checkout}
                            />

                        )
                    })
                }
            </div>
        </div>
    )
}
