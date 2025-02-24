"use client";

import { Pencil } from "lucide-react";
import { Button } from "@/components/ui/button";
import { UpdateEmailForm } from "@/components/account/update-email-form";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import {
  AlertDialog,
  AlertDialogTrigger,
  AlertDialogContent,
} from "@/components/ui/alert-dialog";

import { StaticImageData } from "next/image";

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

interface UserEmailProps {
  userData: UserData;
}

export function EmailDisplay({ userData }: UserEmailProps) {
  return (
    <Card className="relative">
      <CardHeader>
        <CardTitle>Email Address</CardTitle>
      </CardHeader>
      <div className="grid sm:grid-cols-2 gap-x-6 gap-y-4">
        <CardContent>
          <CardDescription>Your current email address</CardDescription>
          <p>{userData.email}</p>
        </CardContent>
      </div>
      <AlertDialog>
        <AlertDialogTrigger asChild>
          <Button
            variant="outline"
            className="absolute top-0 right-0 mt-2 mr-2"
          >
            Edit <Pencil className="h-4 w-4 ml-2" />
          </Button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <UpdateEmailForm />
        </AlertDialogContent>
      </AlertDialog>
    </Card>
  );
}
