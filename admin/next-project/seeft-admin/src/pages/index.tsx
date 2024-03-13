import FailerButton from "@components/common/FailerButton";
import SuccessButton from "@components/common/SuccessButton";
import React from "react";

export default function Home() {
  return (
    <div>
      <FailerButton
        text='リセット'
      >
      </FailerButton>
      <SuccessButton
        text='登録'
      >
      </SuccessButton>
    </div>
  );
}