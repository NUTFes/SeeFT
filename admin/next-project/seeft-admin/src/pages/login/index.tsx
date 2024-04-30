import Image from "next/image";
import { SuccessButton } from "@components/common/SuccessButton";
import { FormInputText } from "@components/common/Form/FormInputText";
import { Form } from "@components/common/Form";

const Login = () => {
  return (
    <div className="h-screen !bg-gradient flex flex-col justify-center items-center">
      <div className="w-[480px]">
        <div className="flex flex-col items-center">
          <Image src="/title.svg" alt="logo" width={436} height={160} />
          <p className="text-2.25xl mb-10">Log in</p>
        </div>
        <Form>
          <FormInputText
            text="学籍番号"
            type="text"
            placeholder="000000"
          ></FormInputText>
          <FormInputText
            text="パスワード"
            type="password"
            placeholder=""
          ></FormInputText>
          <SuccessButton text="ログイン"></SuccessButton>
        </Form>
        <div className="flex flex-col items-center gap-y-4 m-6">
          <a href="">パスワードを変更する</a>
          <a href="../signup/basic">新規登録はこちら</a>
        </div>
      </div>
    </div>
  );
};

export default Login;
