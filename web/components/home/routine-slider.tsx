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
    <Card className="h-full bg-white rounded-xl shadow-md">
      <CardHeader className="pb-2">
        <CardTitle className="text-lg font-bold">Your Routines</CardTitle>
      </CardHeader>
      <CardContent className="overflow-y-auto">
        <Carousel className="h-full">
          <CarouselContent className="h-full">
            {routines.map((routine, index) => (
              <CarouselItem key={index} className="h-full p-2">
                <Card className="shadow-sm rounded-xl h-full flex flex-col">
                  <CardHeader className="pb-2">
                    <CardTitle className="text-xl font-bold">
                      {routine.title}
                    </CardTitle>
                    <CardDescription className="text-gray-600">
                      {routine.description}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="flex-grow overflow-auto py-2">
                    <div className="space-y-3">
                      {routine.exercises.map((exercise, i) => (
                        <div key={i} className="border-b pb-2">
                          <p className="text-md font-medium">{exercise.name}</p>
                          <p className="text-gray-600 text-sm">
                            {exercise.totalSets} sets x {exercise.totalReps} reps
                          </p>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                  <CardFooter className="pt-2">
                    <Button className="w-full">Start Workout</Button>
                  </CardFooter>
                </Card>
              </CarouselItem>
            ))}
          </CarouselContent>
          <CarouselPrevious className="left-1" />
          <CarouselNext className="right-1" />
        </Carousel>
      </CardContent>
    </Card>
  );
}