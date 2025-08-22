const APIURL = 'http://localhost:3000/api/'

import { ItemT } from "./item"

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
  const response = await fetch(`${APIURL}/checkout`)
  if(!response.ok) {
    throw response
  }
  const data = await response.json()
  return data
}

async function fetchHistory(id: number): Promise<CheckoutResponse> {
  const response = await fetch(`${APIURL}/checkout?id=` + id)
  if(!response.ok) {
    throw response
  }
  const data = response.json()
  return data
}

async function createCheckout(item: ItemT, name: string, email: string) {
  const body = {
    "name": name,
    "email": email,
    "id": item.id
  }
  const response = await fetch(`${APIURL}/checkout`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body)
    })
  if(!response.ok) {
    throw response
  }
  const data = response.json()
  return data
}

export { fetchCheckouts, fetchHistory, createCheckout }
