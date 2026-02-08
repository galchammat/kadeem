import { useEffect } from "react"
import { useNavigate } from "react-router"
import { Button } from "@/components/ui/button"
import { useAuth } from "@/contexts/authContext"

export default function LoginPage() {
  const { session, signInWithGoogle, loading } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (!loading && session) {
      navigate("/")
    }
  }, [session, loading, navigate])

  if (loading) {
    return <div className="flex h-screen items-center justify-center">Loading...</div>
  }

  return (
    <div className="flex h-screen w-full flex-col items-center justify-center gap-4">
      <h1 className="text-2xl font-bold">Sign in to Kadeem</h1>
      <Button onClick={() => signInWithGoogle()}>
        Sign in with Google
      </Button>
    </div>
  )
}
