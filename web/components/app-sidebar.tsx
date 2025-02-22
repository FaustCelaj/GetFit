import * as React from "react"
import { Omega, Home, Dumbbell, Search, Settings, LayoutList } from "lucide-react"

import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarHeader,
  SidebarFooter,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarRail,
} from "@/components/ui/sidebar"

import { Button } from "./ui/button"

// This is sample data.
const data = {
  navMain: [
    {
      title: "Home",
      url: "#",
      icon: Home,
      isActive: true,
    },
    {
      title: "Routines",
      url: "#",
      icon: LayoutList,
      items: [
        {
          title: "View My Routines",
          url: "#",
        },
        {
          title: "Create A Routine",
          url: "#",
        },
      ],
    },
    {
      title: "Exercise Library",
      url: "#",
      icon: Search,
      items: [
        {
          title: "Search",
          url: "#",
        },
        {
          title: "View My Exercises",
          url: "#",
        },
        {
          title: "Create New Exercise",
          url: "#",
        },
      ],
    },
    {
      title: "Workout",
      url: "#",
      icon: Dumbbell,
      items: [
        {
          title: "Start Empty Workout",
          url: "#",
        },
        {
          title: "Start From Routine",
          url: "#",
        },
        {
          title: "History",
          url: "#",
        },
      ],
    },
    {
      title: "My Account",
      url: "#",
      icon: Settings,
      items: [
        {
          title: "Setings",
          url: "#",
        },
        {
          title: "Sign Out",
          url: "#",
        },
        {
          title: "Theme",
          url: "#",
        },
      ],
    },
  ],
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <a href="#">
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                  <Omega className="size-4" />
                </div>
                <div className="flex flex-col gap-0.5 leading-none">
                  <span className="font-semibold">GetFit</span>
                  <span className="">Personal Workout Tracker</span>
                </div>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent className="place-content-between">
        <SidebarGroup>
          <SidebarMenu>
            {data.navMain.map((item) => (
              <SidebarMenuItem key={item.title}>
                <SidebarMenuButton asChild>
                  <a href={item.url} className="font-medium">
                  <item.icon />

                    {item.title}
                  </a>
                </SidebarMenuButton>
                {item.items?.length ? (
                  <SidebarMenuSub>
                    {item.items.map((item) => (
                      <SidebarMenuSubItem key={item.title}>
                        {/* what is "isActive" and how to use */}
                        <SidebarMenuSubButton asChild isActive={item.isActive}>
                          <a href={item.url}>{item.title}</a>
                        </SidebarMenuSubButton>
                      </SidebarMenuSubItem>
                    ))}
                  </SidebarMenuSub>
                ) : null}
              </SidebarMenuItem>
            ))}
          </SidebarMenu>
        </SidebarGroup>
        <SidebarFooter>
            <Button className="m-2" variant={"destructive"}>Sign Out</Button>
        </SidebarFooter>
      </SidebarContent>
      <SidebarRail />
    </Sidebar>
  )
}
