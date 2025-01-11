import { useState, useEffect } from "react";
import AddItemModal from "./AddItemModal";
import { PencilSquareIcon, TrashIcon } from "@heroicons/react/24/solid";
import Modal from "./utilities/Modal";
import { NumberInput, TextInput } from "./utilities/Inputs";
import PropTypes from 'prop-types';



function ConfirmDeleteModal({ name, open, onClose, deleteFunction }) {
    return(
        <Modal title={`Delete ${name}`} open={open} onClose={onClose}>
            <>
                <p>Are you sure you want to Delete</p>
                <button className="btn btn-danger mt-1 justify-end" onClick={deleteFunction}>Confirm</button>
            </>
        </Modal>
    )
}

ConfirmDeleteModal.propTypes = {
    name: PropTypes.string,
    open: PropTypes.bool,
    onClose: PropTypes.func,
    deleteFunction: PropTypes.func
}

function EditModal({ id, name, description, quantity, open, onClose, editItem}) {
    const [newName, setName] = useState(name)
    const [newDescription, setDescription] = useState(description)
    const [newQuantity, setQuantity] = useState(quantity)

    return (
        <Modal open={open} onClose={onClose} title={`Edit ${name}`}>
            <>
                <TextInput label="Name" onChange={setName} value={newName} />
                <TextInput label="Description" onChange={setDescription} value={newDescription} />
                <NumberInput label="Quantity" onChange={setQuantity} value={newQuantity}/>
                <button className="btn btn-create mt-1 justify-end" onClick={() => {editItem(id, newName, newDescription, newQuantity); onClose()}}>Edit</button>
            </>
        </Modal>
    )
}

EditModal.propTypes = {
    id: PropTypes.number,
    name: PropTypes.string,
    description: PropTypes.string,
    quantity: PropTypes.number,
    open: PropTypes.bool,
    onClose: PropTypes.func,
    editItem: PropTypes.func
}

function Item({ id, name, description, quantity, deleteFunction, updateFunction }) {
    const [openDelete, setOpenDelete] = useState(false)
    const [openEdit, setOpenEdit] = useState(false)

    return(
        <div className="mt-2 mb-2">
            <div className="bg-ash_gray text-slate-900 text-center p-2 rounded-xl">
                <h1 className="text-xl font-bold mt-2 mb-4">{name}</h1>
                <p className="text-base mb-3">{description}</p>
                <ConfirmDeleteModal name={name} open={openDelete} onClose={() => setOpenDelete(false)} deleteFunction={deleteFunction} />
                <EditModal id={id} name={name} description={description} quantity={quantity} editItem={updateFunction} open={openEdit} onClose={() => setOpenEdit(false)}/>
                <h2 className="text-lg font-semibold pb-2">Quantity: {quantity}</h2>
                <div className="space-x-4">
                    <button className="btn btn-danger" onClick={() => setOpenDelete(true)}><TrashIcon className="size-6"/></button>
                    <button className="btn btn-edit" onClick={() => setOpenEdit(true)}><PencilSquareIcon className="size-6"/></button>
                </div>
            </div>
        </div>
    )
}

Item.propTypes = {
    id: PropTypes.number,
    name: PropTypes.string,
    description: PropTypes.string,
    quantity: PropTypes.number,
    deleteFunction: PropTypes.func,
    updateFunction: PropTypes.func
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
        console.log("CREATING: Name: ", name, ", Description: ", description, ", Quantity: ", quantity)
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

    function updateItem(id, name, description, quantity) {
        console.log("UPDATING: Id: ", id, "Name: ", name, ", Description: ", description, ", Quantity: ", quantity)
        if (name == "" || description == "" || quantity < 1) {
            alert("Please enter all values")
            return
        }
        let body = {
            "Id": id,
            "Name": name,
            "Description": description,
            "Quantity": parseInt(quantity)
        }
        console.log(JSON.stringify(body))
        fetch("http://localhost:8080/items", {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(body)
        })
            .then(response => {
                if(response.ok) {
                    console.log("Updated: ", id)
                }
                else {
                    console.log(response)
                }
            })
            .then(() => {
                fetchData()
            })
    }

    useEffect(() => {
        fetchData()
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

    if (!items) {
        return(
            <div className='m-2'>
                <div className="text-4xl font-bold align-middle text-center pt-10 pb-10">
                    <h1>No Items</h1>
                </div>
                <button className='btn btn-create' onClick={() => setOpen(true)}>Add Item</button>
                <AddItemModal open={open} onClose={() => setOpen(false)} addItem={addItem}/>
            </div>
        )
    }

    return(
        <div className="m-2">
            <div className="grid sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
                {
                    items.map(item => {
                        return(
                            <Item className="w-full" key={item.Id} id={item.Id} name={item.Name} description={item.Description} quantity={item.Quantity} deleteFunction={() => deleteItem(item.Id)} updateFunction={updateItem}/>
                        )}
                    )
                }
            </div>
            <button className='btn btn-create' onClick={() => setOpen(true)}>Add Item</button>

            <AddItemModal open={open} onClose={() => setOpen(false)} addItem={addItem}/>
        </div>
    )
}
