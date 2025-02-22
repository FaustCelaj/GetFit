const example = {
  name: "3/4 Sit-Up",
  force: "pull",
  level: "beginner",
  mechanic: "compound",
  equipment: "body only",
  primaryMuscles: ["abdominals"],
  secondaryMuscles: ["forearms", "biceps"],
  instructions: [
    "Lie down on the floor and secure your feet. Your legs should be bent at the knees.",
    "Place your hands behind or to the side of your head. You will begin with your back on the ground. This will be your starting position.",
    "Flex your hips and spine to raise your torso toward your knees.",
    "At the top of the contraction your torso should be perpendicular to the ground. Reverse the motion, going only ¾ of the way down.",
    "Repeat for the recommended amount of repetitions.",
  ],
  category: "strength",
};

import * as React from "react";

import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import { ExerciseType } from "@/types/exercise";
interface ExerciseSmallProps {
  exercise: ExerciseType;
  onViewMore: () => void;
}

export function ExerciseSmall({ exercise, onViewMore }: ExerciseSmallProps) {
  return (
    <Card className="w-[400px]">
      <CardHeader className="flex flex-row justify-between align-middle">
        <CardTitle className="text-lg">{exercise.name}</CardTitle>
        <Badge className="text-sm">{exercise.category}</Badge>
      </CardHeader>
      <Separator className="my-1" />
      <CardDescription className="py-2 px-4">
        The {exercise.name}, a {exercise.mechanic} movement that primarily
        targets the {exercise.primaryMuscles.join(", ")}.
      </CardDescription>
      <CardFooter className="flex-1 justify-end">
        <Button variant="outline" onClick={onViewMore}>
          View More
        </Button>
      </CardFooter>
    </Card>
  );
}
