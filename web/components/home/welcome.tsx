"use client";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface UserData {
  username: string;
}

export function WelcomeCard({ username }: UserData) {
  const currentDateTime = new Date();

  const formattedDate = new Intl.DateTimeFormat("en-US", {
    weekday: "long",
    month: "long",
    day: "2-digit",
    year: "numeric",
  }).format(currentDateTime);

  const formattedTime = new Intl.DateTimeFormat("en-GB", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  }).format(currentDateTime);

  return (
    <Card className=" text-center p-6 shadow-md">
      <CardHeader className="space-y-2">
        <CardDescription className="text-gray-700 text-sm">
          Welcome back <span className="font-semibold">{username}</span>.
        </CardDescription>
        <p className="text-gray-800 text-sm">{formattedDate}</p>
      </CardHeader>
      <CardContent>
        <CardTitle className="text-5xl font-bold">
          {formattedTime}
        </CardTitle>
      </CardContent>
      <Button className="mt-4">
        Start New Workout
      </Button>
    </Card>
  );
}
