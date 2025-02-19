import { GalleryVerticalEnd } from "lucide-react"

import { LoginForm } from "@/components/login-form"

export default function LoginPage() {
  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-slate-100 p-6 md:p-10 dark:bg-slate-800">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <a href="#" className="flex items-center gap-2 self-center font-medium">
          <div className="flex h-6 w-6 items-center justify-center rounded-md bg-slate-900 text-slate-50 dark:bg-slate-50 dark:text-slate-900">
            <GalleryVerticalEnd className="size-4" />
          </div>
          GETFIT
        </a>
        <LoginForm />
      </div>
    </div>
  )
}
