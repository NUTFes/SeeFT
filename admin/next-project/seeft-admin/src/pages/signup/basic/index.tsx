import Image from "next/image";
import { SuccessButton } from "@components/common/SuccessButton";
import { FormInputText } from "@components/common/Form/FormInputText";
import { Form } from "@components/common/Form";
import { FormInputDropdown } from "@components/common/Form/FormInputDropdown";
import { FormInputRadio } from "@components/common/Form/FormInputRadio";
import { FailerButton } from "@components/common/FailerButton";
import { useRouter } from "next/router";
import { useState } from "react";
import { useForm } from "react-hook-form";

const SignupBasic = () => {
  const router = useRouter();
  const [formData, setFormData] = useState<{ name: string, password: string, passwordConfirmation: string }>({
    name: '',
    password: '',
    passwordConfirmation: '',
  });

  const {
    register,
    formState: { errors, isValid },
    getValues,
    handleSubmit,
  } = useForm<{ name: string, password: string, passwordConfirmation: string }>({
    mode: 'all',
  });

  const userDataHandler =
    (input: string) =>
      (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({ ...formData, [input]: e.target.value });
      };

  const postUser = () => {
    router.push({
      pathname: '/signup/detail',
      query: { name: formData.name, password: formData.password },
    });
  }

  return (
    <div className="h-screen !bg-gradient flex flex-col justify-center items-center">
      <div className="w-[480px]">
        <div className="flex flex-col items-center">
          {/* <Image src="/title.svg" alt="logo" width={436} height={160} /> */}
          <p className="text-2.25xl mb-10">Sign up</p>
          <p className="text-xl mb-10">基本情報</p>
        </div>
        <Form onSubmit={handleSubmit(postUser)}>
          <FormInputText
            text="名前"
            type="text"
            placeholder="技大 太郎"
            value={formData.name}
            onChange={userDataHandler('name')}
          ></FormInputText>
          <FormInputText
            text="パスワード"
            type="password"
            placeholder=""
            value={formData.password}
            onChange={userDataHandler('password')}
          ></FormInputText>
          <FormInputText
            text="パスワード確認"
            type="password"
            placeholder=""
            value={formData.passwordConfirmation}
            onChange={userDataHandler('passwordConfirmation')}
          ></FormInputText>
          <SuccessButton
            text="詳細情報の入力に進む"
          ></SuccessButton>
          <div className='mb-5'>
            <p className='text-red-500'>
              {errors.passwordConfirmation && errors.passwordConfirmation.type === 'correct' && (
                <p>パスワードが一致しません</p>
              )}
            </p>
          </div>
        </Form>
        <div className="flex flex-col items-center gap-y-4 m-6">
          <a href="../login">ログイン画面へ</a>
        </div>
      </div>
    </div>
  );
};

export default SignupBasic;
