import clsx from 'clsx';
import { useRouter } from 'next/router';

import { get } from '@api/api_methods';
import { User, Task, Bureau, Time } from "@type/common";
import MainLayout from '@components/layout/MainLayout';
import Button from '@components/common/Button';
import { destroy } from '@api/user';
import { DeleteButton, EditButton, Select } from '@components/common';
import ListPageLayout from '@components/layout/ListPageLayout';
import { TimeItems, TimeScaleItem } from '@constants/timeItem';

interface Props {
  users: User[];
  tasks: Task[];
  bureaus: Bureau[];
}

export const getServerSideProps = async () => {
  const getUserURL = process.env.SSR_API_URI + '/users';
  const getTaskURL = process.env.SSR_API_URI + '/tasks';
  const getBureauURL = process.env.SSR_API_URI + '/bureaus';
  const userRes = await get(getUserURL);
  const taskRes = await get(getTaskURL);
  const bureauRes = await get(getBureauURL);

  return {
    props: {
      users: userRes,
      tasks: taskRes,
      bureaus: bureauRes,
    },
  };
};

export default function Users(props: Props) {
  const { users, tasks, bureaus } = props;
  const router = useRouter();

  const timeList: Time[] = TimeItems.slice((2 * 16), (3 * 16));

  const addShiftInformation = async (data: Shift) => {
    const addUserInformationUrl = process.env.CSR_API_URI + '/shifts';
    await post(addUserInformationUrl, data);
  };

  const updateShiftInformation = async (data: Shift) => {
    const putShiftInformationUrl = process.env.CSR_API_URI + '/shifts/' + data.id;
    await put(putShiftInformationUrl, data);
  };

  return (
    <ListPageLayout title='シフト一覧'>
      <div className='border border-accent-1 overflow-auto'>
        <div className='flex justify-between items-center p-4'>
          <div>
            {TimeScaleItem[1].time}
          </div>
          <div>
            {TimeScaleItem[2].time}
          </div>
          <div>
            {TimeScaleItem[3].time}
          </div>
        </div>
        <div className='p-5'>
          <table className='mb-5 w-full table-auto border-collapse'>
            <thead>
              <tr>
                <th className='w-1/12 bg-surface-2 border border-accent-1 py-3'>
                  <p className='text-center text-sm text-emphasis'>所属局</p>
                </th>
                <th className='w-2/12 bg-surface-2 border border-accent-1 py-3'>
                  <p className='text-center text-sm text-emphasis'>名前</p>
                </th>
                {timeList.map((time: Time, i: number) => (
                  i % 2 === 0 ?
                    <th className='w-3/64 bg-surface-2 border border-accent-1 py-3'>
                      <p className='text-center text-sm text-emphasis'>{time.time + '-'}</p>
                    </th>
                    :
                    <th className='w-3/64 bg-surface-2 border border-accent-1 py-3' />
                ))}
              </tr>
            </thead>
            <tbody className='border border-x-white-0 border-b-accent-1 border-t-white-0'>
              {users ? users.map((user: User, index) => (
                <tr key={user.id}>
                  <td
                    className={clsx(
                      'px-1 py-2 bg-surface-2',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                    )}
                  >
                    <p className='text-center text-sm text-emphasis'>{bureaus.find((bureau: Bureau) => (bureau.id === user.bureauID))?.bureau}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-2 bg-surface-2',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === users.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                    )}
                  >
                    <p className='text-center text-sm text-emphasis'>{user.name}</p>
                  </td>

                </tr>
              )) :
                'ユーザーが登録されていません'
              }
            </tbody>
          </table>
        </div>
      </div>
      <div className='w-1/3 border border-accent-1 my-2'>
        <div>
          シフトの編集方法
        </div>
        <div className='pl-4'>
          <div>
            <p>1. 入力するシフトを選択</p>
            <div className='flex justify-center items-center gap-4 pl-4'>
              <div className='w-1/3'>シフト検索</div>
              <Select className='w-full'>
                {tasks.map((data) => (
                  <option key={data.id} value={data.id}>
                    {data.task}
                  </option>
                ))}
              </Select>
            </div>
          </div>
          <div>
            <p>2. 活動時間を選択</p>
            <div className='flex justify-center items-center pl-4'>
              <div className='w-1/3'>活動時間</div>
              <Select className='w-full'>
                {TimeItems.map((data) => (
                  <option key={data.id} value={data.id}>
                    {data.time}
                  </option>
                ))}
              </Select>
            </div>
          </div>
          <div>
            <p>3. 開始時刻のセルをクリック</p>
          </div>
        </div>
      </div>
    </ListPageLayout>
  );
}