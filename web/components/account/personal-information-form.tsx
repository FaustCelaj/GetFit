"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { StaticImageData } from "next/image";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";

import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { Input } from "../ui/input";
import { Textarea } from "@/components/ui/textarea";

interface UserData {
  picture: StaticImageData;
  name: string;
  lastName: string;
  username: string;
  email: string;
  title: string;
  bio: string;
  age: number;
}

interface UserPersonalInformationProps {
  userData: UserData;
}

const FormSchema = z.object({
  name: z
    .string()
    .min(2, { message: "First name must be at least 2 characters" })
    .max(50, {
      message: "First name cannot be more than 50 characters.",
    }),
  lastName: z
    .string()
    .min(2, { message: "Last name must be at least 2 characters" })
    .max(50, {
      message: "Last name cannot be more than 50 characters.",
    }),
  age: z.number(),
  title: z.string(),
  bio: z
    .string()
    .max(160, { message: "Bio cannot be more than 140 characters" }),
});

export function UpdatePersonalInformationForm({
  userData,
}: UserPersonalInformationProps) {
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      name: userData.name || "",
      lastName: userData.lastName || "",
      age: userData.age || 0,
      title: userData.title || "",
      bio: userData.bio || "",
    },
  });

  function onSubmit(data: z.infer<typeof FormSchema>) {
    console.log(data);
  }

  return (
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle> Update Personal Information</AlertDialogTitle>
      </AlertDialogHeader>
      <AlertDialogDescription>
        Please fill in the information you would like to change. Only new
        information will be registered upon update.
      </AlertDialogDescription>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <div className="grid gap-2 mb-4">
                  <FormLabel>First Name</FormLabel>
                  <FormControl>
                    <Input placeholder={userData.name} {...field} />
                  </FormControl>
                  <FormMessage />
                </div>
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="lastName"
            render={({ field }) => (
              <FormItem>
                <div className="grid gap-2 mb-4">
                  <FormLabel>Surname</FormLabel>
                  <FormControl>
                    <Input placeholder={userData.lastName} {...field} />
                  </FormControl>
                  <FormMessage />
                </div>
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="age"
            render={({ field }) => (
              <FormItem>
                <div className="grid gap-2 mb-4">
                  <FormLabel>Age</FormLabel>
                  <FormControl>
                    <Input type="number" placeholder={userData.age} max={120} />
                  </FormControl>
                  <FormMessage />
                </div>
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="title"
            render={({ field }) => (
              <FormItem>
                <div className="grid gap-2 mb-4">
                  <FormLabel>Title</FormLabel>
                  <FormControl>
                    <Input placeholder={userData.title} {...field} />
                  </FormControl>
                  <FormMessage />
                </div>
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="bio"
            render={({ field }) => (
              <FormItem>
                <div className="grid gap-2 mb-4">
                  <FormLabel>Bio</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder={userData.bio}
                      className="resize-none"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </div>
              </FormItem>
            )}
          />
          <AlertDialogFooter className="grid grid-cols-2 gap-4">
            <AlertDialogAction type="submit">Update</AlertDialogAction>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
          </AlertDialogFooter>
        </form>
      </Form>
    </AlertDialogContent>
  );
}
