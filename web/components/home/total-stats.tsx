"use client";

import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";

interface Stats {
  workoutsLogged: number;
  totalTime: number;
  totalSets: number;
  totalUniqueExercise: number;
  totalWeight: number;
  totalReps: number;
}

interface TotalStatsCardProps {
  stats: Stats;
}

export function TotalStatsCard({ stats }: TotalStatsCardProps) {
  return (
    <Card className="bg-white text-black p-4 rounded-xl shadow-sm">
      <CardHeader className="flex flex-row justify-between items-start sm:items-center p-2">
        <div>
          <p className="text-xs text-gray-600">Recent View</p>
          <p className="text-sm font-semibold">15 - 22 Feb 2025</p>
        </div>
        <Button variant="outline" size="sm" className= "mt-2 sm:mt-0">
          Select Date
        </Button>
      </CardHeader>

      <Separator className="my-2 bg-gray-200" />

      <CardContent className="p-2">
        <div className="grid grid-cols-2 lg:grid-cols-3 gap-4">
          <StatCard title="Workouts Logged" value={stats.workoutsLogged} />
          <StatCard
            title="Total Duration"
            value={`${Math.floor(stats.totalTime / 60)}h ${stats.totalTime % 60}m`}
          />
          <StatCard title="Total Sets" value={stats.totalSets} />
          <StatCard title="Unique Exercises" value={stats.totalUniqueExercise} />
          <StatCard title="Total Reps" value={stats.totalReps} />
          <StatCard title="Weight Moved" value={`${(stats.totalWeight).toLocaleString()} lbs`} />
        </div>
      </CardContent>
    </Card>
  );
}

function StatCard({ title, value }: { title: string; value: string | number }) {
  return (
    <Card className="bg-gray-100 text-black text-center p-3 rounded-lg shadow-sm">
      <CardHeader className="p-2">
        <CardTitle className="text-xs text-gray-600">{title}</CardTitle>
      </CardHeader>
      <CardContent className="p-1">
        <p className="text-lg xl:text-xl font-bold">{value}</p>
      </CardContent>
    </Card>
  );
}