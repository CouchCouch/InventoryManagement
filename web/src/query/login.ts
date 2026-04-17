import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";

export type LoginInfoT = {
  email: string;
  password: string;
}

export type LoginResponseT = {
   login: "success" |  "fail"
}

const login = async(loginInfo: LoginInfoT): Promise<LoginResponseT> => {
  const response = await fetch(
    `${import.meta.env.VITE_API_URL}/auth/login`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(loginInfo),
      credentials: 'include',
    }
  );
  if (!response.ok) {
    throw new Error('Login failed');
  }
  return response.json()
}

const logout = async(): Promise<void> => {
  await fetch(`${import.meta.env.VITE_API_URL}/auth/logout`, {
    method: 'POST',
    credentials: 'include',
  })
}

export const useLogin = () => {
  const queryClient = useQueryClient();
  const navigate = useNavigate()

  return useMutation({
    mutationFn: login,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['me'] });
      navigate({ to: '/' })
    },
  })
}

export const useLogout = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: logout,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['me'] });
    }
  })
}
