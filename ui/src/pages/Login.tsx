import { Field, FieldGroup, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { login } from "../query/login"
import { useState } from "react"

export default function Login() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  const handleSubmit = () => {
    login(email, password)
  }

  return (
    <>
      <form onSubmit={() => handleSubmit()}>
        <FieldGroup>
          <Field>
            <FieldLabel htmlFor="email-input">Email</FieldLabel>
            <Input
              id="email-input"
              type="email"
              placeholder="abc1234@rit.edu"
              onChange={(e) => setEmail(e.target.value)}
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="password-input">Password</FieldLabel>
            <Input
              id="password-input"
              type="password"
              onChange={(e) => setPassword(e.target.value)}
            />
          </Field>
        </FieldGroup>
      </form>
    </>
  )
}
