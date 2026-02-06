"use client"

import * as React from "react"
import {
  Bot,
  Users,
  FolderKanban,
  ListTodo,
} from "lucide-react"

import { NavMain } from "./nav-main"
import { NavUser } from "./nav-user"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarRail,
} from "@/components/ui/sidebar"
import { useAuth } from "@/hooks/use-auth"
import { Skeleton } from "@/components/ui/skeleton"

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { user, isLoading } = useAuth()

  // Transform user data from API format to component format
  const userData = user
    ? {
        name: user.username,
        email: user.email,
        avatar: user.avatar_url,
      }
    : null

  // Add User Management for admin users and Projects for all users
  const navMainItems = React.useMemo(() => {
    const items = [
      {
        title: "Projects",
        url: "/projects",
        icon: FolderKanban,
        items: [
          {
            title: "All Projects",
            url: "/projects",
          },
        ],
      },
      {
        title: "Tasks",
        url: "/tasks",
        icon: ListTodo,
        items: [
          {
            title: "All Tasks",
            url: "/tasks",
          },
        ],
      },
      {
        title: "AI Agents",
        url: "/agents",
        icon: Bot,
        items: [
          {
            title: "All Agents",
            url: "/agents",
          },
        ],
      },
    ]

    if (user?.role === "admin") {
      items.push({
        title: "User Management",
        url: "/users",
        icon: Users,
        items: [
          {
            title: "All Users",
            url: "/users",
          },
        ],
      })
    }

    return items
  }, [user?.role])

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarContent>
        <NavMain items={navMainItems} />
      </SidebarContent>
      <SidebarFooter>
        {isLoading ? (
          <div className="flex items-center gap-3 px-3 py-2">
            <Skeleton className="h-8 w-8 rounded-lg" />
            <div className="flex flex-col gap-1">
              <Skeleton className="h-4 w-24" />
              <Skeleton className="h-3 w-32" />
            </div>
          </div>
        ) : userData ? (
          <NavUser user={userData} />
        ) : null}
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
