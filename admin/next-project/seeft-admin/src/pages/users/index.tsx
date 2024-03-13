import clsx from 'clsx';
import Head from 'next/head';

import { get } from '@api/api_methods';
import { User, Grade, Department, Bureau } from "@type/common";
import MainLayout from '@components/layout/MainLayout';

interface Props {
  users: User[];
  grades: Grade[];
  departments: Department[];
  bureaus: Bureau[];
}

export const getServerSideProps = async () => {
  const getUserURL = process.env.SSR_API_URI + '/users';
  const getGradeURL = process.env.SSR_API_URI + '/grades';
  const getDepartmentURL = process.env.SSR_API_URI + '/departments';
  const getBureauURL = process.env.SSR_API_URI + '/bureaus';
  const userRes = await get(getUserURL);
  const gradeRes = await get(getGradeURL);
  const departmentRes = await get(getDepartmentURL);
  const bureauRes = await get(getBureauURL);

  return {
    props: {
      users: userRes,
      grades: gradeRes,
      departments: departmentRes,
      bureaus: bureauRes,
    },
  };
};

export default function Uesrs(props: Props) {
  const { users, grades, departments, bureaus } = props;

  return (
    <MainLayout>
      <div className='w-full h-full bg-surface-1/[0.5] flex-col p-8'>
        <div className='items-center'>
          test
        </div>
        <div className='p-5'>
          <table className='mb-5 w-full table-auto border-collapse'>
            <thead>
              <tr>
                <th className='border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                  <p className='text-center text-sm text-black-600'>所属局</p>
                </th>
                <th className='border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                  <p className='text-center text-sm text-black-600'>名前</p>
                </th>
                <th className='border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                  <p className='text-center text-sm text-black-600'>学年</p>
                </th>
                <th className='border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                  <p className='text-center text-sm text-black-600'>学科</p>
                </th>
                <th className='border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                  <p className='text-center text-sm text-black-600'>学籍番号</p>
                </th>
                <th className='border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                  <p className='text-center text-sm text-black-600'>電話番号</p>
                </th>
                <th className='border border-x-white-0 border-b-accent-1 border-t-white-0 py-3' />
                <th className='border border-x-white-0 border-b-accent-1 border-t-white-0 py-3' />
              </tr>
            </thead>
            <tbody className='border border-x-white-0 border-b-accent-1 border-t-white-0'>
              {users ? users.map((user: User, index) => (
                <tr key={user.id}>
                  <td
                    className={clsx(
                      'px-1 py-3',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b py-3',
                    )}
                  >
                    <p className='text-center text-sm text-black-600'>{bureaus.find((bureau: Bureau) => (bureau.id === user.bureauID))?.bureau}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-3',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b py-3',
                    )}
                  >
                    <p className='text-center text-sm text-black-600'>{user.name}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-3',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b py-3',
                    )}
                  >
                    <p className='text-center text-sm text-black-600'>{grades.length - 1 ? grades.find((grade: Grade) => (grade.id === user.gradeID))?.grade : "erorr"}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-3',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b py-3',
                    )}
                  >
                    <p className='text-center text-sm text-black-600'>{departments.length - 1 ? departments.find((department: Department) => (department.id === user.departmentID))?.department : "erorr"}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-3',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b py-3',
                    )}
                  >
                    <p className='text-center text-sm text-black-600'>{user.studentNumber}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-3',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b py-3',
                    )}
                  >
                    <p className='text-center text-sm text-black-600'>{user.tel}</p>
                  </td>
                </tr>
              )) : null}
            </tbody>
          </table>
          table
        </div>
      </div>
    </MainLayout >
  );
}