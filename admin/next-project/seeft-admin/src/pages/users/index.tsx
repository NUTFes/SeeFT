import clsx from 'clsx';
import { useRouter } from 'next/router';

import { get } from '@api/api_methods';
import { User, Grade, Department, Bureau } from "@type/common";
import MainLayout from '@components/layout/MainLayout';
import { MdEdit, MdDeleteForever } from "react-icons/md";
import Button from '@components/common/Button';

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
  const router = useRouter();

  const addUserPageRouter = () => {
    router.push('users/add-user');
  }

  return (
    <MainLayout>
      <div className='w-full h-full bg-white-0 flex-col p-8'>
        <div className='items-center text-xl text-emphasis'>
          ユーザー一覧
        </div>
        <div className='items-center'>
          <div className='text-right pr-4'>
            <Button className='bg-surface-2 border-accent-2 text-right text-emphasis pr-4' onClick={addUserPageRouter}>
              ユーザー追加
            </Button>
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
              {users ? users.map((user: User, index) => (
                <tr key={user.id}>
                  <td
                    className={clsx(
                      'px-1 py-3',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b py-3',
                    )}
                  >
                    <p className='text-center text-sm text-emphasis'>{bureaus.find((bureau: Bureau) => (bureau.id === user.bureauID))?.bureau}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-3',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b py-3',
                    )}
                  >
                    <p className='text-center text-sm text-emphasis'>{user.name}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-3',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b py-3',
                    )}
                  >
                    <p className='text-center text-sm text-emphasis'>{grades.length - 1 ? grades.find((grade: Grade) => (grade.id === user.gradeID))?.grade : "erorr"}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-3',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b py-3',
                    )}
                  >
                    <p className='text-center text-sm text-emphasis'>{departments.length - 1 ? departments.find((department: Department) => (department.id === user.departmentID))?.department : "erorr"}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-3',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b py-3',
                    )}
                  >
                    <p className='text-center text-sm text-emphasis'>{user.studentNumber}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-3',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b py-3',
                    )}
                  >
                    <p className='text-center text-sm text-emphasis'>{user.tel}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-3',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b py-3',
                    )}
                  >
                    <div className='flex justify-items-center gap-4'>
                      <MdEdit />
                      <p className='text-center text-sm text-emphasis'>
                        編集
                      </p>
                    </div>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-3',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b py-3',
                    )}
                  >
                    <div className='flex justify-items-center gap-4'>
                      <MdDeleteForever />
                      <p className='text-center text-sm text-emphasis'>
                        削除
                      </p>
                    </div>
                  </td>
                </tr>
              )) : null}
            </tbody>
          </table>
        </div>
      </div>
    </MainLayout >
  );
}