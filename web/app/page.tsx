import { Button } from "@/components/ui/button";
import Image from "next/image";
import LoginPage from "./login/page";
import SignUpPage from "./signup/page";
import { CustomExerciseForm } from "@/components/custom-exercise";

export default function Home() {
  return (
    <>
      {/* <LoginPage /> */}
      <SignUpPage />
      <CustomExerciseForm />
    </>
  );
}
