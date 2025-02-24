import { Pencil } from "lucide-react";
import { Button } from "@/components/ui/button";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
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

interface UserPersonalInformationProps {
  userData: UserData;
}

export function PersonalInformationDisplay({
  userData,
}: UserPersonalInformationProps) {
  return (
    <Card className="relative">
      <CardHeader>
        <CardTitle>Personal Information</CardTitle>
      </CardHeader>
      <div className="grid sm:grid-cols-2 gap-x-6 gap-y-4">
        <CardContent>
          <CardDescription>First Name</CardDescription>
          <p>{userData.name}</p>
        </CardContent>
        <CardContent>
          <CardDescription>Last Name</CardDescription>
          <p>{userData.lastName}</p>
        </CardContent>
        <CardContent>
          <CardDescription>Age</CardDescription>
          <p>{userData.age}</p>
        </CardContent>
        <CardContent>
          <CardDescription>Title</CardDescription>
          <p>{userData.title}</p>
        </CardContent>
        <CardContent>
          <CardDescription>Bio</CardDescription>
          <p>{userData.bio}</p>
        </CardContent>
      </div>
      <Button variant="outline" className="absolute top-0 right-0 mt-2 mr-2">
        Edit <Pencil className="h-4 w-4 ml-2" />
      </Button>
    </Card>
  );
}
