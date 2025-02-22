import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import { UpdateUsername } from "@/components/account/update-username";
import { UpdateEmail } from "@/components/account/update-email";

const UserData = {
    name: "Faust Celaj",
    username: "Faust2025",
    email: "faustemail@email.com"
    

}

export function AccountCard() {
  return (
    <Card className="p-3">
      {/* <UpdateUsername /> */}
      {/* <UpdateEmail /> */}
    </Card>
  );
}
