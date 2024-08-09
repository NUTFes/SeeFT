import Image from "next/image";
import { useForm } from 'react-hook-form';
import { SuccessButton } from "@components/common/SuccessButton";
import { FormInputText } from "@components/common/Form/FormInputText";
import { Form } from "@components/common/Form";
import { SignIn } from "@type/common";
import { useState } from "react";
import { useRecoilState } from "recoil";
import { authAtom, userAtom } from "@/store/atoms";
import Router from "next/router";
import { signIn } from "@api/signIn";
import { get_with_token } from "@api/api_methods";

const Login = () => {
  const [isSignInNow, setIsSignInNow] = useState<boolean>(false);
  const [, setAuth] = useRecoilState(authAtom);
  const [, setUser] = useRecoilState(userAtom);

  const {
    register,
    formState: { errors, isValid },
    handleSubmit,
  } = useForm<SignIn>({
    mode: 'all',
  });

  const SignIn = async (data: SignIn) => {
    setIsSignInNow(true);
    const signinUrl: string = process.env.CSR_API_URI + '/mail_auth/web_signin';
    const currentUserUrl: string = process.env.CSR_API_URI + '/current_user';

    const req = await signIn(signinUrl, data);
    const res = await req.json();
    const userRes = await get_with_token(currentUserUrl, res.accessToken);
    if (req.status === 200) {
      const authData = {
        isSignIn: true,
        accessToken: res.accessToken,
      };
      setAuth(authData);
      setUser(userRes);
      Router.push('/shifts');
    } else {
      alert(
        'ログインに失敗しました。学籍番号もしくはパスワードが間違っている可能性があります',
      );
      setIsSignInNow(false);
    }
  };

  return (
    <div className="h-screen !bg-gradient flex flex-col justify-center items-center">
      <div className="w-[480px]">
        <div className="flex flex-col items-center">
          <Image src="/title.svg" alt="logo" width={436} height={160} />
          <p className="text-2.25xl mb-10">Log in</p>
        </div>
        <Form onSubmit={handleSubmit(SignIn)}>
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
                required: 'パスワードは必須です',
                minLength: {
                  value: 6,
                  message: 'パスワードは6文字以上で入力してください',
                },
              })}
            />
          </div>
          <div className='mb-5'>
            <p className='text-red-500'>{errors.studentNumber && errors.studentNumber.message}</p>
            <p className='text-red-500'>{errors.password && errors.password.message}</p>
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
