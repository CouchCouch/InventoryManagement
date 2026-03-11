import { Plus } from "lucide-react";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "./ui/dialog";
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";
import { Button } from "./ui/button";
import { Field, FieldGroup } from "./ui/field";
import { Label } from "./ui/label";
import { Input } from "./ui/input";

export default function NewItem() {
  return (
    <Dialog>
      <DialogTrigger>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button size="icon" className="rounded-full absolute bottom-8 right-2 h-20 w-20" onClick={() => console.log("clicked")}>
              <Plus className="w-16 h-16"/>
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            <p>Add a new item</p>
          </TooltipContent>
        </Tooltip>
      </DialogTrigger>
      <DialogContent className="sm:max-w-sm">
        <DialogHeader>
          <DialogTitle>Add Item</DialogTitle>
          <DialogDescription>
            Add items with the given info
          </DialogDescription>
        </DialogHeader>
        <FieldGroup>
          <Field>
            <Label htmlFor="name">Name</Label>
            <Input id="name" name="name" placeholder="Pickle" />
          </Field>
          <Field>
            <Label htmlFor="type">Type</Label>
            <Input id="type" name="type" placeholder="canoe" />
          </Field>
        </FieldGroup>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline">Cancel</Button>
          </DialogClose>
          <Button type="submit">Add Item</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
