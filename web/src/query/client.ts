export async function request(url: string, opts?: RequestInit): Promise<Response> {
  let response: Response
  try {
    response = await fetch(url, opts)
  } catch {
    throw new Error('Unable to connect to the server. Check your connection.')
  }
  return response
}
