import { request } from "./client"

export type UserT = {
  id: string;
  name: string;
  email: string;
}

export type AdminT = {
  user: UserT;
  role: string;
}

const fetchMe = async(): Promise<AdminT> => {
  const response = await request(
    `${import.meta.env.VITE_API_URL}/users/me`,
    {
      method: 'GET',
      credentials: 'include',
    }
  );
  if (response.status === 401) {
    throw new Error('Authentication required. Please log in.');
  }
  if (!response.ok) {
    throw new Error('Failed to fetch user info');
  }
  return await response.json() as Promise<AdminT>;
}

export const meQueryOptions = {
  queryKey: ['me'],
  queryFn: fetchMe,
  retry: false,
}
