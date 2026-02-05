"use client";

import { useAuth } from "@/hooks/use-auth";

export default function HomePage() {
  const { user } = useAuth();

  return (
    <div>
      <h1 className="text-2xl font-bold">Dashboard</h1>
      
      {user && (
        <p className="mt-4 text-muted-foreground">
          Welcome back, <strong>{user.username}</strong>!
        </p>
      )}
      
      <div className="mt-8">
        <p>You are successfully authenticated!</p>
      </div>
    </div>
  );
}
