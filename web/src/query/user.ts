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
  const response = await fetch(
    `${import.meta.env.VITE_API_URL}/users/me`,
    {
      method: 'GET',
      credentials: 'include',
    }
  );
  return await response.json() as Promise<AdminT>;
}

export const meQueryOptions = {
  queryKey: ['me'],
  queryFn: fetchMe,
}
