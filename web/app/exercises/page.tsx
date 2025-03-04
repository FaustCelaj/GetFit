"use client";

import React, { useState } from "react";
import { ExerciseSmall } from "@/components/exersice-small";
import { ExerciseModal } from "@/components/exercise-wrapper";
import { ExerciseType } from "@/types/exercise";

// Mock Data (Replace with API call)
const exercises: ExerciseType[] = [
  {
    _id: "1",
    name: "3/4 Sit-Up",
    force: "pull",
    level: "beginner",
    mechanic: "compound",
    equipment: "body only",
    primaryMuscles: ["abdominals"],
    secondaryMuscles: ["forearms", "biceps"],
    instructions: [
      "Lie down on the floor and secure your feet.",
      "Place your hands behind your head.",
      "Flex your hips and spine to raise your torso.",
      "At the top, your torso should be perpendicular to the ground.",
      "Reverse the motion, going only ¾ of the way down.",
      "Repeat for the recommended reps.",
    ],
    category: "strength",
  },
  {
    _id: "2",
    name: "Axle Deadlift",
    force: "pull",
    level: "intermediate",
    mechanic: "compound",
    equipment: "other",
    primaryMuscles: ["lower back"],
    secondaryMuscles: [
      "forearms",
      "glutes",
      "hamstrings",
      "middle back",
      "quadriceps",
      "traps",
    ],
    instructions: [
      "Approach the bar so that it is centered over your feet. You feet should be about hip width apart. Bend at the hip to grip the bar at shoulder width, allowing your shoulder blades to protract. Typically, you would use an over/under grip.",
      "With your feet and your grip set, take a big breath and then lower your hips and flex the knees until your shins contact the bar. Look forward with your head, keep your chest up and your back arched, and begin driving through the heels to move the weight upward.",
      "After the bar passes the knees, aggressively pull the bar back, pulling your shoulder blades together as you drive your hips forward into the bar.",
      "Lower the bar by bending at the hips and guiding it to the floor.",
    ],
    category: "strongman",
  },
  {
    _id: "3",
    name: "Barbell Shrug",
    force: "pull",
    level: "beginner",
    mechanic: "isolation",
    equipment: "barbell",
    primaryMuscles: ["traps"],
    secondaryMuscles: [],
    instructions: [
      "Stand up straight with your feet at shoulder width as you hold a barbell with both hands in front of you using a pronated grip (palms facing the thighs). Tip: Your hands should be a little wider than shoulder width apart. You can use wrist wraps for this exercise for a better grip. This will be your starting position.",
      "Raise your shoulders up as far as you can go as you breathe out and hold the contraction for a second. Tip: Refrain from trying to lift the barbell by using your biceps.",
      "Slowly return to the starting position as you breathe in.",
      "Repeat for the recommended amount of repetitions.",
    ],
    category: "strength",
  },
];

export default function ExerciseList() {
  const [selectedExercise, setSelectedExercise] = useState<ExerciseType | null>(
    null
  );
  const [modalOpen, setModalOpen] = useState(false);

  const handleViewMore = (exercise: ExerciseType) => {
    setSelectedExercise(exercise);
    setModalOpen(true);
  };

  return (
    <div className="flex flex-col items-center space-y-4">
      {exercises.map((exercise) => (
        <ExerciseSmall
          key={exercise._id}
          exercise={exercise}
          onViewMore={() => handleViewMore(exercise)}
        />
      ))}

      <ExerciseModal
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        exercise={selectedExercise}
      />
    </div>
  );
}
