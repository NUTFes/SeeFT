import Image from "next/image";
import { useForm } from 'react-hook-form';
import { SuccessButton } from "@components/common/SuccessButton";
import { FormInputText } from "@components/common/Form/FormInputText";
import { Form } from "@components/common/Form";
import { LoginType } from "@type/common";
import { useState } from "react";

const Login = () => {
  const [isSignInNow, setIsSignInNow] = useState<boolean>(false);

  const {
    register,
    formState: { errors, isValid },
    handleSubmit,
  } = useForm<LoginType>({
    mode: 'all',
  });

  return (
    <div className="h-screen !bg-gradient flex flex-col justify-center items-center">
      <div className="w-[480px]">
        <div className="flex flex-col items-center">
          <Image src="/title.svg" alt="logo" width={436} height={160} />
          <p className="text-2.25xl mb-10">Log in</p>
        </div>
        <Form>
          <div className="flex w-full h-10 items-center">
            <p className="w-40">学籍番号</p>
            <input
              type="text"
              placeholder="000000"
              className="flex-grow h-full bg-transparent border-b border-solid border-accent-1"
              {...register("studentNumber", {
                required: '入力は必須です',
                pattern: {
                  value: /^\d{8}$/,
                  message: '8桁の学籍番号を入力してください',
                },
              })}
            />
          </div>
          <div className="flex w-full h-10 items-center">
            <p className="w-40">パスワード</p>
            <input
              type="password"
              placeholder=""
              className="flex-grow h-full bg-transparent border-b border-solid border-accent-1"
              {...register("password", {
                required: '入力は必須です',
                pattern: {
                  value: /^.$/,
                  message: '所定の形式で入力してください',
                },
              })}
            />
          </div>
          {isSignInNow ? (
            <SuccessButton text='ログイン中' />
          ) : (
            <SuccessButton disabled={!isValid} text="ログイン" />
          )}
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
