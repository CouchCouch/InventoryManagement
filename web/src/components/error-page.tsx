import { Button } from "./ui/button"
import { RouterButton } from "./ui/router-button"

type ErrorPageProps = {
  error: Error
  refetch: () => void
}

const AUTH_ERROR = 'Authentication required. Please log in.'

export function ErrorPage({ error, refetch }: ErrorPageProps) {
  const isAuthError = error.message === AUTH_ERROR

  return (
    <div className="flex items-center justify-center h-full space-y-2">
      <div className="block text-center space-y-2">
        <h1>{error?.name}</h1>
        <h2>{error.message}</h2>
        {isAuthError ? (
          <RouterButton to="/login" variant="secondary" className="text-xl">
            Go to Login
          </RouterButton>
        ) : (
          <Button variant="secondary" onClick={() => refetch()} className="text-xl">Retry</Button>
        )}
      </div>
    </div>
  )
}
