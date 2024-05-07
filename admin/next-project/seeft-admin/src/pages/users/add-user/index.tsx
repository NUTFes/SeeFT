import React, { useState } from 'react';
import { useRouter } from 'next/router';

import { get } from '@api/api_methods';
import { User, Grade, Department, Bureau } from "@type/common";
import { Input, Select } from '@components/common';
import { post } from '@api/user';
import InformationPageLayout from '@components/layout/InformationPageLayout';

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

export default function Users(props: Props) {
  const { grades, departments, bureaus } = props;
  const router = useRouter();

  const [formData, setFormData] = useState<User>({
    id: 0,
    name: '',
    mail: '',
    gradeID: 1,
    departmentID: 1,
    bureauID: 1,
    roleID: 1,
    studentNumber: 0,
    tel: '',
    password: '',
  });

  const handler = (input: string) =>
    (e: React.ChangeEvent<HTMLSelectElement> | React.ChangeEvent<HTMLInputElement>) => {
      setFormData({ ...formData, [input]: e.target.value });
    }

  const addUserInformation = async (data: User) => {
    const addUserInformationUrl = process.env.CSR_API_URI + '/users';
    await post(addUserInformationUrl, data);
    router.push('/users');
  };

  return (
    <InformationPageLayout title='ユーザー登録' submitText='登録' onClick={() => { addUserInformation(formData); }}>
      <div className='flex w-full items-center'>
        <div className='flex w-1/4'>学籍番号</div>
        <div className='col-span-4 w-full'>
          <Input className='w-full' value={formData.studentNumber} onChange={handler('studentNumber')} />
        </div>
      </div>
      <div className='flex w-full items-center'>
        <div className='flex w-1/4 whitespace-nowrap'>パスワード</div>
        <input
          type='password'
          className='rounded-full border border-primary-1 px-4 py-2 col-span-2 w-full'
          onChange={handler('Password')}
        />
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