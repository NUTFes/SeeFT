import { FailerButton } from "@components/common/FailerButton";
import { SuccessButton } from "@components/common/SuccessButton";
import React, { useState } from "react";
import Login from "./login";
import SignupBasic from "./signup/basic";

export default function Home() {
  const [isMember, setIsMember] = useState(true);
  const cardContent = (isMember: boolean) => {
    if (isMember) {
      return (
        <>
          <Login />
        </>
      );
    } else {
      return (
        <>
          <SignupBasic />
        </>
      );
    }
  };
  return (
    <div>
      {cardContent(isMember)}
    </div>
  );
}