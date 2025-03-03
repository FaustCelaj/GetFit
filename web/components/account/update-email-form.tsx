"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { toast } from "@/hooks/use-toast";
import { Input } from "@/components/ui/input";

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
} from "@/components/ui/form";

const FormSchema = z.object({
  email: z.string().email({ message: "Invalid email address" }),
});

export function UpdateEmailForm() {
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      email: "",
    },
  });

  function onSubmit(data: z.infer<typeof FormSchema>) {
    console.log(data);
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
    <AlertDialogContent className="max-w-96 md:max-w-lg rounded-md">
      <AlertDialogHeader>
        <AlertDialogTitle>Update Email Address</AlertDialogTitle>
      </AlertDialogHeader>
      <AlertDialogDescription>
        Please enter a new username and press update to confirm the changes.
      </AlertDialogDescription>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <div className="grid gap-2">
                    <Input
                      className="w-full my-2"
                      placeholder="new email"
                      {...field}
                    />
                  </div>
                </FormControl>
              </FormItem>
            )}
          />
          <AlertDialogFooter className="grid grid-cols-2 gap-4">
            <AlertDialogAction type="submit">Update</AlertDialogAction>
            <AlertDialogCancel >Cancel</AlertDialogCancel>
          </AlertDialogFooter>
        </form>
      </Form>
    </AlertDialogContent>
  );
}
