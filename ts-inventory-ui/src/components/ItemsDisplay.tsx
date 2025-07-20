import { useEffect } from "react"
import { useState } from "react"
import { ItemT, fetchItems } from "../types/item"
import Modal from "../utilities/Modal"
import AddItemModal from "./AddItemModal"
import { PencilSquareIcon, TrashIcon } from "@heroicons/react/24/solid"
import { NumberInput, TextInput } from "../utilities/Inputs"
import { Link } from "react-router"

interface ConfirmDeleteProps {
  item: ItemT,
  open: boolean,
  onClose: () => void,
  deleteFunction: () => void
}

const ConfirmDelete = ({ item, open, onClose, deleteFunction }: ConfirmDeleteProps) => {
  return (
    <Modal title={`Delete ${item.name}`} open={open} onClose={onClose}>
      <>
        <p>Are you sure you want to Delete</p>
        <button className="btn btn-danger mt-1 justify-end" onClick={() => {deleteFunction(); onClose(); }}>Confirm</button>
      </>
    </Modal>

  )
}

interface EditModalProps {
  item: ItemT
  open: boolean,
  onClose: () => void,
  editFunction: (item: ItemT, name: string, description: string, quantity: number) => void
}

function EditModal({ item, open, onClose, editFunction}: EditModalProps) {
  const [name, setName] = useState(item.name)
  const [description, setDescription] = useState(item.description)
  const [quantity, setQuantity] = useState(item.quantity)

  return (
    <Modal open={open} onClose={onClose} title={`Edit ${name}`}>
      <>
        <TextInput label="Name" onChange={setName} value={name} />
        <TextInput label="Description" onChange={setDescription} value={description} />
        <NumberInput label="Quantity" onChange={setQuantity} value={quantity}/>
        <button className="btn btn-create mt-2 justify-end" onClick={() => {editFunction(item, name, description, quantity); onClose()}}>Save Changes</button>
      </>
    </Modal>
  )
}


interface ItemProps {
  item: ItemT,
  deleteFunction: () => void,
  editFunction: (item: ItemT, name: string, description: string, quantity: number) => void
}

function Item({ item, deleteFunction, editFunction }: ItemProps) {
  const [openDelete, setOpenDelete] = useState(false)
  const [openEdit, setOpenEdit] = useState(false)


  return(
    <div className="w-full">
      <ConfirmDelete item={item} open={openDelete} onClose={() => setOpenDelete(false)} deleteFunction={deleteFunction} />
      <EditModal item={item} open={openEdit} onClose={() => setOpenEdit(false)} editFunction={editFunction} />
      <div className="bg-l_bg1 dark:bg-bg1 text-center p-2 rounded-lg">
        <h1 className="text-xl font-bold mt-2 mb-2"><Link to={"/items/"+item.id}>{item.name}</Link></h1>
        <p className="text-base mb-2">{item.description}</p>
        <h2 className="text-lg font-semibold pb-2">Quantity: {item.quantity}</h2>
        <div className="space-x-4">
          <button className="btn btn-danger" onClick={() => setOpenDelete(true)}><TrashIcon className="size-6"/></button>
          <button className="btn btn-create" onClick={() => setOpenEdit(true)}><PencilSquareIcon className="size-6"/></button>
        </div>

      </div>
    </div>
  )
}


export default function ItemDisplay() {
  const [items, setItems] = useState<ItemT[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")
  const [open, setOpen] = useState(false)

  useEffect(() => {
    updateItems()
  }, [])

  const updateItems = () => {
    fetchItems().
      then(response => {
        setItems(response.items)
        setLoading(false)
      })
      .catch(e => {
        setError(e)
        setLoading(false)
        console.log(e)
      })
  }

  const addItem = (name: string, description: string, quantity: number): void => {
    if(name.trim().length < 1 || description.trim().length < 1 || quantity < 1) {
      alert("Please enter all values")
    }
    const body = {
      "name": name,
      "description": description,
      "quantity": quantity,
    }
    fetch(
      'http://localhost:8080/items',
      {
        method: "POST",
        headers: {"Content-Type": "application/json",},
        body: JSON.stringify(body)
      }
    )
      .then(response => {
        if (!response.ok) {
          throw response
        }
      })
      .catch(e => {
        alert(e)
        console.log(e)
      })
      .finally(() => {
        updateItems()
      })
  }

  const deleteItem = (item: ItemT): void => {
    fetch(
      "http://localhost:8080/items?id=" + item.id,
      {
        method: "DELETE"
      }
    )
      .then(response => {
        if(!response.ok) {
          throw response
        }
      })
      .catch(e => {
        alert(e)
        console.log(e)
      })
      .finally(() => {
        updateItems()
      })
  }

  const editItem = (item: ItemT, name: string, description: string, quantity: number): void => {
    if(name.trim().length < 1 || quantity < 1) {
      alert("Please enter all values")
    }
    const body = {
      "id": item.id,
      "name": name,
      "description": description,
      "quantity": quantity,
    }
    fetch(
      'http://localhost:8080/items',
      {
        method: "PUT",
        headers: {"Content-Type": "application/json",},
        body: JSON.stringify(body)
      }
    )
      .then(response => {
        if (!response.ok) {
          throw response
        }
      })
      .catch(e => {
        alert(e)
        console.log(e)
      })
      .finally(() => {
        updateItems()
      })

  }

  if(loading) {
    return (
      <>
        <img src = "../src/assets/dragon-beaver.png" className="absolute left-0 right-0 mx-auto h-24 animate-bounce" />
      </>
    )
  }

  if(error != "") {
    return (
      <Modal title="Failed to load items" open={true} onClose={():void => {location.reload()}} >
        <p>{error.toString()}</p>
        <button className="btn btn-danger" onClick={() => {location.reload()}}>Reload</button>
      </Modal>
    )
  }

  return (
    <div className="m-2">
      <div className="grid sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-2">
        {
          items.map(item => {
            return(
              <Item
                key={item.id}
                item={item}
                deleteFunction={() => {deleteItem(item)}}
                editFunction={editItem}
              />
            )}
          )
        }
      </div>
      <button className='btn btn-create mt-2' onClick={() => setOpen(true)}>Add Item</button>
      <AddItemModal open={open} onClose={() => {setOpen(false); updateItems();}} addItem={addItem}/>
    </div>
  )
}
