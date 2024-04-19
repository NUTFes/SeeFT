import Image from "next/image";
import { SuccessButton } from "@components/common/SuccessButton";
import { FormInputText } from "@components/common/Form/FormInputText";
import { Form } from "@components/common/Form";
import { FormInputDropdown } from "@components/common/Form/FormInputDropdown";
import { FormInputRadio } from "@components/common/Form/FormInputRadio";
import { FailerButton } from "@components/common/FailerButton";

const SignupDetail = () => {
  const bureaus = [
    { value: 1, label: "総務局" },
    { value: 2, label: "企画局" },
    { value: 3, label: "渉外局" },
    { value: 4, label: "情報局" },
    { value: 5, label: "制作局" },
    { value: 6, label: "財務局" },
  ];

  const positions = [
    { value: 1, label: "どうしよう" },
    { value: 2, label: "これって" },
    { value: 3, label: "局によって" },
    { value: 4, label: "変わっちゃう" },
    { value: 5, label: "よなー" },
  ];

  const courses = [
    { value: 1, label: "電気" },
    { value: 2, label: "情経" },
    { value: 3, label: "機械" },
    { value: 4, label: "環社" },
    { value: 5, label: "生物" },
    { value: 6, label: "物質" },
    { value: 7, label: "物生" },
  ];

  const grades = [
    { value: 1, label: "B1" },
    { value: 2, label: "B2" },
    { value: 3, label: "B3" },
    { value: 4, label: "B4" },
    { value: 5, label: "M1" },
    { value: 6, label: "M2" },
  ];

  const curriculums = [
    { value: 1, label: "新カリ" },
    { value: 2, label: "旧カリ" },
  ];

  const sexes = [
    { value: 1, label: "男性" },
    { value: 2, label: "女性" },
    { value: 3, label: "その他" },
  ];

  return (
    <div className="h-screen !bg-gradient flex flex-col justify-center items-center">
      <div className="w-[480px]">
        <div className="flex flex-col items-center">
          {/* <Image src="/title.svg" alt="logo" width={436} height={160} /> */}
          <p className="text-2.25xl mb-10">Sign up</p>
          <p className="text-xl mb-10">詳細情報</p>
        </div>
        <Form>
          <FormInputText
            text="名前"
            type="text"
            placeholder="技大 太郎"
          ></FormInputText>
          <FormInputDropdown
            text="所属局"
            options={bureaus}
          ></FormInputDropdown>
          <FormInputDropdown
            text="役職"
            options={positions}
          ></FormInputDropdown>
          <FormInputDropdown text="課程" options={courses}></FormInputDropdown>
          <FormInputDropdown text="学年" options={grades}></FormInputDropdown>
          <FormInputText
            text="電話番号"
            type="text"
            placeholder="000-0000-0000"
          ></FormInputText>
          <FormInputRadio
            text="カリキュラム"
            options={curriculums}
          ></FormInputRadio>
          <FormInputRadio text="性別" options={sexes}></FormInputRadio>
          <div className="flex gap-x-6">
            <FailerButton text="戻る"></FailerButton>
            <SuccessButton text="登録"></SuccessButton>
          </div>
        </Form>
      </div>
    </div>
  );
};

export default SignupDetail;
