"use client";

import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { cn } from "@/lib/utils";

import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

// Schema for form validation
const customExerciseSchema = z.object({
  name: z
    .string()
    .min(2, {
      message: "Exercise name must be at least 2 characters.",
    })
    .max(50, {
      message: "Exercise name cannot be more than 50 characters.",
    }),
  force: z.union([z.enum(["static", "pull", "push"]), z.null()]),
  level: z.union([z.enum(["beginner", "intermediate", "expert"]), z.null()]),
  mechanic: z.union([z.enum(["isolation", "compound"]), z.null()]),
  equipment: z.union([
    z.enum([
      "medicine ball",
      "dumbbell",
      "body only",
      "bands",
      "kettlebells",
      "foam roll",
      "cable",
      "machine",
      "barbell",
      "exercise ball",
      "e-z curl bar",
      "other",
    ]),
    z.null(),
  ]),
  primaryMuscles: z.union([
    z.enum([
      "abdominals",
      "abductors",
      "adductors",
      "biceps",
      "calves",
      "chest",
      "forearms",
      "glutes",
      "hamstrings",
      "lats",
      "lower back",
      "middle back",
      "neck",
      "quadriceps",
      "shoulders",
      "traps",
      "triceps",
    ]),
    z.null(),
  ]),
  secondaryMuscles: z.union([
    z.enum([
      "abdominals",
      "abductors",
      "adductors",
      "biceps",
      "calves",
      "chest",
      "forearms",
      "glutes",
      "hamstrings",
      "lats",
      "lower back",
      "middle back",
      "neck",
      "quadriceps",
      "shoulders",
      "traps",
      "triceps",
    ]),
    z.null(),
  ]),
  category: z.enum([
    "powerlifting",
    "strength",
    "stretching",
    "cardio",
    "olympic_weightlifting",
    "strongman",
    "plyometrics",
  ]),
});

// TODO: add toast for when the submit is complete and successful

export function CustomExerciseForm({
  className,
  ...props
}: React.ComponentPropsWithoutRef<"div">) {
  const form = useForm<z.infer<typeof customExerciseSchema>>({
    resolver: zodResolver(customExerciseSchema),
    defaultValues: {
      name: "",
      force: null,
      level: null,
      mechanic: null,
      equipment: null,
      primaryMuscles: null,
      secondaryMuscles: null,
      category: "strength",
    },
  });

  // 2. Define a submit handler.
  function onSubmit(values: z.infer<typeof customExerciseSchema>) {
    // Do something with the form values.
    console.log(values);
  }

  return (
    <div
      className={cn("flex flex-col gap-6 w-full max-w-md", className)}
      {...props}
    >
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Add Custom Exercise</CardTitle>
          <CardDescription>
            Please enter the details of your custom exercise.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
              {/* Name Field */}
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <div className="grid gap-2">
                      <FormLabel>Exercise Name</FormLabel>
                      <FormControl>
                        <Input placeholder="ex. Tibialis Raises" {...field} />
                      </FormControl>
                      <FormDescription>
                        What is your exercise called?
                      </FormDescription>
                      <FormMessage />
                    </div>
                  </FormItem>
                )}
              />
              {/* Force Field */}
              <FormField
                control={form.control}
                name="force"
                render={({ field }) => (
                  <FormItem>
                    <div className="grid gap-2">
                      <FormLabel>Exercise Force</FormLabel>
                      <FormControl>
                        <Select
                          onValueChange={(value) => {
                            field.onChange(value === "none" ? null : value);
                          }}
                          defaultValue={field.value || "none"}
                        >
                          <SelectTrigger
                            className={cn(
                              "w-full",
                              !field.value && "text-gray-400"
                            )}
                          >
                            <SelectValue placeholder="Select a force" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>Forces</SelectLabel>
                              <SelectItem value="none">None</SelectItem>
                              <SelectItem value="static">Static</SelectItem>
                              <SelectItem value="pull">Pull</SelectItem>
                              <SelectItem value="push">Push</SelectItem>
                            </SelectGroup>
                          </SelectContent>
                        </Select>
                      </FormControl>
                      <FormDescription>
                        What type of force does this exercise involve? (e.g.,
                        static, pull, push)
                      </FormDescription>
                      <FormMessage />
                    </div>
                  </FormItem>
                )}
              />
              {/* Level Field */}
              <FormField
                control={form.control}
                name="level"
                render={({ field }) => (
                  <FormItem>
                    <div className="grid gap-2">
                      <FormLabel>Exercise Level</FormLabel>
                      <FormControl>
                        <Select
                          onValueChange={(value) => {
                            field.onChange(value === "none" ? null : value);
                          }}
                          defaultValue={field.value || "none"}
                        >
                          <SelectTrigger
                            className={cn(
                              "w-full",
                              !field.value && "text-gray-400"
                            )}
                          >
                            <SelectValue placeholder="Select a level" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>Levels</SelectLabel>
                              <SelectItem value="none">None</SelectItem>
                              <SelectItem value="beginner">Beginner</SelectItem>
                              <SelectItem value="intermediate">
                                Intermediate
                              </SelectItem>
                              <SelectItem value="expert">Expert</SelectItem>
                            </SelectGroup>
                          </SelectContent>
                        </Select>
                      </FormControl>
                      <FormDescription>
                        What is the difficulty level of this exercise? (e.g.,
                        beginner, intermediate, expert)
                      </FormDescription>
                      <FormMessage />
                    </div>
                  </FormItem>
                )}
              />
              {/* Mechanic Field */}
              <FormField
                control={form.control}
                name="mechanic"
                render={({ field }) => (
                  <FormItem>
                    <div className="grid gap-2">
                      <FormLabel>Exercise Mechanic</FormLabel>
                      <FormControl>
                        <Select
                          onValueChange={(value) => {
                            field.onChange(value === "none" ? null : value);
                          }}
                          defaultValue={field.value || "none"}
                        >
                          <SelectTrigger
                            className={cn(
                              "w-full",
                              !field.value && "text-gray-400"
                            )}
                          >
                            <SelectValue placeholder="Select a level" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>Mechanics</SelectLabel>
                              <SelectItem value="none">None</SelectItem>
                              <SelectItem value="isolation">
                                Isolation
                              </SelectItem>
                              <SelectItem value="compound">Compound</SelectItem>
                            </SelectGroup>
                          </SelectContent>
                        </Select>
                      </FormControl>
                      <FormDescription>
                        What type of movement does this exercise use? (e.g.,
                        isolation, compound)
                      </FormDescription>
                      <FormMessage />
                    </div>
                  </FormItem>
                )}
              />
              {/* Equipment Field */}
              <FormField
                control={form.control}
                name="equipment"
                render={({ field }) => (
                  <FormItem>
                    <div className="grid gap-2">
                      <FormLabel>Exercise Equipment</FormLabel>
                      <FormControl>
                        <Select
                          onValueChange={(value) => {
                            field.onChange(value === "none" ? null : value);
                          }}
                          defaultValue={field.value || "none"}
                        >
                          <SelectTrigger
                            className={cn(
                              "w-full",
                              !field.value && "text-gray-400"
                            )}
                          >
                            <SelectValue placeholder="Select the equipment used" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>Mechanics</SelectLabel>
                              <SelectItem value="none">None</SelectItem>
                              <SelectItem value="medicine ball">
                                medicine ball
                              </SelectItem>
                              <SelectItem value="dumbbell">Dumbbell</SelectItem>
                              <SelectItem value="body only">
                                Body only
                              </SelectItem>
                              <SelectItem value="bands">Bands</SelectItem>
                              <SelectItem value="kettlebells">
                                Kettlebells
                              </SelectItem>
                              <SelectItem value="foam roll">
                                Foam roll
                              </SelectItem>
                              <SelectItem value="cable">Cable</SelectItem>
                              <SelectItem value="machine">Machine</SelectItem>
                              <SelectItem value="barbell">Barbell</SelectItem>
                              <SelectItem value="exercise ball">
                                Exercise ball
                              </SelectItem>
                              <SelectItem value="e-z curl bar">
                                E-Z Curl bar
                              </SelectItem>
                              <SelectItem value="other">Other</SelectItem>
                            </SelectGroup>
                          </SelectContent>
                        </Select>
                      </FormControl>
                      <FormDescription>
                        What type of movement does this exercise use? (e.g.,
                        isolation, compound)
                      </FormDescription>
                      <FormMessage />
                    </div>
                  </FormItem>
                )}
              />
              {/* Primary Muscles Field */}
              <FormField
                control={form.control}
                name="primaryMuscles"
                render={({ field }) => (
                  <FormItem>
                    <div className="grid gap-2">
                      <FormLabel>Primary Muscles</FormLabel>
                      <FormControl>
                        <Select
                          onValueChange={(value) => {
                            field.onChange(value === "none" ? null : value);
                          }}
                          defaultValue={field.value || "none"}
                        >
                          <SelectTrigger
                            className={cn(
                              "w-full",
                              !field.value && "text-gray-400"
                            )}
                          >
                            <SelectValue placeholder="Select primary muscles" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>Primary Muscles</SelectLabel>
                              <SelectItem value="none">None</SelectItem>
                              <SelectItem value="abdominals">
                                Abdominals
                              </SelectItem>
                              <SelectItem value="abductors">
                                Abductors
                              </SelectItem>
                              <SelectItem value="adductors">
                                Adductors
                              </SelectItem>
                              <SelectItem value="biceps">Biceps</SelectItem>
                              <SelectItem value="calves">Calves</SelectItem>
                              <SelectItem value="chest">Chest</SelectItem>
                              <SelectItem value="forearms">Forearms</SelectItem>
                              <SelectItem value="glutes">Glutes</SelectItem>
                              <SelectItem value="hamstrings">
                                Hamstrings
                              </SelectItem>
                              <SelectItem value="lats">Lats</SelectItem>
                              <SelectItem value="lower back">
                                Lower Back
                              </SelectItem>
                              <SelectItem value="middle back">
                                Middle Back
                              </SelectItem>
                              <SelectItem value="neck">Neck</SelectItem>
                              <SelectItem value="quadriceps">
                                Quadriceps
                              </SelectItem>
                              <SelectItem value="shoulders">
                                Shoulders
                              </SelectItem>
                              <SelectItem value="traps">Traps</SelectItem>
                              <SelectItem value="triceps">Triceps</SelectItem>
                            </SelectGroup>
                          </SelectContent>
                        </Select>
                      </FormControl>
                      <FormDescription>
                        Which muscles are primarily targeted by this exercise?
                      </FormDescription>
                      <FormMessage />
                    </div>
                  </FormItem>
                )}
              />
              {/* Secondary Muscles Field */}
              <FormField
                control={form.control}
                name="secondaryMuscles"
                render={({ field }) => (
                  <FormItem>
                    <div className="grid gap-2">
                      <FormLabel>Secondary Muscles</FormLabel>
                      <FormControl>
                        <Select
                          onValueChange={(value) => {
                            field.onChange(value === "none" ? null : value);
                          }}
                          defaultValue={field.value || "none"}
                        >
                          <SelectTrigger
                            className={cn(
                              "w-full",
                              !field.value && "text-gray-400"
                            )}
                          >
                            <SelectValue placeholder="Select secondary muscles" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>Secondary Muscles</SelectLabel>
                              <SelectItem value="none">None</SelectItem>
                              <SelectItem value="abdominals">
                                Abdominals
                              </SelectItem>
                              <SelectItem value="abductors">
                                Abductors
                              </SelectItem>
                              <SelectItem value="adductors">
                                Adductors
                              </SelectItem>
                              <SelectItem value="biceps">Biceps</SelectItem>
                              <SelectItem value="calves">Calves</SelectItem>
                              <SelectItem value="chest">Chest</SelectItem>
                              <SelectItem value="forearms">Forearms</SelectItem>
                              <SelectItem value="glutes">Glutes</SelectItem>
                              <SelectItem value="hamstrings">
                                Hamstrings
                              </SelectItem>
                              <SelectItem value="lats">Lats</SelectItem>
                              <SelectItem value="lower back">
                                Lower Back
                              </SelectItem>
                              <SelectItem value="middle back">
                                Middle Back
                              </SelectItem>
                              <SelectItem value="neck">Neck</SelectItem>
                              <SelectItem value="quadriceps">
                                Quadriceps
                              </SelectItem>
                              <SelectItem value="shoulders">
                                Shoulders
                              </SelectItem>
                              <SelectItem value="traps">Traps</SelectItem>
                              <SelectItem value="triceps">Triceps</SelectItem>
                            </SelectGroup>
                          </SelectContent>
                        </Select>
                      </FormControl>
                      <FormDescription>
                        Which muscles are secondarily involved in this exercise?
                      </FormDescription>
                      <FormMessage />
                    </div>
                  </FormItem>
                )}
              />

              {/* Category Field */}
              <FormField
                control={form.control}
                name="category"
                render={({ field }) => (
                  <FormItem>
                    <div className="grid gap-2">
                      <FormLabel>Exercise Category</FormLabel>
                      <FormControl>
                        <Select
                          onValueChange={field.onChange}
                          defaultValue={field.value}
                        >
                          <SelectTrigger className="w-full">
                            <SelectValue placeholder="Select a category" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>Categories</SelectLabel>
                              <SelectItem value="powerlifting">
                                Powerlifting
                              </SelectItem>
                              <SelectItem value="strength">Strength</SelectItem>
                              <SelectItem value="stretching">
                                Stretching
                              </SelectItem>
                              <SelectItem value="cardio">Cardio</SelectItem>
                              <SelectItem value="olympic_weightlifting">
                                Olympic Weightlifting
                              </SelectItem>
                              <SelectItem value="strongman">
                                Strongman
                              </SelectItem>
                              <SelectItem value="plyometrics">
                                Plyometrics
                              </SelectItem>
                            </SelectGroup>
                          </SelectContent>
                        </Select>
                      </FormControl>
                      <FormDescription>
                        What category does this exercise belong to? (e.g.,
                        strength, cardio, stretching)
                      </FormDescription>
                      <FormMessage />
                    </div>
                  </FormItem>
                )}
              />

              {/* Buttons */}
              <div className="flex justify-center gap-6">
                <Button type="button" variant="destructive">
                  Cancel
                </Button>
                <Button type="submit">Submit</Button>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
