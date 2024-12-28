import { useState, useEffect } from "react";
import AddItemModal from "./AddItemModal";
import { TrashIcon } from "@heroicons/react/24/solid";
import Modal from "./utilities/Modal";



function ConfirmDeleteModal({ name, open, onClose, deleteFunction }) {
    return(
        <Modal title={`Delete ${name}`} open={open} onClose={onClose}>
            <p>Are you sure you want to Delete</p>
            <button className="btn btn-danger mt-1 justify-end" onClick={deleteFunction}>Confirm</button>
        </Modal>
    )
}

function Item({ id, name, description, quantity, deleteFunction }) {
    const [open, setOpen] = useState(false)

    return(
        <div className="mt-2 mb-2">
            <div className="bg-ash_gray text-slate-900 text-center p-2 rounded-xl hover:bg-ash_gray-400">
                <h1 className="text-xl font-bold mt-2 mb-4">{name}</h1>
                <p className="text-base mb-3">{description}</p>
                <h2 className="text-lg font-semibold pb-2">Quantity: {quantity}</h2>
                <button className="relative btn btn-danger" onClick={() => setOpen(true)}><TrashIcon className="size-6"/></button>
                <ConfirmDeleteModal id={id} name={name} open={open} onClose={() => setOpen(false)} deleteFunction={deleteFunction} />
            </div>
        </div>
    )
}

export default function ItemDisplay() {
    const [items, setItems] = useState(Array([]))
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)
    const [open, setOpen] = useState(false)

    function deleteItem(id) {
        console.log("deleting item: ", id)
        fetch("http://localhost:8080/items?id=" + id, {method: "DELETE"})
        .then(response => {
            if(response.ok) {
                console.log("Deleted: ", id)
            }
            else {
                console.log(response)
            }
        })
        .then(() => {
            fetchData()
        })
    }

    function fetchData() {
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
    }

    function addItem(name, description, quantity) {
        console.log("Name: ", name, ", Description: ", description, ", Quantity: ", quantity)
        if (name == "" || description == "" || quantity < 1) {
            alert("Please enter all values")
            return
        }
        let body = {
            "Name": name,
            "Description": description,
            "Quantity": parseInt(quantity)
        }
        console.log(JSON.stringify(body))
        fetch("http://localhost:8080/items", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(body)
        })
            .then(response => {
                if(response.ok) {
                    return response.json()
                }
                throw response
            })
            .then(id => {
                console.log(id)
            })
            .catch(error => {
                console.error("could not add item: ", error)
            })
            .finally(() => {
                fetchData()
            })
    }

    useEffect(() => {
        fetchData()
    }, [""]);

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
                            <Item className="w-full" key={item.Id} id={item.Id} name={item.Name} description={item.Description} quantity={item.Quantity} deleteFunction={() => deleteItem(item.Id)} />
                        )}
                    )
                }
            </div>
            <button className='btn btn-create' onClick={() => setOpen(true)}>Add Item</button>

            <AddItemModal open={open} onClose={() => setOpen(false)} addItem={addItem}/>
        </div>
    )
}
