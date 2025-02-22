"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { toast } from "@/hooks/use-toast";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";

const FormSchema = z.object({
  username: z
    .string()
    .min(2, {
      message: "Exercise name must be at least 2 characters.",
    })
    .max(30, {
      message: "Exercise name cannot be more than 30 characters.",
    }),
});

export function UpdateUsername() {
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      username: "",
    },
  });

  function onSubmit(data: z.infer<typeof FormSchema>) {
    toast({
      title: "You submitted the following values:",
      description: (
        <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
          <code className="text-white">{JSON.stringify(data, null, 2)}</code>
        </pre>
      ),
    });
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="w-2/3 space-y-6">
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Update your username</FormLabel>
              <FormControl>
                <div className="flex w-full max-w-sm items-center space-x-2">
                  {/* dynamically render with the logged in users name */}
                  <Input placeholder="new username" {...field} />
                  <Button type="submit">Submit</Button>
                </div>
              </FormControl>
              <FormDescription className="gap-1">
                This is your public display name.
                <br />
                Please enter a new username and press update to confirm the
                changes.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
      </form>
    </Form>
  );
}
