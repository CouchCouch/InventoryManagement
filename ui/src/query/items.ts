export type Item = {
  id: string;
  name: string;
  type: string;
  notes: string;
  date_purchased: Date;
}

const fetchItems = async(): Promise<Item[]> => {
  const response = await fetch('http://localhost:3000/api/items');
  return await response.json() as Promise<Item[]>;
}

const fetchItem = async(id?: string): Promise<Item[]> => {
  const response = await fetch(`/api/items?id=${id}`);
  return await response.json() as Promise<Item[]>;
}

export const itemsQueryOptions = {
  queryKey: ['items'],
  queryFn: fetchItems,
}

export const itemQueryOptions = {
  queryKey: ['items'],
  queryFn: fetchItem,
}
