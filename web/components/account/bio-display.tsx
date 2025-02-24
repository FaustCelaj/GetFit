import { Pencil } from "lucide-react";
import { Button } from "@/components/ui/button";

import {
  Card,
  CardContent,
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

interface UserBioProps {
  userData: UserData;
}

export function BioDisplay({ userData }: UserBioProps) {
  return (
    <Card className="relative">
      <CardContent className="flex flex-col gap-2 p-8 sm:flex-row sm:items-center sm:gap-6 sm:py-4">
        <img
          className="mx-auto block h-24 rounded-full sm:mx-0 sm:shrink-0"
          src={userData.picture.src}
          alt="User profile"
        />
        {/* text */}
        <div className="space-y-2 text-center sm:text-left">
          <p className="text-lg font-semibold text-black">
            {userData.username}
          </p>
          <p className="text-base font-medium text-gray-500">
            {userData.title}
          </p>
          <p className="text-sm font-medium text-gray-400">{userData.bio}</p>
        </div>
      </CardContent>
      <Button variant="outline" className="absolute top-0 right-0 mt-2 mr-2">
        Edit <Pencil className="h-4 w-4 ml-2" />
      </Button>
    </Card>
  );
}
