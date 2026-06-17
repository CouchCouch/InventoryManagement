import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { request } from "./client";

export type ItemT = {
  id: string;
  name: string;
  type: string;
  notes: string | null;
  date_purchased: Date | null;
}

export type ItemTypeT = {
  name: string;
  filter: string;
}

export const itemKeys = {
  all: ['items'] as const,
  lists: (type?: string) => [...itemKeys.all, 'list', type] as const,
  details: () => [...itemKeys.all, 'detail'] as const,
  detail: (id: string) => [...itemKeys.details(), id] as const,
};

const fetchItems = async(filter?: string): Promise<ItemT[]> => {
  let requestUrl = `${import.meta.env.VITE_API_URL}/items`;

  if (filter && filter !== "Select Type") {
    requestUrl += `?type=${filter}`;
  }

  const response = await request(requestUrl);
  if (!response.ok) {
    throw new Error('Failed to fetch items');
  }
  return await response.json() as Promise<ItemT[]>;
}

const fetchItem = async(id?: string): Promise<ItemT[]> => {
  const response = await request(`${import.meta.env.VITE_API_URL}/items?id=${id}`);
  if (!response.ok) {
    throw new Error('Failed to fetch item');
  }
  return await response.json() as Promise<ItemT[]>;
}

const createItem = async(item: ItemT): Promise<ItemT> => {
  if (item.notes?.trim() === '') { item.notes = null }
  const response = await request(
    `${import.meta.env.VITE_API_URL}/items`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(item),
      credentials: 'include',
    }
  )
  if (!response.ok) {
    throw new Error('Failed to create item');
  }
  return response.json() as Promise<ItemT>;
}

const deleteItem = async(id: string): Promise<void> => {
  const response = await request(`${import.meta.env.VITE_API_URL}/items?id=${id}`, {
    method: 'DELETE',
    credentials: 'include',
  })
  if (!response.ok) {
    throw new Error('Failed to delete item');
  }
}

const fetchTypes = async(): Promise<string[]> => {
  const response = await request(`${import.meta.env.VITE_API_URL}/items/types`);
  if (!response.ok) {
    throw new Error('Failed to fetch types');
  }
  return await response.json() as Promise<string[]>;
}

export const useItems = (type?: string) => {
  return useQuery({
    queryKey: itemKeys.lists(type),
    queryFn: () => fetchItems(type),
  });
};

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
