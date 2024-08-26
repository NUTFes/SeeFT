import clsx from 'clsx';
import { useRouter } from 'next/router';

import { get } from '@api/api_methods';
import { User, Grade, Department, Bureau } from "@type/common";
import MainLayout from '@components/layout/MainLayout';
import Button from '@components/common/Button';
import { destroy } from '@api/user';
import { DeleteButton, EditButton, Select } from '@components/common';
import ListPageLayout from '@components/layout/ListPageLayout';
import { useMemo, useState } from 'react';

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

export default function Users(props: Props) {
  const { users, grades, departments, bureaus } = props;
  const [filteredBureau, setFilteredBureau] = useState<number>(0);
  const router = useRouter();

  const addUserPageRouter = () => {
    router.push('users/add-user');
  }

  const userDetailPageRouter = (user: User) => {
    router.push('users/' + user.id + '/detail-user');
  }

  const destroyUserInformation = async (data: User) => {
    const destroyUserInformationUrl = process.env.CSR_API_URI + '/users';
    await destroy(destroyUserInformationUrl, data);
    router.reload();
  };

  const filterBureauHandler = () =>
    (e: React.ChangeEvent<HTMLSelectElement>) => {
      setFilteredBureau(Number(e.target.value));
    };

  const filteredUsers = useMemo(() => {
    return filteredBureau === 0 ? users?.sort((a: User, b: User) => a.bureauID - b.bureauID)
      : users.filter((user: User) => (
        user.bureauID === filteredBureau
      ))
  }, [filteredBureau]);
  
  return (
    <ListPageLayout title='ユーザー一覧'>
      <div className='items-center'>
        <div className='w-full flex justify-center items-center gap-6 p-1 '>
          <div className='w-1/6 ml-auto'>
            <Select className="w-full" value={filteredBureau} onChange={filterBureauHandler()}>
              <option key={0} value={0}>全局</option>
              {bureaus.map((data) => (
                <option key={data.id} value={data.id}>
                  {data.bureau}
                </option>
              ))}
            </Select>
          </div>
          <div className='text-right pr-4 ml-auto'>
            <Button className='bg-surface-2 border-accent-2 text-right text-emphasis pr-4 hover:bg-surface-1' onClick={addUserPageRouter}>
              ユーザー追加
            </Button>
          </div>
        </div>
      </div>
      <div className='p-5'>
        <table className='mb-5 w-full table-auto border-collapse'>
          <thead>
            <tr>
              <th className='w-1/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                <p className='text-center text-sm text-emphasis'>所属局</p>
              </th>
              <th className='w-2/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                <p className='text-center text-sm text-emphasis'>名前</p>
              </th>
              <th className='w-1/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                <p className='text-center text-sm text-emphasis'>学年</p>
              </th>
              <th className='w-2/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                <p className='text-center text-sm text-emphasis'>学科</p>
              </th>
              <th className='w-2/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                <p className='text-center text-sm text-emphasis'>学籍番号</p>
              </th>
              <th className='w-2/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                <p className='text-center text-sm text-emphasis'>電話番号</p>
              </th>
              <th className='w-1/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3' />
              <th className='w-1/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3' />
            </tr>
          </thead>
          <tbody className='border border-x-white-0 border-b-accent-1 border-t-white-0'>
            {filteredUsers ? filteredUsers.map((user: User, index) => (
              <tr key={user.id}>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === filteredUsers.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <p className='text-center text-sm text-emphasis'>{bureaus.find((bureau: Bureau) => (bureau.id === user.bureauID))?.bureau}</p>
                </td>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === filteredUsers.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <p className='text-center text-sm text-emphasis'>{user.name}</p>
                </td>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === filteredUsers.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <p className='text-center text-sm text-emphasis'>{grades.length - 1 ? grades.find((grade: Grade) => (grade.id === user.gradeID))?.grade : "erorr"}</p>
                </td>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === filteredUsers.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <p className='text-center text-sm text-emphasis'>{departments.length - 1 ? departments.find((department: Department) => (department.id === user.departmentID))?.department : "erorr"}</p>
                </td>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === filteredUsers.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <p className='text-center text-sm text-emphasis'>{user.studentNumber}</p>
                </td>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === filteredUsers.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <p className='text-center text-sm text-emphasis'>{user.tel}</p>
                </td>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === filteredUsers.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <EditButton onClick={() => { userDetailPageRouter(user) }}>
                    編集
                  </EditButton>
                </td>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === users.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <DeleteButton onClick={() => { destroyUserInformation(user) }} >
                    削除
                  </DeleteButton>
                </td>
              </tr>
            )) : 
            <tr>
              <td colSpan={8} className='text-center text-emphasis py-3'>ユーザーが存在しません</td>
            </tr>}
          </tbody>
        </table>
      </div>
    </ListPageLayout>
  );
}