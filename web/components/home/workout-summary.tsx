"use client";

import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { CircleFadingPlus } from "lucide-react";
import { Button } from "../ui/button";

export function WorkoutSummaryCard() {
  return (
    <Card className="p-4 ">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-4">
        <div>
          <p className="text-xs text-gray-600">
            Thursday Night Workout (default)
          </p>
          <p className="text-md font-semibold">Routine - Legs and Abs</p>
        </div>
        <p className="text-xs text-gray-600 mt-1 sm:mt-0">Thursday Feb. 21</p>
      </div>
      <div className="grid grid-cols-2 lg:grid-cols-3 gap-4">
        <StatCard title="Duration" value="45 min" />
        <StatCard title="Total Reps" value="150" />
        <StatCard title="Total Weight Moved" value="12,500 lbs" />
        <Button variant="secondary">
            Expand <CircleFadingPlus className="h-4 w-4" />
          </Button>
      </div>
    </Card>
  );
}

function StatCard({ title, value }: { title: string; value: string | number }) {
  return (
    <Card className="bg-gray-100 text-center p-3 shadow-sm">
      <CardHeader className="p-2">
        <CardTitle className="text-xs text-gray-600">{title}</CardTitle>
      </CardHeader>
      <CardContent className="p-1">
        <p className="text-lg font-bold">{value}</p>
      </CardContent>
    </Card>
  );
}
