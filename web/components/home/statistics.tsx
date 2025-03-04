"use client";

import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { WorkoutSummaryCard } from "./workout-summary";
import { TotalStatsCard } from "./total-stats";

interface Stats {
  workoutsLogged: number;
  totalTime: number;
  totalSets: number;
  totalUniqueExercise: number;
  totalWeight: number;
  totalReps: number;
}

interface StatisticsProps {
  stats: Stats;
}

export function StatisticsCard({ stats }: StatisticsProps) {
  return (
    <Card className=" h-fit shadow-md">
      <CardHeader className="pb-1 pt-4 px-6">
        <CardTitle className="text-xl font-bold">Your Statistics</CardTitle>
      </CardHeader>
      <Separator className="bg-gray-200" />
      <CardContent className="grid grid-cols-1 gap-6 p-6">
        <div>
          <CardDescription className="mb-3 text-gray-600">
            Previous Workout
          </CardDescription>
          <WorkoutSummaryCard />
        </div>
        
        <div>
          <CardDescription className="mb-3 text-gray-600">
            Your Total Stats
          </CardDescription>
          <TotalStatsCard stats={stats} />
        </div>
      </CardContent>
    </Card>
  );
}