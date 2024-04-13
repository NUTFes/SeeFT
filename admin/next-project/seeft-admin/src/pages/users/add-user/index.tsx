import clsx from 'clsx';
import Head from 'next/head';

import { get } from '@api/api_methods';
import { User, Grade, Department, Bureau } from "@type/common";
import MainLayout from '@components/layout/MainLayout';
import { MdEdit, MdDeleteForever } from "react-icons/md";
import Button from '@components/common/Button';
import Input from '@components/common/Input';
import Select from '@components/common/Select';

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

export default function Uesrs(props: Props) {
  const { grades, departments, bureaus } = props;

  return (
    <MainLayout>
      <div className='mx-auto relative md:w-1/2 h-full bg-white-0 p-8'>
        <div className=''>
          <div className='mx-auto w-fit text-xl text-emphasis mb-8'>
            ユーザー登録
          </div>
          <div className='mb-10 flex flex-col gap-3'>
            <div className='my-4 grid grid-cols-4 flex-col items-center justify-items-end gap-5 text-base text-emphasis'>
              <p>学籍番号</p>
              <div className='col-span-4 w-full'>
                <Input className='w-full' />
              </div>
              <p>パスワード</p>
              <div className='col-span-4 w-full'>
                <Input className='w-full' />
              </div>
              <p>名前</p>
              <div className='col-span-4 w-full'>
                <Input className='w-full' />
              </div>
              <p>所属局</p>
              <div className='col-span-4 w-full'>
                <Select className='w-full'>
                  {bureaus.map((data) => (
                    <option key={data.id} value={data.id}>
                      {data.bureau}
                    </option>
                  ))}
                </Select>
              </div>
              <p>課程</p>
              <div className='col-span-4 w-full'>
                <Select className='w-full'>
                  {departments.length > 1 ? departments.map((data) => (
                    <option key={data.id} value={data.id}>
                      {data.department}
                    </option>
                  )) : null}
                </Select>
              </div>
              <p>学年</p>
              <div className='col-span-4 w-full'>
                <Select className='w-full'>
                  {grades.map((data) => (
                    <option key={data.id} value={data.id}>
                      {data.grade}
                    </option>
                  ))}
                </Select>
              </div>
              <p>電話番号</p>
              <div className='col-span-4 w-full'>
                <Input className='w-full' />
              </div>
              <p>メールアドレス</p>
              <div className='col-span-4 w-full'>
                <Input className='w-full' />
              </div>
            </div>
          </div>
          <div className='mx-auto w-fit text-emphasis mb-8'>
            <Button>
              登録
            </Button>
          </div>
        </div>
      </div >
    </MainLayout >
  );
}