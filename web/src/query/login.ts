export const login = async(email: string, password: string): Promise<string> => {
  const response = await fetch(
    `http://localhost:3000/api/auth/login`,
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
