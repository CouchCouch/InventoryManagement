import { useState } from "react";
import Modal from "../utilities/Modal";
import { NumberInput, TextInput } from "../utilities/Inputs";

interface AddItemModalProps {
  open: boolean,
  onClose: () => void,
  addItem: (name: string, description: string, quantity: number) => void
}

function AddItemModal({ open, onClose, addItem }: AddItemModalProps) {
  const [name, setName] = useState("")
  const [description, setDescription] = useState("")
  const [quantity, setQuantity] = useState(0)

  const resetValues = () => {
    setName("")
    setDescription("")
    setQuantity(0)
  }

  const handleClose = () => {
    addItem(name, description, quantity)
    resetValues();
    onClose();
  }

  return (
    <Modal open={open} onClose={onClose} title="Add Item">
      <>
        <TextInput label="Name" onChange={setName} value={name}/>
        <TextInput label="Description" onChange={setDescription} value={description} />
        <NumberInput label="Quantity" onChange={setQuantity} />
        <button className="btn btn-create mt-1 justify-end" onClick={() => handleClose()} value={quantity}>Add</button>
      </>
    </Modal>
  )
}


export default AddItemModal
