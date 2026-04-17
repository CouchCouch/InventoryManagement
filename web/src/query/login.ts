export const login = async(email: string, password: string): Promise<string> => {
  const response = await fetch(
    `${import.meta.env.VITE_API_URL}/auth/login`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email: email , password: password }),
      credentials: 'include',
    }
  );
  if (!response.ok) {
    throw new Error('Network response was not ok');
  }
  return 'success'
}

export const loginQueryOptions = {
  queryKey: ['login'],
  queryFn: login,
}
