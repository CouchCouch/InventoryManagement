const APIURL = 'http://localhost:3000/api/'

export interface ItemT {
  id: number,
  name: string,
  description: string,
  quantity: number
}

export interface ItemResponse {
  code: number,
  items: ItemT[]
}

export interface NewItemResponse {
  code: number,
  id: number
}

async function fetchItems(): Promise<ItemResponse> {
  const response = await fetch(`${APIURL}items`)
  if(!response.ok) {
    throw response
  }
  const data = await response.json()
  return data
}

async function fetchItem(id: number): Promise<ItemResponse> {
  const response = await fetch(`${APIURL}/items?id=` + id)
  if(!response.ok) {
    throw response
  }
  const data = await response.json()
  return data
}

export { fetchItems, fetchItem }
