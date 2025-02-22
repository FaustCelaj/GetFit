import LoginPage from "./login/page";
import SignUpPage from "./signup/page";
import { CustomExerciseForm } from "@/components/custom-exercise";
import { ExerciseList } from "./exercises/page";
import { ModeToggle } from "@/components/account/theme-toggle";
import { AccountCard } from "./account/page";



export default function Home() {
  return (
    <>
      {/* <LoginPage />
      <SignUpPage /> */}
      {/* <ModeToggle/> */}
      {/* <CustomExerciseForm /> */}
      {/* <ExerciseList/> */}
      <AccountCard/>
    </>
  );
}
