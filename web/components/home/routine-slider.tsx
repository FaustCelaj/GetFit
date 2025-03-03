"use client";

import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselPrevious,
  CarouselNext,
} from "@/components/ui/carousel";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface WorkoutRoutine {
  title: string;
  description: string;
  exercises: Exercise[];
}

interface Exercise {
  name: string;
  totalSets?: number;
  totalReps?: number;
}

interface RoutineProps {
  routines: WorkoutRoutine[];
}

export function RoutineSlider({ routines }: RoutineProps) {
  return (
    <Carousel className="p-6 w-full max-w-xl max-h-max bg-red-400">
      <CarouselContent>
        {routines.map((routine, index) => (
          <CarouselItem key={index} className="p-4">
            <Card className="shadow-lg rounded-xl">
              <CardHeader>
                <CardTitle className="felx felx-row text-xl font-bold">
                  {routine.title}
                </CardTitle>
                <CardDescription className="text-gray-500">
                  {routine.description}
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4 max-h-max overflow-scroll">
                {routine.exercises.map((exercise, i) => (
                  <div
                    key={i}
                    className="border-b pb-2"
                  >
                    <p className="text-lg font-medium">{exercise.name}</p>
                    <p className="text-gray-600">
                      {exercise.totalSets} sets x {exercise.totalReps} reps
                    </p>
                  </div>
                ))}
              </CardContent>
              <CardFooter>
                <Button className="w-full mt-4">Start New Workout</Button>
              </CardFooter>
            </Card>
          </CarouselItem>
        ))}
      </CarouselContent>
      <CarouselPrevious />
      <CarouselNext />
    </Carousel>
  );
}
