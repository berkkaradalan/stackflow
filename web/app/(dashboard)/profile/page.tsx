"use client";

import { useState } from "react";
import { useAuth } from "@/hooks/use-auth";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Loader2, User, Mail, Image, Lock, Eye, EyeOff } from "lucide-react";

export default function ProfilePage() {
  const { user, updateProfile, isLoading, error } = useAuth();
  const [showPassword, setShowPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [formData, setFormData] = useState({
    username: user?.username || "",
    email: user?.email || "",
    avatar_url: user?.avatar_url || "",
    old_password: "",
    new_password: "",
  });
  const [successMessage, setSuccessMessage] = useState("");

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    setSuccessMessage("");
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSuccessMessage("");

    // Prepare data - only include changed fields
    const dataToUpdate: {
      username?: string;
      email?: string;
      avatar_url?: string;
      old_password?: string;
      new_password?: string;
    } = {};

    if (formData.username !== user?.username) {
      dataToUpdate.username = formData.username;
    }
    if (formData.email !== user?.email) {
      dataToUpdate.email = formData.email;
    }
    if (formData.avatar_url !== user?.avatar_url) {
      dataToUpdate.avatar_url = formData.avatar_url;
    }
    if (formData.old_password && formData.new_password) {
      dataToUpdate.old_password = formData.old_password;
      dataToUpdate.new_password = formData.new_password;
    }

    // Only update if there are changes
    if (Object.keys(dataToUpdate).length === 0) {
      setSuccessMessage("No changes to save.");
      return;
    }

    try {
      await updateProfile(dataToUpdate);
      setSuccessMessage("Profile updated successfully!");
      // Clear password fields after successful update
      setFormData((prev) => ({
        ...prev,
        old_password: "",
        new_password: "",
      }));
    } catch {
      // Error is handled by auth context
    }
  };

  if (!user) {
    return (
      <div className="flex min-h-[400px] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Profile</h1>
        <p className="text-muted-foreground">
          Manage your account settings and update your profile information.
        </p>
      </div>

      <Separator />

      {/* Profile Preview Card */}
      <Card>
        <CardHeader>
          <CardTitle>Profile Preview</CardTitle>
          <CardDescription>
            This is how your profile appears across the application.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-4">
            <Avatar className="h-20 w-20">
              <AvatarImage
                src={user.avatar_url}
                alt={user.username}
                className="object-cover"
              />
              <AvatarFallback className="text-2xl">
                {user.username.charAt(0).toUpperCase()}
              </AvatarFallback>
            </Avatar>
            <div>
              <h3 className="text-xl font-semibold">{user.username}</h3>
              <p className="text-muted-foreground">{user.email}</p>
              <div className="mt-2 flex items-center gap-2">
                <span className="inline-flex items-center rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary">
                  {user.role}
                </span>
                <span
                  className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${
                    user.is_active
                      ? "bg-green-100 text-green-800"
                      : "bg-red-100 text-red-800"
                  }`}
                >
                  {user.is_active ? "Active" : "Inactive"}
                </span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Edit Profile Form */}
      <Card>
        <CardHeader>
          <CardTitle>Edit Profile</CardTitle>
          <CardDescription>
            Update your personal information and password.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
              <div className="rounded-lg bg-destructive/10 p-3 text-sm text-destructive">
                {error}
              </div>
            )}
            {successMessage && (
              <div className="rounded-lg bg-green-100 p-3 text-sm text-green-800">
                {successMessage}
              </div>
            )}

            <div className="space-y-4">
              {/* Username Field */}
              <div className="space-y-2">
                <Label htmlFor="username" className="flex items-center gap-2">
                  <User className="h-4 w-4" />
                  Username
                </Label>
                <Input
                  id="username"
                  name="username"
                  value={formData.username}
                  onChange={handleInputChange}
                  placeholder="Enter your username"
                />
              </div>

              {/* Email Field */}
              <div className="space-y-2">
                <Label htmlFor="email" className="flex items-center gap-2">
                  <Mail className="h-4 w-4" />
                  Email
                </Label>
                <Input
                  id="email"
                  name="email"
                  type="email"
                  value={formData.email}
                  onChange={handleInputChange}
                  placeholder="Enter your email"
                />
              </div>

              {/* Avatar URL Field */}
              <div className="space-y-2">
                <Label htmlFor="avatar_url" className="flex items-center gap-2">
                  <Image className="h-4 w-4" />
                  Avatar URL
                </Label>
                <Input
                  id="avatar_url"
                  name="avatar_url"
                  type="url"
                  value={formData.avatar_url}
                  onChange={handleInputChange}
                  placeholder="https://example.com/avatar.png"
                />
              </div>

              <Separator />

              {/* Password Change Section */}
              <div className="space-y-2">
                <Label className="flex items-center gap-2 text-base font-semibold">
                  <Lock className="h-4 w-4" />
                  Change Password
                </Label>
                <p className="text-sm text-muted-foreground">
                  Leave blank if you don&apos;t want to change your password.
                </p>
              </div>

              {/* Old Password Field */}
              <div className="space-y-2">
                <Label htmlFor="old_password">Current Password</Label>
                <div className="relative">
                  <Input
                    id="old_password"
                    name="old_password"
                    type={showPassword ? "text" : "password"}
                    value={formData.old_password}
                    onChange={handleInputChange}
                    placeholder="Enter current password"
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? (
                      <EyeOff className="h-4 w-4 text-muted-foreground" />
                    ) : (
                      <Eye className="h-4 w-4 text-muted-foreground" />
                    )}
                  </Button>
                </div>
              </div>

              {/* New Password Field */}
              <div className="space-y-2">
                <Label htmlFor="new_password">New Password</Label>
                <div className="relative">
                  <Input
                    id="new_password"
                    name="new_password"
                    type={showNewPassword ? "text" : "password"}
                    value={formData.new_password}
                    onChange={handleInputChange}
                    placeholder="Enter new password"
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                    onClick={() => setShowNewPassword(!showNewPassword)}
                  >
                    {showNewPassword ? (
                      <EyeOff className="h-4 w-4 text-muted-foreground" />
                    ) : (
                      <Eye className="h-4 w-4 text-muted-foreground" />
                    )}
                  </Button>
                </div>
              </div>
            </div>

            <div className="flex justify-end gap-4">
              <Button
                type="button"
                variant="outline"
                onClick={() =>
                  setFormData({
                    username: user.username,
                    email: user.email,
                    avatar_url: user.avatar_url,
                    old_password: "",
                    new_password: "",
                  })
                }
                disabled={isLoading}
              >
                Reset
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Saving...
                  </>
                ) : (
                  "Save Changes"
                )}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      {/* Account Info Card */}
      <Card>
        <CardHeader>
          <CardTitle>Account Information</CardTitle>
          <CardDescription>
            Your account details and statistics.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <dl className="space-y-4">
            <div className="flex justify-between border-b pb-2">
              <dt className="text-muted-foreground">User ID</dt>
              <dd className="font-medium">{user.id}</dd>
            </div>
            <div className="flex justify-between border-b pb-2">
              <dt className="text-muted-foreground">Role</dt>
              <dd className="font-medium capitalize">{user.role}</dd>
            </div>
            <div className="flex justify-between border-b pb-2">
              <dt className="text-muted-foreground">Status</dt>
              <dd className="font-medium">
                {user.is_active ? "Active" : "Inactive"}
              </dd>
            </div>
            <div className="flex justify-between border-b pb-2">
              <dt className="text-muted-foreground">Member Since</dt>
              <dd className="font-medium">
                {new Date(user.created_at).toLocaleDateString()}
              </dd>
            </div>
            <div className="flex justify-between">
              <dt className="text-muted-foreground">Last Updated</dt>
              <dd className="font-medium">
                {new Date(user.updated_at).toLocaleDateString()}
              </dd>
            </div>
          </dl>
        </CardContent>
      </Card>
    </div>
  );
}
