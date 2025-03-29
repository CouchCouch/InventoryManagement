export interface CheckoutT {
    id: number,
    itemId: number,
    itemName: string,
    name: string,
    email: string,
    date: Date,
    returned: boolean
}

export interface CheckoutResponse {
    code: number,
    checkouts: CheckoutT[]
}

async function fetchCheckouts(): Promise<CheckoutResponse> {
    const response = await fetch('http://localhost:8080/checkout')
    if(!response.ok) {
        throw response
    }
    const data = await response.json()
    return data
}

async function fetchHistory(id: number): Promise<CheckoutResponse> {
    const response = await fetch('http://localhost:8080/checkout?id=' + id)
    if(!response.ok) {
        throw response
    }
    const data = await response.json()
    return data
}

export { fetchCheckouts, fetchHistory }
