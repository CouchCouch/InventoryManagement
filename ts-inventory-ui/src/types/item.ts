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
    const response = await fetch('http://localhost:8080/items')
    if(!response.ok) {
        throw response
    }
    const data = await response.json()
    return data
}

export { fetchItems }
