import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { request } from "./client";
import type { ItemT } from "./items";
import type { UserT } from "./user";

export type CheckoutItemT = {
  item: ItemT;
  returnDate: Date;
}

export type CheckoutT = {
  id: number;
  user: UserT;
  checkout_date: Date;
  // if a checkout is for personal use
  personal: boolean;
  notes: string | null;
  created_by: UserT;
  items: CheckoutItemT[];
}

export const checkoutKeys = {
  all: ['checkout'] as const,
  lists: (type?: string) => [...checkoutKeys.all, 'list', type] as const,
  details: () => [...checkoutKeys.all, 'detail'] as const,
  detail: (id: number) => [...checkoutKeys.details(), id] as const,
};

const fetchCheckouts = async(): Promise<CheckoutT[]> => {
  const response = await request(
    `${import.meta.env.VITE_API_URL}/checkouts`,
    {
      method: 'GET',
      credentials: 'include',
    }
  );
  if (response.status === 401) {
    throw new Error('Authentication required. Please log in.');
  }
  if (!response.ok) {
    throw new Error('Failed to fetch checkouts');
  }
  return await response.json() as Promise<CheckoutT[]>;
}

const fetchCheckout = async(id: number): Promise<CheckoutT> => {
  const response = await request(`${import.meta.env.VITE_API_URL}/checkouts?id=${id}`, {
    credentials: 'include',
  });
  if (response.status === 401) {
    throw new Error('Authentication required. Please log in.');
  }
  if (!response.ok) {
    throw new Error('Failed to fetch checkout');
  }
  return await response.json() as Promise<CheckoutT>;
}

const createCheckout = async(items: ItemT[]): Promise<ItemT> => {
  const response = await request(
    `${import.meta.env.VITE_API_URL}/items`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(items),
      credentials: 'include',
    }
  )
  if (!response.ok) {
    throw new Error('Failed to create checkout');
  }
  return response.json() as Promise<ItemT>;
}

const returnCheckout = async({ id, items }: { id: number; items: string[] }): Promise<void> => {
  const response = await request(`${import.meta.env.VITE_API_URL}/items?id=${id}`, {
    method: 'PUT',
    credentials: 'include',
    body: JSON.stringify({
      "id": id,
      "items": items
    })
  })
  if (!response.ok) {
    throw new Error('Failed to return checkout');
  }
}

export const useCheckouts = () => {
  return useQuery({
    queryKey: checkoutKeys.lists(),
    queryFn: () => fetchCheckouts(),
    retry: false,
  });
};

export const useCheckout = (id: number) => {
  return useQuery({
    queryKey: checkoutKeys.detail(id),
    queryFn: () => fetchCheckout(id),
    enabled: !!id,
    retry: false,
  });
};

export const useCreateItem = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createCheckout,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: checkoutKeys.all });
    },
  });
};

export const useReturnCheckout = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: returnCheckout,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: checkoutKeys.all });
    }
  })
}
