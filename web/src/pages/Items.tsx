import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Combobox, ComboboxContent, ComboboxEmpty, ComboboxInput, ComboboxItem, ComboboxList } from "@/components/ui/combobox";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { ErrorPage } from "@/components/error-page";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { LoadingPage } from "@/components/loading-page";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useCreateItem, useDeleteItem, useItems, useTypes, type ItemT } from "@/query/items";
import { meQueryOptions, type UserT } from "@/query/user";
import { useQuery, type UseMutateFunction } from "@tanstack/react-query";
import { Image, Plus, ShoppingCart, Trash2 } from "lucide-react";
import { useState } from "react";

const NewItem = () => {
  const [open, setOpen] = useState(false)
  const [id, setId] = useState('')
  const [idError, setIdError] = useState(false)
  const [name, setName] = useState('')
  const [nameError, setNameError] = useState(false)
  const [type, setType] = useState('')
  const [typeError, setTypeError] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  const [datePurchased, setDatePurchased] = useState<Date | null>(null)
  const [notes, setNotes] = useState('')
  const { mutate: addItem, isPending, error } = useCreateItem();

  const types = useTypes()
  const displayedTypes = Array.isArray(types.data) ? [...types.data] : []
  if (searchQuery && !displayedTypes.includes(searchQuery)) {
    displayedTypes.push(searchQuery)
  }


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
    console.log(datePurchased)
    addItem(
      { id, name, type, notes, date_purchased: datePurchased },
      {
        onSuccess: () => {
          resetState()
          setOpen(false)
        }
      }
    )

  }

  const resetState = () => {
    setId('')
    setName('')
    setType('')
    setDatePurchased(null)
    setNotes('')
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <Tooltip>
        <TooltipTrigger render={
          <DialogTrigger render={
            <Button size="icon" className="rounded-full fixed bottom-8 right-2 h-20 w-20 shadow-xl" onClick={() => console.log("clicked")} />
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
                value={id}
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
                value={name}
                required
              />
            </Field>
            <Field>
              <FieldLabel htmlFor="type">Type<span className="text-destructive">*</span></FieldLabel>
              <Combobox
                id="type"
                name="type"
                items={displayedTypes}
                onValueChange={(val) => {
                  setType(val as string);
                  setTypeError(false);
                }}
                onInputValueChange={(val) => {
                  setSearchQuery(val);
                  setType(val); // Back-up: captures raw text in case they submit without clicking
                  setTypeError(false);
                }}
                value={type}
                aria-invalid={typeError}
                required
              >
                <ComboboxInput
                  placeholder="Select a type or enter a new one"
                  onChange={(e) => {
                    const val = e.target.value;
                    setSearchQuery(val);
                    setType(val);
                    setTypeError(false);
                  }}
                />
                <ComboboxContent>
                  <ComboboxEmpty>
                    {searchQuery || "No Types Found."}
                  </ComboboxEmpty>
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
              <Input
                id="date_purchased"
                name="date_purchased"
                type="date"
                onChange={(e) => {setDatePurchased(new Date(e.target.value))}}
                value={
                  datePurchased
                    ? `${datePurchased.getFullYear()}-${String(datePurchased.getMonth() + 1).padStart(2, '0')}-${String(datePurchased.getDate()).padStart(2, '0')}`
                    : ""
                }
              />
            </Field>
            <Field>
              <FieldLabel htmlFor="notes">Notes</FieldLabel>
              <Textarea
                id="notes"
                name="notes"
                onChange={(e) => {setNotes(e.target.value)}}
                className="resize-none h-20"
                value={notes}
              />
            </Field>
          </FieldGroup>
          <DialogFooter>
            {error && <p>{error?.name}</p>}
            <Button variant="destructive" disabled={isPending} onClick={() => resetState()}>Clear</Button>
            <DialogClose render={<Button variant="outline" />}>
              Cancel
            </DialogClose>
            <Button type="submit" disabled={isPending}>Add Item</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

type FilterProps = {
  value: string,
  onValueChange: (value: string | null) => void
  filterItems?: string[],
  defaultValue: string
}

const Filter = ({ value, onValueChange, filterItems, defaultValue }: FilterProps ) => {
  return(
    <Select value={value} onValueChange={onValueChange}>
      <SelectTrigger className="w-full max-w-48">
        <SelectValue />
      </SelectTrigger>
      <SelectContent alignItemWithTrigger={false} >
        <SelectItem key="none" value={defaultValue}>{defaultValue}</SelectItem>
        {filterItems?.map((item) =>
          <SelectItem key={item} value={item}>
            {item}
          </SelectItem>
        )}
      </SelectContent>
    </Select>
  )
}

type ItemProps = {
  item: ItemT
  user: UserT | undefined
  deleteFunction: UseMutateFunction<void, Error, string, unknown>
}

const Item = ({item, user, deleteFunction}: ItemProps) =>
  <Card key={item.id} className="pt-0">
    <div className="bg-muted flex items-center justify-center h-40 rounded-t-lg">
      <Image className="w-16 h-16 text-muted-foreground" />
    </div>
    <CardHeader>
      <CardTitle>
        {item.name}
      </CardTitle>
      <CardDescription>
        <p>ID: {item.id.toUpperCase()}</p>
        <p>Type: {item.type}</p>
        <p>Notes: {item.notes}</p>
      </CardDescription>
    </CardHeader>
    {
      user &&
        <CardFooter className={`justify-end grid ${user ? "grid-cols-2" : "grid-cols-1"} gap-2`}>
          <Button variant="destructive" onClick={() => deleteFunction(item.id)}>
            <Trash2 />
          </Button>
          <Button variant="secondary" >
            <ShoppingCart />
          </Button>
        </CardFooter>
    }
  </Card>


export default function Items() {
  const [typeFilter, setTypeFilter] = useState<string>('Select Type')

  const { data: items, error, refetch, isFetching} = useItems(typeFilter)

  const { mutate: deleteItem } = useDeleteItem()
  const { data: user } = useQuery(meQueryOptions)
  const { data: types }  = useTypes()
  if (error) return <ErrorPage error={error} refetch={refetch} />

  if (!items || isFetching) return <LoadingPage />

  return (
    <div className="m-2 space-y-2">
      <div className="grid-cols-5">
        <Filter
          value={typeFilter}
          defaultValue="Select Type"
          onValueChange={(value: string | null) => setTypeFilter(value || "Select Type")}
          filterItems={types}
        />
      </div>
      <div className="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-2">
        {
          items.map((item) =>
            <Item key={item.id} item={item} user={user?.user} deleteFunction={deleteItem} />
          )
        }
      </div>
      {user?.user && <NewItem />}
    </div>
  )
}
