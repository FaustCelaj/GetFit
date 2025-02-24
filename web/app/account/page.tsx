"use client";

import { BioDisplay } from "@/components/account/bio-display";
import { PersonalInformationDisplay } from "@/components/account/personal-information-display";
import { EmailDisplay } from "@/components/account/email-display";
import { PasswordDisplay } from "@/components/account/password-display";

import userImage from "@/public/userImage.jpg";

const userData = {
  picture: userImage,
  name: "Emily",
  lastName: "Perez",
  username: "EmilyPerez2001",
  email: "EmilyPerez@email.com",
  password: "123abc456",
  title: "Fitness Enthusiast ",
  bio: "Just strated lifting 1 year ago, feeling better than ever!",
  age: 24,
  country: "Canada",
  city: "Toronto",
};

export function AccountCard() {
  return (
    <div className="flex flex-col gap-y-6 ">
      <p className="text-lg font-semibold text-black">My Profile</p>

      <BioDisplay userData={userData} />
      <PersonalInformationDisplay userData={userData} />
      <EmailDisplay userData={userData} />
      <PasswordDisplay userData={userData} />
    </div>
  );
}
