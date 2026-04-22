import { createLink } from '@tanstack/react-router'
import { Button, type buttonVariants } from '@/components/ui/button'
import { Button as ButtonPrimitive } from "@base-ui/react/button"
import { forwardRef } from 'react'
import type { VariantProps } from 'class-variance-authority'

type ButtonProps = ButtonPrimitive.Props & VariantProps<typeof buttonVariants>

// Create a router-compatible Button
export const RouterButton = createLink(
  forwardRef<HTMLButtonElement, ButtonProps>((props, ref) => {
    return <Button ref={ref} {...props} />
  }),
)
