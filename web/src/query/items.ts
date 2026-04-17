import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export type ItemT = {
  id: string;
  name: string;
  type: string;
  notes: string | null;
  date_purchased: Date | null;
}

export const itemKeys = {
  all: ['items'] as const,
  lists: () => [...itemKeys.all, 'list'] as const,
  details: () => [...itemKeys.all, 'detail'] as const,
  detail: (id: string) => [...itemKeys.details(), id] as const,
};

const fetchItems = async(): Promise<ItemT[]> => {
  const response = await fetch(`${import.meta.env.VITE_API_URL}/items`);
  return await response.json() as Promise<ItemT[]>;
}

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

export const useItems = () => {
  return useQuery({
    queryKey: itemKeys.lists(),
    queryFn: fetchItems,
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
