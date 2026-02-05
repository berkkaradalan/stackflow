"use client";

import { Button } from "@/components/ui/button";
import { useAuth } from "@/hooks/use-auth";

export default function HomePage() {
  const { user, logout } = useAuth();

  return (
    <div className="p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <div className="flex items-center gap-4">
          {user && (
            <span className="text-muted-foreground">
              Welcome, <strong>{user.username}</strong> ({user.role})
            </span>
          )}
          <Button variant="outline" onClick={() => logout()}>
            Logout
          </Button>
        </div>
      </div>
      
      <div className="mt-8">
        <p>You are successfully authenticated!</p>
      </div>
    </div>
  );
}
