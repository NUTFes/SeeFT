import clsx from 'clsx';
import { useRouter } from 'next/router';

import { get } from '@api/api_methods';
import { Task, User, Place } from "@type/common";
import { destroy } from '@api/task';
import { Button, DeleteButton, EditButton } from '@components/common';
import ListPageLayout from '@components/layout/ListPageLayout';

interface Props {
  tasks: Task[];
  users: User[];
  places: Place[];
}

export const getServerSideProps = async () => {
  const getTaskURL = process.env.SSR_API_URI + '/tasks';
  const getUserURL = process.env.SSR_API_URI + '/users';
  const getPlaceURL = process.env.SSR_API_URI + '/places';
  const taskRes = await get(getTaskURL);
  const userRes = await get(getUserURL);
  const placeRes = await get(getPlaceURL);

  return {
    props: {
      tasks: taskRes,
      users: userRes,
      places: placeRes,
    },
  };
};

export default function Users(props: Props) {
  const { tasks, users, places } = props;
  const router = useRouter();

  const addTaskPageRouter = () => {
    router.push('tasks/add-task');
  }

  const taskDetailPageRouter = (task: Task) => {
    router.push('tasks/' + task.id + '/detail-task');
  }

  const destroyTaskInformation = async (data: Task) => {
    const destroyTaskInformationUrl = process.env.CSR_API_URI + '/tasks';
    await destroy(destroyTaskInformationUrl, data);
    router.reload();
  };

  return (
    <ListPageLayout title='タスク一覧'>
      <div className='items-center'>
        <div className='text-right pr-4'>
          <Button className='bg-surface-2 border-accent-2 text-right text-emphasis pr-4 hover:bg-surface-1' onClick={addTaskPageRouter}>
            タスク追加
          </Button>
        </div>
      </div>
      <div className='p-5'>
        <table className='mb-5 w-full table-auto border-collapse'>
          <thead>
            <tr>
              <th className='w-1/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                <p className='text-center text-sm text-emphasis'>タスク名</p>
              </th>
              <th className='w-2/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                <p className='text-center text-sm text-emphasis'>集合場所</p>
              </th>
              <th className='w-2/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                <p className='text-center text-sm text-emphasis'>マニュアルURL</p>
              </th>
              <th className='w-2/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                <p className='text-center text-sm text-emphasis'>責任者</p>
              </th>
              <th className='w-1/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                <p className='text-center text-sm text-emphasis'>カラー</p>
              </th>
              <th className='w-2/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                <p className='text-center text-sm text-emphasis'>備考</p>
              </th>
              <th className='w-1/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3' />
              <th className='w-1/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3' />
            </tr>
          </thead>
          <tbody className='border border-x-white-0 border-b-accent-1 border-t-white-0'>
            {tasks ? tasks.map((task: Task, index) => (
              <tr key={task.id}>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === users.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <p className='text-center text-sm text-emphasis'>{task.task}</p>
                </td>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === users.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <p className='text-center text-sm text-emphasis'>{places.find((place: Place) => { place.id === task.placeID })?.place}</p>
                </td>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === users.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <p className='text-center text-sm text-emphasis'>{task.url}</p>
                </td>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === users.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <p className='text-center text-sm text-emphasis'>{users.length ? users.find((user: User) => (user.id === task.superviserID))?.name : "erorr"}</p>
                </td>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === users.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <div
                    className="mx-auto w-4 h-4 rounded-full border border-gray-400"
                    style={{ background: ('#' + task.color) }}
                  />
                </td>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === users.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <p className='text-center text-sm text-emphasis'>{task.remark}</p>
                </td>
                <td
                  className={clsx(
                    'px-1 py-2',
                    index === 0 ? 'pb-3 pt-4' : 'py-3',
                    index === users.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                  )}
                >
                  <EditButton onClick={() => { taskDetailPageRouter(task) }}>
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
                  <DeleteButton onClick={() => { destroyTaskInformation(task) }} >
                    削除
                  </DeleteButton>
                </td>
              </tr>
            )) : null}
          </tbody>
        </table>
      </div>
    </ListPageLayout>
  );
}