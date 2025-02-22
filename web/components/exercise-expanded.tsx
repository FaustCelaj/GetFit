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

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ExerciseType } from "@/types/exercise";

interface ExerciseExpandedProps {
  exercise: ExerciseType;
}

export function ExerciseExpanded({ exercise }: ExerciseExpandedProps) {
  return (
    <Tabs defaultValue="summary" className="w-[400px]">
      <TabsList className="grid w-full grid-cols-2">
        <TabsTrigger value="summary">Summary</TabsTrigger>
        <TabsTrigger value="how-to">How To</TabsTrigger>
      </TabsList>

      {/* Summary Tab */}
      <TabsContent value="summary">
        <Card>
          <CardHeader>
            <CardTitle>{exercise.name}</CardTitle>
            <CardDescription>
              <strong>Primary Muscles:</strong> {exercise.primaryMuscles.join(", ")}
            </CardDescription>
            <CardDescription>
              <strong>Secondary Muscles:</strong> {exercise.secondaryMuscles.join(", ")}
            </CardDescription>
            <CardDescription>
              <strong>Force:</strong> {exercise.force}
            </CardDescription>
            <CardDescription>
              <strong>Level:</strong> {exercise.level}
            </CardDescription>
            <CardDescription>
              <strong>Mechanic:</strong> {exercise.mechanic}
            </CardDescription>
            <CardDescription>
              <strong>Equipment:</strong> {exercise.equipment}
            </CardDescription>
            <CardDescription>
              <strong>Category:</strong> {exercise.category}
            </CardDescription>
          </CardHeader>
        </Card>
      </TabsContent>

      {/* How To Tab */}
      <TabsContent value="how-to">
        <Card>
          <CardHeader>
            <CardTitle>{exercise.name}</CardTitle>
            <CardDescription>
              <strong>Instructions:</strong>
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ol className="list-decimal pl-5">
              {exercise.instructions.map((step, index) => (
                <li key={index} className="mb-2">
                  {step}
                </li>
              ))}
            </ol>
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>
  );
}