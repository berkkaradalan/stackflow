"use client";

import { LoginForm } from "@/components/login-form";
import { useRedirectIfAuthenticated } from "@/hooks/use-auth";

export default function LoginPage() {
  // Redirect to dashboard if already logged in
  useRedirectIfAuthenticated();

  return (
    <div className="bg-background flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="w-full max-w-sm">
        <LoginForm />
      </div>
    </div>
  );
}
