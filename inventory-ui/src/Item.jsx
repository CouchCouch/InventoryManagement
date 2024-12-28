import { useState, useEffect } from "react";
import AddItemModal from "./AddItemModal";

function Item({ name, description, quantity }) {
    return(
        <div className="mt-2 mb-2">
            <div className="bg-ash_gray text-slate-900 text-center p-2 rounded-xl hover:bg-ash_gray-400">
                <h1 className="text-xl font-bold mt-2 mb-4">{name}</h1>
                <p className="text-base mb-3">{description}</p>
                <h2 className="text-lg font-semibold pb-2">Quantity: {quantity}</h2>
            </div>
        </div>
    )
}

export default function ItemDisplay() {
    const [items, setItems] = useState(Array([]))
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)
    const [open, setOpen] = useState(false)

    useEffect(() => {
        fetch("http://localhost:8080/items")
            .then(response => {
                if(response.ok) {
                    return response.json()
                }
                throw response
            })
            .then(data => {
                console.log(data)
                setItems(data['Items'])
                console.log(items)
            })
            .catch(error => {
                console.error("error fetching data: ", error)
                setError(error)
            })
            .finally(() => {
                setLoading(false)
            })
    }, []);

    if (loading) return(
        <div className="text-center">
            <h1 className="text-4xl font-bold bg-red-700 align-middle pt-10 pb-10">loading...</h1>
        </div>
    )
    if (error) return (
        <div className="text-4xl font-bold bg-red-700 align-middle text-center pt-10 pb-10">
            <h1>Error! please try reloading</h1>
        </div>
    )

    if (items.length < 1) {
        return(
        <div className="text-4xl font-bold bg-red-700 align-middle text-center pt-10 pb-10">
            <h1>No Data</h1>
        </div>
        )
    }

    return(
        <div className="m-2">
            <div className="grid grid-cols-4 gap-4">
                {
                    items.map(item => {
                        return(
                            <Item className="w-full" key={item.Id} name={item.Name} description={item.Description} quantity={item.Quantity} />
                        )}
                    )
                }
            </div>
            <button className='btn btn-create' onClick={() => setOpen(true)}>
              Show Modal
            </button>

            <AddItemModal open={open} onClose={() => setOpen(false)} />
        </div>
    )
}
