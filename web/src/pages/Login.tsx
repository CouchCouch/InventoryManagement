import { Field, FieldGroup, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { login } from "../query/login"
import { useState } from "react"
import { Button } from "@/components/ui/button"
import { redirect } from '@tanstack/react-router'

export default function Login() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  const handleSubmit = (e: SubmitEvent<HTMLFormElement>) => {
    e.preventDefault()
    login(email, password).then((response) => {
      if (response === 'success') {
        console.log('Login successful')
        throw redirect({
          to: '/items',
        })
      }
    })
  }

  return (
    <>
      <form onSubmit={(e) => handleSubmit(e)} className="m-2">
        <h2 className="text-center my-4">Login</h2>
        <FieldGroup className="w-1/3 m-auto">
          <Field>
            <FieldLabel htmlFor="email-input">Email</FieldLabel>
            <Input
              id="email-input"
              type="email"
              placeholder="Email"
              onChange={(e) => setEmail(e.target.value)}
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="password-input">Password</FieldLabel>
            <Input
              id="password-input"
              type="password"
              placeholder="Password"
              onChange={(e) => setPassword(e.target.value)}
            />
          </Field>
          <Button type="submit" className="w-full">
            Login
          </Button>
        </FieldGroup>
      </form>
    </>
  )
}
