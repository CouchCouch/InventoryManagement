import { useState } from "react";
import Modal from "./utilities/Modal";
import { NumberInput, TextInput } from "./utilities/Inputs";

function AddItemModal({ open, onClose }) {
    const [name, setName] = useState("")
    const [description, setDescription] = useState("")
    const [quantity, setQuantity] = useState(0)

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
                setName("")
                setDescription("")
                setQuantity(0)
                onClose()
            })
    }

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