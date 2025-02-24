"use client";
import { useState } from "react";

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import { Pencil, Eye, EyeClosed } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { StaticImageData } from "next/image";

interface UserData {
  picture: StaticImageData;
  name: string;
  lastName: string;
  username: string;
  password: string;
  email: string;
  title: string;
  bio: string;
  age: number;
}

interface UserPasswordProps {
  userData: UserData;
}

export function PasswordDisplay({ userData }: UserPasswordProps) {
  const [showPassword, setShowPassword] = useState(false);

  const togglePasswordVisability = () => {
    console.log("Toggling password visibility. Current state:", showPassword);
    setShowPassword(!showPassword);
  };

  return (
    <Card className="relative">
      <CardHeader>
        <CardTitle>Password</CardTitle>
      </CardHeader>
      <div className="grid sm:grid-cols-2 gap-x-6 gap-y-4">
        <CardContent>
          <CardDescription>Your current password</CardDescription>
          <div className="flex w-full max-w-sm items-center space-x-2 pt-2">
            <Input
              id="password"
              type={showPassword ? "text" : "password"}
              placeholder={userData.password}
              readOnly
            />
            <Button variant="ghost" onClick={togglePasswordVisability}>
              {showPassword ? <EyeClosed /> : <Eye />}
            </Button>
          </div>
        </CardContent>
      </div>
      <Button variant="outline" className="absolute top-0 right-0 mt-2 mr-2">
        Edit <Pencil className="h-4 w-4 ml-2" />
      </Button>
    </Card>
  );
}
