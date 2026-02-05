"use client"

import * as React from "react"
import {
  AudioWaveform,
  BookOpen,
  Bot,
  Command,
  Frame,
  GalleryVerticalEnd,
  Map,
  PieChart,
  Settings2,
  SquareTerminal,
  Users,
  FolderKanban,
} from "lucide-react"

import { NavMain } from "./nav-main"
import { NavProjects } from "./nav-projects"
import { NavUser } from "./nav-user"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
} from "@/components/ui/sidebar"
import { useAuth } from "@/hooks/use-auth"
import { Skeleton } from "@/components/ui/skeleton"
import Link from "next/link"

// Static navigation data
const navData = {
  teams: [
    {
      name: "Acme Inc",
      logo: GalleryVerticalEnd,
      plan: "Enterprise",
    },
    {
      name: "Acme Corp.",
      logo: AudioWaveform,
      plan: "Startup",
    },
    {
      name: "Evil Corp.",
      logo: Command,
      plan: "Free",
    },
  ],
  navMain: [
    {
      title: "Playground",
      url: "#",
      icon: SquareTerminal,
      isActive: true,
      items: [
        {
          title: "History",
          url: "#",
        },
        {
          title: "Starred",
          url: "#",
        },
        {
          title: "Settings",
          url: "#",
        },
      ],
    },
    {
      title: "Models",
      url: "#",
      icon: Bot,
      items: [
        {
          title: "Genesis",
          url: "#",
        },
        {
          title: "Explorer",
          url: "#",
        },
        {
          title: "Quantum",
          url: "#",
        },
      ],
    },
    {
      title: "Documentation",
      url: "#",
      icon: BookOpen,
      items: [
        {
          title: "Introduction",
          url: "#",
        },
        {
          title: "Get Started",
          url: "#",
        },
        {
          title: "Tutorials",
          url: "#",
        },
        {
          title: "Changelog",
          url: "#",
        },
      ],
    },
    {
      title: "Settings",
      url: "#",
      icon: Settings2,
      items: [
        {
          title: "General",
          url: "#",
        },
        {
          title: "Team",
          url: "#",
        },
        {
          title: "Billing",
          url: "#",
        },
        {
          title: "Limits",
          url: "#",
        },
      ],
    },
  ],
  projects: [
    {
      name: "Design Engineering",
      url: "#",
      icon: Frame,
    },
    {
      name: "Sales & Marketing",
      url: "#",
      icon: PieChart,
    },
    {
      name: "Travel",
      url: "#",
      icon: Map,
    },
  ],
}

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
      ...navData.navMain,
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
      {/* <SidebarHeader>
        <TeamSwitcher teams={navData.teams} />
      </SidebarHeader> */}
      <SidebarContent>
        <NavMain items={navMainItems} />
        <NavProjects projects={navData.projects} />
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
