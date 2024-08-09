import Image from "next/image";
import { SuccessButton } from "@components/common/SuccessButton";
import { FormInputText } from "@components/common/Form/FormInputText";
import { Form } from "@components/common/Form";
import { FormInputDropdown } from "@components/common/Form/FormInputDropdown";
import { FormInputRadio } from "@components/common/Form/FormInputRadio";
import { FailerButton } from "@components/common/FailerButton";
import { authAtom, userAtom } from "@/store/atoms";
import { useRecoilState } from "recoil";
import { useState } from "react";
import { Bureau, Department, Grade, Shift, Date, Time, User, Weather } from "@type/common";
import { get, get_with_token } from "@api/api_methods";
import { useRouter } from "next/router";
import { signUp } from "@api/signUp";
import { useForm } from "react-hook-form";
import { post } from "@api/user";
import { post as shiftPost } from '@api/shift';
import { YearItem } from "@constants/yearItem";
import { DateItem } from "@constants/dateItem";
import { TimeItems } from "@constants/timeItem";
import { WeatherItem } from "@constants/weatherItem";

interface Props {
  grades: Grade[];
  departments: Department[];
  bureaus: Bureau[];
}

export const getServerSideProps = async () => {
  const getGradeURL = process.env.SSR_API_URI + '/grades';
  const getDepartmentURL = process.env.SSR_API_URI + '/departments';
  const getBureauURL = process.env.SSR_API_URI + '/bureaus';
  const gradeRes = await get(getGradeURL);
  const departmentRes = await get(getDepartmentURL);
  const bureauRes = await get(getBureauURL);

  return {
    props: {
      grades: gradeRes,
      departments: departmentRes,
      bureaus: bureauRes,
    },
  };
};

const SignupDetail = (props: Props) => {
  const positions = [
    { value: 1, label: "どうしよう" },
    { value: 2, label: "これって" },
    { value: 3, label: "局によって" },
    { value: 4, label: "変わっちゃう" },
    { value: 5, label: "よなー" },
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

  const iniiShiftData: Shift = {
    id: 0,
    taskID: 1,
    userID: 0,
    yearID: YearItem[YearItem.length - 1].id,
    dateID: 1,
    timeID: 1,
    weatherID: 1,
    isAttendance: false
  };

  const { grades, departments, bureaus } = props;
  const router = useRouter();
  const { name, password } = router.query;
  const [, setAuth] = useRecoilState(authAtom);
  const [, setUser] = useRecoilState(userAtom);

  // 新規登録中フラグ
  const [isSignUpNow, setIsSignUpNow] = useState<boolean>(false);

  const [postUserData, setPostUserData] = useState<User>({
    id: 0,
    name: name ? (Array.isArray(name) ? name[0] : name) : '',
    mail: '',
    gradeID: 1,
    departmentID: 1,
    bureauID: 1,
    roleID: 1,
    studentNumber: 0,
    tel: '',
    password: password ? (Array.isArray(password) ? password[0] : password) : '',
  });

  const {
    register,
    formState: { errors, isValid },
    getValues,
    handleSubmit,
  } = useForm<User>({
    mode: 'all',
  });

  const userDataHandler =
    (input: string) =>
      (e: React.ChangeEvent<HTMLSelectElement> | React.ChangeEvent<HTMLInputElement>) => {
        if (input === 'studentNumber') {
          setPostUserData({ ...postUserData, [input]: Number(e.target.value) });
        }
        else {
          setPostUserData({ ...postUserData, [input]: e.target.value });
        }
      };

  const SignUp = async (data: User) => {
    setIsSignUpNow(true);
    const userUrl: string = process.env.CSR_API_URI + '/users';
    const signUpUrl: string = process.env.CSR_API_URI + '/mail_auth/web_signup';
    // userのpost時のResに登録したデータが返ってこないので以下で用意
    const getRes = await get(userUrl);
    const userID: number = getRes[getRes.length - 1].id + 1;
    // signIn には登録したuserのIDが必要なので先にUserをpost
    // await post(userUrl, postUserData);
    // signUp
    const req = await signUp(signUpUrl, postUserData);
    const res = await req.json();
    // state用のuserのデータ
    const userData: User = {
      id: userID,
      name: postUserData.name,
      mail: postUserData.mail,
      gradeID: Number(postUserData.gradeID),
      departmentID: Number(postUserData.departmentID),
      bureauID: Number(postUserData.bureauID),
      roleID: postUserData.roleID,
      studentNumber: Number(postUserData.studentNumber),
      tel: postUserData.tel,
      password: postUserData.password,
    };
    console.log(userData);
    if (req.status === 200) {
      // state用のauthのデータ
      const authData = {
        isSignIn: true,
        accessToken: res.accessToken,
      };
      const initShiftInformationUrl = process.env.CSR_API_URI + '/shifts-admin';
      iniiShiftData['userID'] = userID;
      DateItem.map((date: Date) => {
        iniiShiftData['dateID'] = date.id;
        TimeItems.map((time: Time) => {
          iniiShiftData['timeID'] = time.id
          WeatherItem.map(async (weather: Weather) => {
            iniiShiftData['weatherID'] = weather.id;
            await shiftPost(initShiftInformationUrl, iniiShiftData);
          })
        })
      });
      setAuth(authData);
      setUser(userData);
      router.push('/shifts');
    } else {
      alert(
        '新規登録に失敗しました。メールアドレスがすでに登録されている可能性があります',
      );
      setIsSignUpNow(false);
    }
  };

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
            value={postUserData.name}
            onChange={userDataHandler('name')}
          ></FormInputText>
          <FormInputText
            text="学籍番号"
            type="number"
            placeholder="xxxxxxxx"
            value={String(postUserData.studentNumber)}
            onChange={userDataHandler('studentNumber')}
          ></FormInputText>
          <FormInputDropdown text="所属局" onChange={userDataHandler('bureauID')}>
            {bureaus.map((bureau) => (
              <option key={bureau.id} value={bureau.id}>
                {bureau.bureau}
              </option>
            ))}
          </FormInputDropdown>
          {/* <FormInputDropdown text="役職">
            {positions.map((position) => (
              <option key={position.value} value={position.value}>
                {position.label}
              </option>
            ))}
          </FormInputDropdown> */}
          <FormInputDropdown text="課程" onChange={userDataHandler('departmentID')}>
            {departments.map((department) => (
              <option key={department.id} value={department.id} >
                {department.department}
              </option>
            ))}
          </FormInputDropdown>
          <FormInputDropdown text="学年" onChange={userDataHandler('gradeID')}>
            {grades.map((grade) => (
              <option key={grade.id} value={grade.id}>
                {grade.grade}
              </option>
            ))}
          </FormInputDropdown>
          <FormInputText
            text="電話番号"
            type="text"
            placeholder="000-0000-0000"
            value={postUserData.tel}
            onChange={userDataHandler('tel')}
          ></FormInputText>
          {/* <FormInputRadio
            text="カリキュラム"
            options={curriculums}
          ></FormInputRadio>
          <FormInputRadio text="性別" options={sexes} onChange={() => userDataHandler('sexID')} ></FormInputRadio> */}
          <div className="flex gap-x-6">
            <FailerButton text="戻る" onClick={() => { router.push('/signup/basic') }}></FailerButton>
            <SuccessButton text="登録" onClick={handleSubmit(SignUp)}></SuccessButton>
          </div>
        </Form>
      </div>
    </div >
  );
};

export default SignupDetail;
