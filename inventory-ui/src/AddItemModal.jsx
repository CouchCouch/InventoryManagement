import { useState } from "react";
import Modal from "./utilities/Modal";
import { NumberInput, TextInput } from "./utilities/Inputs";

function AddItemModal({ open, onClose, addItem}) {
    const [name, setName] = useState("")
    const [description, setDescription] = useState("")
    const [quantity, setQuantity] = useState(0)

    return (
        <Modal open={open} onClose={onClose} title="Add Item">
            <TextInput label="Name" onChange={setName} />
            <TextInput label="Description" onChange={setDescription} />
            <NumberInput label="Quantity" onChange={setQuantity} />
            <button className="btn btn-create mt-1 justify-end" onClick={() => addItem(name, description, quantity)}>Create</button>
        </Modal>
    )
}

export default AddItemModal