"use client";

import LoginPage from "./login/page";
import SignUpPage from "./signup/page";
import { CustomExerciseForm } from "@/components/custom-exercise";
import { ExerciseList } from "./exercises/page";
import { ModeToggle } from "@/components/account/theme-toggle";
import { Account } from "./account/page";
import { WelcomeCard } from "@/components/home/welcome";
import { RoutineSlider } from "@/components/home/routine-slider";

const data = [
  {
    title: "Chest and Tri's",
    description: "Push day to start off the week, high intensity 💪",
    exercises: [
      { name: "Flat Bench Press", totalSets: 3, totalReps: 18 },
      { name: "Incline Bench Press", totalSets: 3, totalReps: 18 },
      { name: "Pec Deck Chest Flys", totalSets: 3, totalReps: 24 },
      { name: "Rope Tricep Extensions", totalSets: 3, totalReps: 30 },
      { name: "Skullcrushers", totalSets: 3, totalReps: 24 },
      { name: "Tricep Kickbacks", totalSets: 3, totalReps: 30 },
    ],
  },
  {
    title: "Back and Bi's",
    description: "Focus on pulling strength and bicep endurance 🏋️",
    exercises: [
      { name: "Deadlift", totalSets: 3, totalReps: 12 },
      { name: "Pull-Ups", totalSets: 3, totalReps: 15 },
      { name: "Lat Pulldown", totalSets: 3, totalReps: 18 },
      { name: "Seated Row", totalSets: 3, totalReps: 18 },
      { name: "Barbell Deadlifts", totalSets: 4, totalReps: 20 },
      { name: "Dumbbell Bicep Curls", totalSets: 3, totalReps: 24 },
      { name: "Hammer Curls", totalSets: 3, totalReps: 24 },
    ],
  },
  {
    title: "Leg Day",
    description: "Heavy compound lifts with quad and hamstring focus 🦵",
    exercises: [
      { name: "Squats", totalSets: 4, totalReps: 20 },
      { name: "Leg Press", totalSets: 4, totalReps: 24 },
      { name: "Romanian Deadlifts", totalSets: 3, totalReps: 18 },
      { name: "Leg Extensions", totalSets: 3, totalReps: 30 },
      { name: "Lying Hamstring Curls", totalSets: 3, totalReps: 30 },
      { name: "Calf Raises", totalSets: 4, totalReps: 40 },
      { name: "Leg Adduction", totalSets: 3, totalReps: 36 },
      { name: "Tibialis Extensions", totalSets: 3, totalReps: 30 },
    ],
  },
  {
    title: "Shoulders and Abs",
    description: "Build strong delts and a solid core ⚡",
    exercises: [
      { name: "Overhead Shoulder Press", totalSets: 3, totalReps: 18 },
      { name: "Lateral Raises", totalSets: 3, totalReps: 24 },
      { name: "Front Raises", totalSets: 3, totalReps: 24 },
      { name: "Face Pulls", totalSets: 3, totalReps: 30 },
      { name: "Hanging Leg Raises", totalSets: 3, totalReps: 20 },
    ],
  },
  {
    title: "Full Body Strength",
    description: "A mix of compound lifts for overall strength gains 🏋️‍♂️",
    exercises: [
      { name: "Deadlifts", totalSets: 4, totalReps: 20 },
      { name: "Squats", totalSets: 4, totalReps: 20 },
      { name: "Bench Press", totalSets: 3, totalReps: 18 },
      { name: "Farmer’s Walk (meters)", totalSets: 3, totalReps: 50 },
    ],
  },
  {
    title: "Hypertrophy Upper Body",
    description: "Higher volume training for muscle growth 🔥",
    exercises: [
      { name: "Dumbbell Bench Press", totalSets: 4, totalReps: 24 },
      { name: "Lat Pulldown", totalSets: 4, totalReps: 24 },
      { name: "Seated Shoulder Press", totalSets: 3, totalReps: 18 },
      { name: "Incline Dumbbell Press", totalSets: 3, totalReps: 18 },
      { name: "Face Pulls", totalSets: 3, totalReps: 30 },
      { name: "Preacher Curls", totalSets: 3, totalReps: 24 },
    ],
  },
  {
    title: "Hypertrophy Lower Body",
    description: "Volume-based leg training for definition and endurance 🦿",
    exercises: [
      { name: "Squats", totalSets: 4, totalReps: 24 },
      { name: "Bulgarian Split Squats", totalSets: 3, totalReps: 30 },
      { name: "Leg Press", totalSets: 3, totalReps: 24 },
      { name: "Romanian Deadlifts", totalSets: 3, totalReps: 18 },
      { name: "Calf Raises", totalSets: 4, totalReps: 40 },
      { name: "Seated Hamstring Curls", totalSets: 3, totalReps: 24 },
    ],
  },
  {
    title: "Active Recovery & Mobility",
    description:
      "A mix of stretching, core work, and light activity for recovery 🧘",
    exercises: [
      { name: "Foam Rolling", totalSets: 3, totalReps: 5 },
      { name: "Bird Dogs", totalSets: 3, totalReps: 20 },
      { name: "Glute Bridges", totalSets: 3, totalReps: 20 },
      { name: "Hanging Leg Raises", totalSets: 3, totalReps: 20 },
      { name: "Jump Rope"},
    ],
  },
];

export default function Home() {
  return (
    <div className="bg-slate-100 w-full h-full p-5 flex-row">
      <WelcomeCard username={"EmilyPerez2020"} />
      <RoutineSlider routines={data}/>
    </div>
  );
}
