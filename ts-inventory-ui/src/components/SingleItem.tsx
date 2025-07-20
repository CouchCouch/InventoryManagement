import { useEffect } from "react"
import { useState } from "react"
import { fetchItem, ItemT } from "../types/item"
//import { useParams } from "react-router"

export default function Item({ id }: { id: number }) {
  //const params = useParams()
  //const id = Number(params.id)
  const [item, setItem] = useState<ItemT>()
  const [error, setError] = useState("")

  useEffect(() => {
    fetchItem(id)
      .then(response => {
        setItem(response.items[0])
      })
      .catch(e => {
        setError(e)
        console.log(e)
      })
  }, [id])

  if(error) {
    return (
      <div>
        <p>{error.toString()}</p>
        <button className="btn btn-danger" onClick={() => {location.reload()}}>Reload</button>
      </div>
    )
  }

  if(!item) {
    return (
      <div>
        <p>No Such Item</p>
        <button className="btn btn-danger" onClick={() => {location.reload()}}>Reload</button>
      </div>
    )
  }
  return(
    <div className="mt-2 mb-2 w-full">
      <div className="text-start">
        <h1 className="text-xl font-bold mt-2 mb-4">{item.name}</h1>
        <p className="mb-3"><span className="font-semibold">Descriptors:</span> {item.description}</p>
        <p className="pb-2"><span className="font-semibold">Quantity:</span> {item.quantity}</p>
      </div>
    </div>
  )
}
