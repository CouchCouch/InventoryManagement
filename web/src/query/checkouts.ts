import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { ItemT } from "./items";
import type { AdminT, UserT } from "./user";

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

/*
export const itemKeys = {
  all: ['items'] as const,
  lists: (type?: string) => [...itemKeys.all, 'list', type] as const,
  details: () => [...itemKeys.all, 'detail'] as const,
  detail: (id: string) => [...itemKeys.details(), id] as const,
};
*/

const fetchCheckouts = async(): Promise<CheckoutT[]> => {
  const response = await fetch(
    `${import.meta.env.VITE_API_URL}/checkouts`,
    {
      method: 'GET',
      credentials: 'include',
    }
  );
  if (!response.ok) {
    throw new Error('Failed to fetch checkouts');
  }
  return await response.json() as Promise<CheckoutT[]>;
}

/*
const fetchItem = async(id?: string): Promise<ItemT[]> => {
  const response = await fetch(`${import.meta.env.VITE_API_URL}/items?id=${id}`);
  return await response.json() as Promise<ItemT[]>;
}

const createItem = async(item: ItemT): Promise<ItemT> => {
  const response = await fetch(
    `${import.meta.env.VITE_API_URL}/items`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(item),
      credentials: 'include',
    }
  )
  return response.json() as Promise<ItemT>;
}

const deleteItem = async(id: string): Promise<void> => {
  await fetch(`${import.meta.env.VITE_API_URL}/items?id=${id}`, {
    method: 'DELETE',
    credentials: 'include',
  })
}

const fetchTypes = async(): Promise<string[]> => {
  const response = await fetch(`${import.meta.env.VITE_API_URL}/items/types`);
  return await response.json() as Promise<string[]>;
}
*/

export const useCheckouts = () => {
  return useQuery({
    queryKey: ['checkouts'],
    queryFn: () => fetchCheckouts(),
  });
};

/*
export const useItem = (id: string) => {
  return useQuery({
    queryKey: itemKeys.detail(id),
    queryFn: () => fetchItem(id),
    enabled: !!id,
  });
};

export const useCreateItem = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createItem,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: itemKeys.lists() });
    },
  });
};

export const useTypes = () => {
  return useQuery({
    queryKey: ['itemTypes'],
    queryFn: fetchTypes,
  })
}

export const useDeleteItem = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteItem,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: itemKeys.lists() });
    }
  })
}
*/
