"use client";

import React, { useEffect, useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { ExerciseExpanded } from "@/components/exercise-expanded";
import { ExerciseType } from "@/types/exercise";

interface ExerciseModalProps {
  isOpen: boolean;
  onClose: () => void;
  exercise: ExerciseType | null;
}

export function ExerciseModal({ isOpen, onClose, exercise }: ExerciseModalProps) {
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) return null; // Prevent hydration mismatch

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>{exercise?.name}</DialogTitle>
        </DialogHeader>
        <ExerciseExpanded exercise={exercise} />
      </DialogContent>
    </Dialog>
  );
}