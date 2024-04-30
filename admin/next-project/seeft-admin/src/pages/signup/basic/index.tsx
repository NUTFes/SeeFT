import Image from "next/image";
import { SuccessButton } from "@components/common/SuccessButton";
import { FormInputText } from "@components/common/Form/FormInputText";
import { Form } from "@components/common/Form";
import { FormInputDropdown } from "@components/common/Form/FormInputDropdown";
import { FormInputRadio } from "@components/common/Form/FormInputRadio";
import { FailerButton } from "@components/common/FailerButton";
import { useRouter } from "next/router";

const SignupBasic = () => {
  return (
    <div className="h-screen !bg-gradient flex flex-col justify-center items-center">
      <div className="w-[480px]">
        <div className="flex flex-col items-center">
          {/* <Image src="/title.svg" alt="logo" width={436} height={160} /> */}
          <p className="text-2.25xl mb-10">Sign up</p>
          <p className="text-xl mb-10">基本情報</p>
        </div>
        <Form>
          <FormInputText
            text="名前"
            type="text"
            placeholder="技大 太郎"
          ></FormInputText>
          <FormInputText
            text="パスワード"
            type="password"
            placeholder=""
          ></FormInputText>
          <FormInputText
            text="パスワード確認"
            type="password"
            placeholder=""
          ></FormInputText>
          <SuccessButton
            text="詳細情報の入力に進む"
          ></SuccessButton>
        </Form>
        <div className="flex flex-col items-center gap-y-4 m-6">
          <a href="../login">ログイン画面へ</a>
        </div>
      </div>
    </div>
  );
};

export default SignupBasic;
