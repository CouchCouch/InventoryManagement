import { Plus } from "lucide-react";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "./ui/dialog";
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";
import { Button } from "./ui/button";
import { Field, FieldGroup, FieldLabel } from "./ui/field";
import { Input } from "./ui/input";
import { useState } from "react";
import { useCreateItem, useTypes } from "@/query/items";
import { Combobox, ComboboxContent, ComboboxEmpty, ComboboxInput, ComboboxItem, ComboboxList } from "./ui/combobox";
import { Textarea } from "./ui/textarea";

export default function NewItem() {
  const [id, setId] = useState('')
  const [idError, setIdError] = useState(false)
  const [name, setName] = useState('')
  const [nameError, setNameError] = useState(false)
  const [type, setType] = useState('')
  const [typeError, setTypeError] = useState(false)
  const [datePurchased, setDatePurchased] = useState<Date | null>(null)
  const [notes, setNotes] = useState<string | null>(null)
  const { mutate: addItem, isPending } = useCreateItem();

  const types = useTypes()

  const validateForm = (): boolean => {
    console.log("validating form with values:", {id, name, type})
    var valid = true
    setIdError(false);
    setNameError(false);
    setTypeError(false);

    if (!id || !/^[A-Za-z]{2}-\d{2}-\d{2}$/.test(id)) {
      setIdError(true)
      valid = false
    }
    if (!name || name.trim() === '') {
      setNameError(true)
      valid = false
    }
    if (!type || type.trim() === '') {
      setTypeError(true)
      valid = false
    }
    return valid
  }

  const onSubmit = (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault()
    console.log("form submitted")
    if (!validateForm()) return
    addItem({id, name, type, notes, date_purchased: datePurchased})
  }

  const resetState = () => {
    setId('')
    setName('')
    setType('')
    setDatePurchased(null)
    setNotes(null)
  }

  return (
    <Dialog onOpenChange={() => resetState()}>
      <Tooltip>
        <TooltipTrigger render={
          <DialogTrigger render={
            <Button size="icon" className="rounded-full absolute bottom-8 right-2 h-20 w-20" onClick={() => console.log("clicked")} />
          }/>
        }>
          <Plus className="w-16 h-16"/>
        </TooltipTrigger>
        <TooltipContent>
          <p>Add a new item</p>
        </TooltipContent>
      </Tooltip>
      <DialogContent className="md:max-w-md sm:max-w-sm">
        <form onSubmit={(e) => onSubmit(e)}>
          <DialogHeader>
            <DialogTitle>Add Item</DialogTitle>
            <DialogDescription>
              Add items with the given info
            </DialogDescription>
          </DialogHeader>
          <FieldGroup className="py-4">
            <Field>
              <FieldLabel htmlFor="id">ID<span className="text-destructive">*</span></FieldLabel>
              <Input
                id="id"
                name="id"
                onChange={(e) => {setId(e.target.value)}}
                placeholder="XX-00-00"
                aria-invalid={idError}
                required
              />
            </Field>
            <Field>
              <FieldLabel htmlFor="name">Name<span className="text-destructive">*</span></FieldLabel>
              <Input
                id="name"
                name="name"
                onChange={(e) => {setName(e.target.value)}}
                placeholder="Pickle"
                aria-invalid={nameError}
                required
              />
            </Field>
            <Field>
              <FieldLabel htmlFor="type">Type<span className="text-destructive">*</span></FieldLabel>
              <Combobox
                id="type"
                name="type"
                items={types.data}
                onValueChange={(val) => {
                  setType(val as string);
                  setTypeError(false);
                }}
                aria-invalid={typeError}
                required
              >
                <ComboboxInput placeholder="Select a type or enter a new one" />
                <ComboboxContent>
                  <ComboboxEmpty>No Types Found.</ComboboxEmpty>
                  <ComboboxList>
                    {(item: string) => (
                      <ComboboxItem key={item} value={item}>
                        {item}
                      </ComboboxItem>
                    )}
                  </ComboboxList>
                </ComboboxContent>
              </Combobox>
            </Field>
            <Field>
              <FieldLabel htmlFor="date_purchased">Date Purchased</FieldLabel>
              <Input id="date_purchased" name="date_purchased" type="date" onChange={(e) => {setDatePurchased(new Date(e.target.value))}} />
            </Field>
            <Field>
              <FieldLabel htmlFor="notes">Notes</FieldLabel>
              <Textarea id="notes" name="notes" onChange={(e) => {setNotes(e.target.value)}} className="resize-none h-20"/>
            </Field>
          </FieldGroup>
          <DialogFooter>
            <DialogClose render={<Button variant="outline" />} onClick={() => resetState()}>
              Cancel
            </DialogClose>
            <Button type="submit" disabled={isPending}>Add Item</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
