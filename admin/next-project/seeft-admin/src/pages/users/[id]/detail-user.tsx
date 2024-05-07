import React, { useState } from 'react';
import { useRouter } from 'next/router';

import { get } from '@api/api_methods';
import { User, Grade, Department, Bureau } from "@type/common";
import { Input, Select } from '@components/common';
import { put } from '@api/user';
import InformationPageLayout from '@components/layout/InformationPageLayout';

interface Props {
  grades: Grade[];
  departments: Department[];
  bureaus: Bureau[];
  user: User;
}

export const getServerSideProps = async (
  { params }: { params: { id: string } }) => {
  const userID = params.id;
  const getUserURL = process.env.SSR_API_URI + '/users/' + userID;
  const getGradeURL = process.env.SSR_API_URI + '/grades';
  const getDepartmentURL = process.env.SSR_API_URI + '/departments';
  const getBureauURL = process.env.SSR_API_URI + '/bureaus';
  const userRes = await get(getUserURL);
  const gradeRes = await get(getGradeURL);
  const departmentRes = await get(getDepartmentURL);
  const bureauRes = await get(getBureauURL);

  return {
    props: {
      grades: gradeRes,
      departments: departmentRes,
      bureaus: bureauRes,
      user: userRes,
    },
  };
};

export default function Users(props: Props) {
  const { grades, departments, bureaus, user } = props;
  const router = useRouter();

  const [formData, setFormData] = useState<User>({
    id: user.id,
    name: user.name,
    mail: user.mail,
    gradeID: user.gradeID,
    departmentID: user.departmentID,
    bureauID: user.bureauID,
    roleID: user.roleID,
    studentNumber: user.studentNumber,
    tel: user.tel,
    password: '',
  });

  const handler = (input: string) =>
    (e: React.ChangeEvent<HTMLSelectElement> | React.ChangeEvent<HTMLInputElement>) => {
      setFormData({ ...formData, [input]: e.target.value });
    }

  const updateUserInformation = async (data: User) => {
    const putUserInformationUrl = process.env.CSR_API_URI + '/users/' + data.id;
    await put(putUserInformationUrl, data);
    router.push('/users');
  };

  return (
    <InformationPageLayout title='ユーザー詳細' submitText='編集' onClick={() => { updateUserInformation(formData); }}>
      <div className='flex w-full items-center'>
        <div className='flex w-1/4'>学籍番号</div>
        <div className='col-span-4 w-full'>
          <Input className='w-full' value={formData.studentNumber} onChange={handler('studentNumber')} />
        </div>
      </div>
      <div className='flex w-full items-center'>
        <div className='flex w-1/4'>名前</div>
        <div className='col-span-4 w-full'>
          <Input className='w-full' value={formData.name} onChange={handler('name')} />
        </div>
      </div>
      <div className='flex w-full items-center'>
        <div className='flex w-1/4'>所属局</div>
        <div className='col-span-4 w-full'>
          <Select className='w-full' value={formData.bureauID} onChange={handler('bureauID')}>
            {bureaus.map((data) => (
              <option key={data.id} value={data.id}>
                {data.bureau}
              </option>
            ))}
          </Select>
        </div>
      </div>
      <div className='flex w-full items-center'>
        <div className='flex w-1/4'>課程</div>
        <div className='col-span-4 w-full'>
          <Select className='w-full' value={formData.departmentID} onChange={handler('departmentID')}>
            {departments.length > 1 ? departments.map((data) => (
              <option key={data.id} value={data.id}>
                {data.department}
              </option>
            )) : null}
          </Select>
        </div>
      </div>
      <div className='flex w-full items-center'>
        <div className='flex w-1/4'>学年</div>
        <div className='col-span-4 w-full'>
          <Select className='w-full' value={formData.gradeID} onChange={handler('gradeID')}>
            {grades.map((data) => (
              <option key={data.id} value={data.id}>
                {data.grade}
              </option>
            ))}
          </Select>
        </div>
      </div>
      <div className='flex w-full items-center'>
        <div className='flex w-1/4'>電話番号</div>
        <div className='col-span-4 w-full'>
          <Input className='w-full' value={formData.tel} onChange={handler('tel')} />
        </div>
      </div>
      <div className='flex w-full items-center'>
        <div className='flex w-1/4'>メールアドレス</div>
        <div className='col-span-4 w-full'>
          <Input className='w-full' value={formData.mail} onChange={handler('mail')} />
        </div>
      </div>
    </InformationPageLayout>
  );
}